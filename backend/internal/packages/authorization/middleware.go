package authorization

import (
	"net/http"

	"iiitn-career-portal/internal/config"

	"github.com/gin-gonic/gin"
)

func RequireAuth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1️⃣ Read cookie
		tokenStr, err := c.Cookie("portal_token")
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "authentication required"},
			)
			return
		}

		// 2️⃣ Verify JWT
		claims, err := VerifyPortalJWT(tokenStr, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid session"},
			)
			return
		}

		// 3️⃣ Extract user_id
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid token"},
			)
			return
		}

		// 4️⃣ Extract role (STRICT)
		role, ok := claims["role"].(string)
		if !ok || role == "" {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "invalid role"},
			)
			return
		}

		// 5️⃣ Extract college_id (SAFE)
		var collegeID *uint
		if cid, ok := claims["college_id"].(float64); ok {
			val := uint(cid)
			collegeID = &val
		}

		// 6️⃣ Attach auth context (POINTER)
		c.Set("auth", &AuthContext{
			UserID:    uint(userIDFloat),
			Role:      role,
			CollegeID: collegeID,
		})

		c.Next()
	}
}
