package models

import "time"

type ApplicationIntent struct {
	ID uint `gorm:"primaryKey"`

	JobID     uint `gorm:"index;not null;uniqueIndex:uniq_intent"`
	StudentID uint `gorm:"index;not null;uniqueIndex:uniq_intent"`
	CollegeID uint `gorm:"index;not null"`

	ExpiresAt time.Time

	CreatedAt time.Time
}
