package models

import (
	"time"

	"gorm.io/datatypes"
)

type Job struct {
	ID uint `gorm:"primaryKey"`

	CollegeID uint `gorm:"not null;index"`

	Title string `gorm:"not null"`

	Company string `gorm:"not null"`

	JobType JobType   `gorm:"type:varchar(20);not null"`
	Domain  JobDomain `gorm:"type:varchar(20);not null"`

	EligibleBatches datatypes.JSON `gorm:"not null"`

	CTC     *float64
	Stipend *float64

	RegistrationFormURL *string `gorm:"type:text"`

	Description string `gorm:"type:text"`
	Rounds      datatypes.JSON

	IsActive bool `gorm:"default:true"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
