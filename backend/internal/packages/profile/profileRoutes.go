package profile

import (
	"iiitn-career-portal/internal/config"
	"iiitn-career-portal/internal/models"
	"iiitn-career-portal/internal/packages/authorization"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, cfg config.Config) {
	profile := rg.Group("/profile")
	profile.Use(authorization.RequireRole(string(models.Student)))
	{
		profile.GET("", GetProfile(db))
		profile.PATCH("", UpdateProfile(db))
	}
}
