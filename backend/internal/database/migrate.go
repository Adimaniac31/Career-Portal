package database

import (
	"log"

	"iiitn-career-portal/internal/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.College{},
		&models.User{},
		&models.StudentProfile{},
		&models.Job{},
		&models.Application{},
		&models.Notification{},
		&models.ApplicationIntent{},
		&models.Application{},
	)
	if err != nil {
		log.Fatal("migration failed:", err)
	}
}
