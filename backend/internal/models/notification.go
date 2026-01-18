package models

import (
	"time"

	"gorm.io/datatypes"
)

type Notification struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"index;not null"`

	Type NotificationType `gorm:"type:varchar(50);not null"`

	Payload datatypes.JSON `gorm:"not null"`

	IsRead bool `gorm:"default:false"`

	CreatedAt time.Time
}
