package service

import (
	"errors"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestTaskServiceImp_CreateTask(t *testing.T) {
	t.Run("Should create a valid task", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		task := &model.Task{
			UserID:  1,
			Summary: "Do something",
		}

		err := service.CreateTask(task)

		assert.NoError(t, err)
		assert.Len(t, repo.Tasks, 1)

		foundTask, err := service.FindTaskByID(1)

		assert.NotNil(t, foundTask)
		assert.NoError(t, err)
	})

	t.Run("Should return an error when creating a task with empty summary", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		task := &model.Task{
			UserID:  1,
			Summary: "",
		}

		err := service.CreateTask(task)

		assert.Error(t, err)
		assert.Equal(t, errors.New("invalid task summary"), err)
		assert.Empty(t, repo.Tasks)
	})

	t.Run("Should return an error when creating a task with empty user ID", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		task := &model.Task{
			UserID:  0,
			Summary: "Some summary",
		}

		err := service.CreateTask(task)

		assert.Error(t, err)
		assert.Equal(t, errors.New("invalid user ID"), err)
		assert.Empty(t, repo.Tasks)
	})
}

func TestTaskServiceImp_FindTasks(t *testing.T) {
	t.Run("Should return a task by ID", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		task := &model.Task{
			UserID:  1,
			Summary: "Do something",
		}

		repo.Tasks[1] = task

		foundTask, err := service.FindTaskByID(1)

		assert.NotNil(t, foundTask)
		assert.NoError(t, err)
		assert.Equal(t, task, foundTask)
	})

	t.Run("Should return an error when getting a task by non-existent ID", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		foundTask, err := service.FindTaskByID(999)

		assert.Nil(t, foundTask)
		assert.Error(t, err)
		assert.Equal(t, errors.New("task does not exist"), err)
	})

	t.Run("Should list tasks by user ID", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		userID := uint(1)
		task1 := &model.Task{
			UserID:  userID,
			Summary: "Task 1",
		}
		task2 := &model.Task{
			UserID:  userID,
			Summary: "Task 2",
		}

		repo.Tasks[1] = task1
		repo.Tasks[2] = task2

		tasks, err := service.FindTasksByUser(userID)

		assert.NoError(t, err)
		assert.Len(t, tasks, 2)
		assert.Contains(t, tasks, task1)
		assert.Contains(t, tasks, task2)
	})

	t.Run("Should return an error when listing tasks by non-existent user ID", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		tasks, err := service.FindTasksByUser(0)

		assert.Nil(t, tasks)
		assert.Error(t, err)
		assert.Equal(t, errors.New("no tasks found for user"), err)
	})
}

func TestTaskServiceImp_DeleteTask(t *testing.T) {
	t.Run("Should delete a task with valid ID", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		originalTask := &model.Task{
			UserID:  1,
			Summary: "Do something",
		}

		_ = service.CreateTask(originalTask)

		taskID := uint(1)

		err := service.DeleteTask(taskID)

		assert.NoError(t, err)
	})

	t.Run("Should return an error when deleting a task with invalid ID", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		taskID := uint(0)

		err := service.DeleteTask(taskID)

		assert.Error(t, err)
		assert.Equal(t, errors.New("invalid task ID"), err)
	})
}

func TestTaskServiceImp_UpdateTask(t *testing.T) {
	t.Run("Should update a task with valid ID", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		originalTask := &model.Task{
			UserID:  1,
			Summary: "Do something",
		}

		_ = service.CreateTask(originalTask)

		taskToUpdate := &model.Task{
			Model:   gorm.Model{ID: 1},
			Summary: "Updated summary",
		}

		err := service.UpdateTask(taskToUpdate)

		assert.NoError(t, err)
	})

	t.Run("Should return an error when updating a task with invalid ID", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		service := NewTaskService(repo)

		task := &model.Task{
			Model: gorm.Model{ID: 0},
		}

		err := service.UpdateTask(task)

		assert.Error(t, err)
		assert.Equal(t, errors.New("invalid task ID"), err)
	})
}
