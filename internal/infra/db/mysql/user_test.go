package mysql

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&model.User{}, &model.Task{})
	return db
}

func TestUserRepository(t *testing.T) {
	dbConnect := setupTestDB()
	db, err := dbConnect.DB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepository := NewUserRepository(dbConnect)
	newUser := model.User{Username: "testuser", Password: "testpassword", Role: model.Manager}
	foundUser := &model.User{}

	t.Run("Should create user", func(t *testing.T) {
		err = userRepository.Create(&newUser)
		assert.NoError(t, err)
	})

	t.Run("Should return user When exists", func(t *testing.T) {
		foundUser, err := userRepository.FindByUsername("testuser")
		assert.NoError(t, err)
		assert.NotNil(t, foundUser)
		assert.Equal(t, newUser.Username, foundUser.Username)
		assert.Equal(t, newUser.Role, foundUser.Role)
		assert.Equal(t, newUser.Password, foundUser.Password)

		foundUserByID, err := userRepository.FindByID(foundUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, foundUserByID)
		assert.Equal(t, foundUser.ID, foundUserByID.ID)
	})

	t.Run("Should be able to update existing user", func(t *testing.T) {
		foundUser.Role = model.Technician
		err = userRepository.Update(foundUser)
		assert.NoError(t, err)

		updatedUser, err := userRepository.FindByID(foundUser.ID)
		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.Equal(t, foundUser.Role, updatedUser.Role)
	})

	t.Run("Should be able to delete existing user", func(t *testing.T) {
		err = userRepository.Delete(foundUser.ID)
		assert.NoError(t, err)

		_, err = userRepository.FindByID(foundUser.ID)
		assert.Error(t, err)
	})
}
