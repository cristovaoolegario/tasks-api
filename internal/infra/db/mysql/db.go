package mysql

import (
	"fmt"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(connectionString string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{}, &model.Task{})

	return db
}
