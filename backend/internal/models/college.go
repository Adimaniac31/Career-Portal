package models

import "time"

type College struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex"`

	Domain string `gorm:"uniqueIndex"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
