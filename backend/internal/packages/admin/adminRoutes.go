package admin

import (
	"iiitn-career-portal/internal/packages/authorization"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB) {

	// Admin-only
	admin := rg.Group("/colleges")
	admin.Use(authorization.RequireRole("admin"))
	{
		admin.POST("", CreateCollege(db))
	}
}
