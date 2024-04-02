package infra

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// User model
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"type:enum('manager', 'technician');not null"`
	Tasks    []Task `gorm:"foreignKey:UserID"`
}

// Task model
type Task struct {
	gorm.Model
	Summary       string    `gorm:"type:varchar(2500);not null"`
	PerformedDate time.Time `gorm:"not null"`
	UserID        uint      // Foreign key for User
}

func InitDB(connectionString string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{}, &Task{})

	return db
}
