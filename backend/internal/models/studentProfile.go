package models

import "time"

type StudentProfile struct {
	UserID          uint     `gorm:"primaryKey"`
	CGPA            *float32 // nullable
	ResumeURL       *string  // nullable (MinIO URL)
	ProfileComplete bool     `gorm:"default:false"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
