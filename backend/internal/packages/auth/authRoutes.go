package auth

import (
	"iiitn-career-portal/internal/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(rg *gin.RouterGroup, cfg config.Config, db *gorm.DB) {

	auth := rg.Group("/auth")
	{
		auth.POST("/signup", Signup(db, cfg))
		auth.POST("/login", Login(db, cfg))
		auth.GET("/sso/login", SSOLogin(cfg))
		auth.GET("/sso/callback", SSOCallback(db, cfg))
		auth.GET("/me", Me(db, cfg))
	}
}
