package models

import "time"

type Application struct {
	ID uint `gorm:"primaryKey"`

	JobID uint `gorm:"not null;index;uniqueIndex:uniq_student_job"`
	Job   Job  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	StudentID uint `gorm:"not null;index;uniqueIndex:uniq_student_job"`
	Student   User `gorm:"foreignKey:StudentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	CollegeID uint `gorm:"not null;index"`

	Status ApplicationStatus `gorm:"type:varchar(20);not null"`

	ResumeSnapshotURL string `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
