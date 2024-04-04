package model

import (
	"gorm.io/gorm"
	"time"
)

// Task model
type Task struct {
	gorm.Model
	Summary       string    `gorm:"type:varchar(2500);not null"`
	PerformedDate time.Time `gorm:"not null"`
	UserID        uint
	User          User `gorm:"foreignKey:UserID"`
}
