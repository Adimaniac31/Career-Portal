package main

import (
	"iiitn-career-portal/internal/cache"
	"iiitn-career-portal/internal/config"
	"iiitn-career-portal/internal/database"
	"iiitn-career-portal/internal/packages/admin"
	"iiitn-career-portal/internal/packages/auth"
	"iiitn-career-portal/internal/packages/authorization"
	"iiitn-career-portal/internal/packages/colleges"
	"iiitn-career-portal/internal/packages/jobs"
	"iiitn-career-portal/internal/packages/profile"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	db := database.Connect(cfg.DB)

	database.Migrate(db)
	redisClient := cache.NewRedisClient(cfg.Redis)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		auth.RegisterRoutes(api, cfg, db)
		colleges.RegisterRoutes(api, db)
		protected := api.Group("/")
		protected.Use(authorization.RequireAuth(cfg))
		{
			admin.RegisterRoutes(protected, db)
			profile.RegisterRoutes(protected, db, cfg)
			jobs.RegisterRoutes(protected, db, redisClient)
		}
	}

	log.Println("Server running on port:", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
