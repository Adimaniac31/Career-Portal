package models

import (
	"time"
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	KeycloakID string `gorm:"uniqueIndex;not null"`
	Email      string `gorm:"uniqueIndex"`
	Name       string

	CollegeID *uint
	College   College
	Role      string `gorm:"type:varchar(20);default:'student'" json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
