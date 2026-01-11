package main

import (
	"iiitn-career-portal/internal/config"
	"iiitn-career-portal/internal/database"
	"iiitn-career-portal/internal/packages/auth"
	"iiitn-career-portal/internal/packages/authorization"
	"iiitn-career-portal/internal/packages/colleges"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	// 1️⃣ Connect DB
	db := database.Connect(cfg.DB)

	// 2️⃣ Run migrations
	database.Migrate(db)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		auth.RegisterRoutes(api, cfg, db)
		protected := api.Group("/")
		protected.Use(authorization.RequireAuth(cfg))
		{
			colleges.RegisterRoutes(protected, db)
		}
		// colleges.RegisterRoutes(api, db)
	}

	log.Println("Server running on port:", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
