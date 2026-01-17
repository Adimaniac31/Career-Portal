package authorization

import "github.com/gin-gonic/gin"

func GetAuth(c *gin.Context) AuthContext {
	return c.MustGet("auth").(AuthContext)
}
