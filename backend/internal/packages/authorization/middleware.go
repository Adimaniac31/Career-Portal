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
		// log.Println(tokenStr)
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

		// 3️⃣ Extract claims safely
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		role, _ := claims["role"].(string)

		var collegeID *uint
		if cid, exists := claims["college_id"]; exists && cid != nil {
			val := uint(cid.(float64))
			collegeID = &val
		}

		// 4️⃣ Attach auth context
		c.Set("auth", AuthContext{
			UserID:    uint(userIDFloat),
			Role:      role,
			CollegeID: collegeID,
		})

		// 5️⃣ Continue
		c.Next()
	}
}
