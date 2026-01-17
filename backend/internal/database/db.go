package database

import (
	"log"
	"strings" // Add this

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) *gorm.DB {
	// Trim whitespace and hidden carriage returns
	cleanDSN := strings.TrimSpace(dsn)
	log.Println(cleanDSN)

	db, err := gorm.Open(postgres.Open(cleanDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	return db
}
