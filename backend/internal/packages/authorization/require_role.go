package authorization

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCtx, ok := c.Get("auth")
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				gin.H{"error": "authentication required"},
			)
			return
		}

		user := authCtx.(AuthContext)

		for _, role := range allowedRoles {
			if user.Role == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(
			http.StatusForbidden,
			gin.H{"error": "insufficient permissions"},
		)
	}
}
