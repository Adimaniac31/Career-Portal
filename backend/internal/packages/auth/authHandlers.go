package auth

import (
	"encoding/json"
	"errors"
	"iiitn-career-portal/internal/config"
	"iiitn-career-portal/internal/models"
	"iiitn-career-portal/internal/packages/authorization"
	"iiitn-career-portal/internal/packages/keycloak"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GeneratePortalJWT(user models.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    user.ID,
		"role":       user.Role,
		"college_id": user.CollegeID,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func Signup(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email     string `json:"email" binding:"required,email"`
			Password  string `json:"password" binding:"required,min=8"`
			Name      string `json:"name" binding:"required"`
			CollegeID uint   `json:"college_id" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 1️⃣ Fetch college
		var college models.College
		if err := db.First(&college, req.CollegeID).Error; err != nil {
			c.JSON(400, gin.H{"error": "invalid college"})
			return
		}

		// 2️⃣ Validate domain
		if !strings.HasSuffix(req.Email, "@"+college.Domain) {
			c.JSON(400, gin.H{"error": "email domain not allowed"})
			return
		}

		// 3️⃣ Check duplicate
		var existing models.User
		if err := db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			c.JSON(400, gin.H{"error": "user already exists"})
			return
		}

		// 4️⃣ Create Keycloak user
		kcUserID, err := keycloak.CreateUser(
			cfg,
			req.Email,
			req.Password,
			req.Name,
		)
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to create identity"})
			return
		}

		// Ensure cleanup on ANY failure below
		defer func() {
			if err != nil {
				keycloak.DeleteUser(cfg, kcUserID)
			}
		}()

		// 5️⃣ Assign role
		err = keycloak.AssignRealmRole(cfg, kcUserID, "student")
		if err != nil {
			c.JSON(500, gin.H{"error": "role assignment failed"})
			return
		}

		// 6️⃣ DB transaction
		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(500, gin.H{"error": "db transaction failed"})
			return
		}

		user := models.User{
			KeycloakID: kcUserID,
			Email:      req.Email,
			Name:       req.Name,
			CollegeID:  &college.ID,
			Role:       "student",
		}

		if err = tx.Create(&user).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "db insert failed"})
			return
		}

		if err = tx.Commit().Error; err != nil {
			c.JSON(500, gin.H{"error": "db commit failed"})
			return
		}

		// 7️⃣ Success
		c.JSON(201, gin.H{
			"message": "Signup successful. Please login.",
		})
	}
}

func Login(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 1️⃣ Get token from Keycloak
		tokenResp, err := keycloak.PasswordGrant(cfg, req.Email, req.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid credentials"})
			return
		}

		// 2️⃣ Verify token
		claims, err := keycloak.VerifyAccessToken(cfg, tokenResp.AccessToken)
		if err != nil {
			log.Println(err)
			c.JSON(401, gin.H{"error": "invalid token"})
			return
		}

		// 3️⃣ Lookup DB user
		var user models.User
		if err := db.Where("keycloak_id = ?", claims.Subject).
			First(&user).Error; err != nil {
			c.JSON(403, gin.H{"error": "user not registered"})
			return
		}

		// 4️⃣ Issue portal JWT
		portalToken, err := GeneratePortalJWT(user, cfg.JWTSecret)
		if err != nil {
			c.JSON(500, gin.H{"error": "login failed"})
			return
		}

		// 5️⃣ Set cookie
		c.SetCookie(
			"portal_token",
			portalToken,
			86400,
			"/",
			"",
			false,
			true, // HttpOnly
		)

		c.JSON(200, gin.H{"message": "login successful"})
	}
}

func SSOLogin(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		q := url.Values{}
		q.Set("client_id", cfg.ClientID)
		q.Set("response_type", "code")
		q.Set("scope", "openid email profile")
		q.Set("redirect_uri", cfg.BackendBaseURL+"/api/auth/sso/callback")

		authURL := cfg.BaseURL +
			"/realms/" + cfg.Realm +
			"/protocol/openid-connect/auth?" +
			q.Encode()

		c.Redirect(302, authURL)
	}
}

func SSOCallback(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(400, gin.H{"error": "missing code"})
			return
		}

		// 1️⃣ Exchange code for token
		token, err := exchangeCode(cfg, code)
		if err != nil {
			c.JSON(401, gin.H{"error": "sso exchange failed"})
			return
		}

		// 2️⃣ Verify token (same verifier you already trust)
		claims, err := keycloak.VerifyAccessToken(cfg, token.AccessToken)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			return
		}

		// 3️⃣ Find user in DB
		var user models.User
		if err := db.Where("keycloak_id = ?", claims.Subject).
			First(&user).Error; err != nil {
			c.JSON(403, gin.H{"error": "user not registered"})
			return
		}

		// 4️⃣ Issue portal JWT
		portalToken, err := GeneratePortalJWT(user, cfg.JWTSecret)
		if err != nil {
			c.JSON(500, gin.H{"error": "login failed"})
			return
		}

		// 5️⃣ Set cookie
		c.SetCookie(
			"portal_token",
			portalToken,
			86400,
			"/",
			"",
			false,
			true,
		)

		// 6️⃣ Redirect to frontend
		c.Redirect(302, cfg.FrontendURL)
	}
}

func exchangeCode(cfg config.Config, code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("redirect_uri", cfg.BackendBaseURL+"/api/auth/sso/callback")

	resp, err := http.PostForm(
		cfg.BaseURL+"/realms/"+cfg.Realm+
			"/protocol/openid-connect/token",
		data,
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	// log.Println("[SSO TOKEN EXCHANGE]", resp.StatusCode, string(body))

	if resp.StatusCode != 200 {
		return nil, errors.New("code exchange failed")
	}

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func Me(db *gorm.DB, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1️⃣ Read cookie
		tokenStr, err := c.Cookie("portal_token")
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthenticated"})
			return
		}

		// 2️⃣ Verify portal JWT (unchanged)
		claims, err := authorization.VerifyPortalJWT(tokenStr, cfg.JWTSecret)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			return
		}

		// 3️⃣ Build auth context
		authCtx, err := authorization.BuildAuthContext(claims)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid auth context"})
			return
		}

		// 4️⃣ Fetch minimal user info
		var user struct {
			ID        uint
			Name      string
			Email     string
			Role      string
			CollegeID uint
		}

		err = db.
			Model(&models.User{}).
			Select("id, name, email, role, college_id").
			Where("id = ?", authCtx.UserID).
			First(&user).Error

		if err != nil {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		// 5️⃣ Profile completion check
		var profileComplete bool
		db.
			Model(&models.StudentProfile{}).
			Select("profile_complete").
			Where("user_id = ?", user.ID).
			Scan(&profileComplete)

		// 6️⃣ Respond
		c.JSON(200, gin.H{
			"id":               user.ID,
			"name":             user.Name,
			"email":            user.Email,
			"role":             user.Role,
			"college_id":       user.CollegeID,
			"profile_complete": profileComplete,
		})
	}
}
