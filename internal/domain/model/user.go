package model

import "gorm.io/gorm"

// User model
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     Role   `gorm:"not null"`
	Tasks    []Task `gorm:"foreignKey:UserID"`
}

type Role string

const (
	Manager    Role = "manager"
	Technician Role = "technician"
)
