package mysql

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
)

func TestTaskRepository(t *testing.T) {
	dbConnect := setupTestDB()
	db, err := dbConnect.DB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepository := NewUserRepository(dbConnect)
	newUser := &model.User{Username: "testuser", Password: "testpassword", Role: model.Manager}
	_ = userRepository.Create(newUser)

	repo := NewTaskRepository(dbConnect)

	t.Run("Should create task when valid", func(t *testing.T) {
		task := &model.Task{
			UserID:        1,
			Summary:       "Do something",
			PerformedDate: time.Now(),
		}

		err = repo.Create(task)
		assert.NoError(t, err)
	})

	t.Run("Should find existing user", func(t *testing.T) {
		// Find task by ID
		foundTask, err := repo.FindByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, foundTask)
	})

	t.Run("Should update existing task", func(t *testing.T) {
		testTask := &model.Task{
			Model:         gorm.Model{ID: 1},
			UserID:        1,
			Summary:       "Do something",
			PerformedDate: time.Now(),
		}

		err = repo.Update(testTask)

		assert.NoError(t, err)
	})

	t.Run("Should delete exiting task", func(t *testing.T) {
		err = repo.Delete(1)
		assert.NoError(t, err)
	})

}