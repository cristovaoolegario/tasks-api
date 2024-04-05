package controller

import (
	"encoding/json"
	"errors"
	"github.com/cristovaoolegario/tasks-api/internal/auth"
	"github.com/cristovaoolegario/tasks-api/internal/domain/dto"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/domain/service"
	"github.com/cristovaoolegario/tasks-api/internal/infra/kafka"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"
)

func TestTaskController_CreateTaskHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("Should create the task properly", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 1, nil
			}}
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, mockAuth, nil)

		task := dto.Task{Summary: "Test task", UserID: 1}
		w := performRequest(taskController.CreateTaskHandler, http.MethodPost, "/users", "", "", task)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response dto.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(0), response.ID)
		assert.Equal(t, "Test task", response.Summary)
	})

	t.Run("Should return 400 When Body is malformed", func(t *testing.T) {
		taskController := NewTaskController(nil, nil, nil)

		w := performRequest(taskController.CreateTaskHandler, http.MethodPost, "/users", "", "", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 500 When UserId is not set", func(t *testing.T) {
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 0, errors.New("test error")
			}}
		taskController := NewTaskController(nil, mockAuth, nil)

		task := dto.Task{Summary: "Test task", UserID: 1}
		w := performRequest(taskController.CreateTaskHandler, http.MethodPost, "/users", "", "", task)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTaskController_UpdateTaskHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("Should return 400 When path id is malformed", func(t *testing.T) {
		taskController := NewTaskController(nil, nil, nil)

		w := performRequest(taskController.UpdateTaskHandler, http.MethodPut, "/users", "", "", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 400 When body is malformed", func(t *testing.T) {
		taskController := NewTaskController(nil, nil, nil)

		w := performRequest(taskController.UpdateTaskHandler, http.MethodPut, "/users", "", "1", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 400 When id Isn't from an existing task", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, nil, nil)

		task := dto.Task{Summary: "Test task", UserID: 1}
		w := performRequest(taskController.UpdateTaskHandler, http.MethodPut, "/users", "", "1", task)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 500 When UserId is not set in claims", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		repo.Tasks = map[uint]*model.Task{
			1: {Model: gorm.Model{ID: 1}},
		}
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 0, errors.New("test error")
			}}
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, mockAuth, nil)

		task := dto.Task{Summary: "Test task", UserID: 1}
		w := performRequest(taskController.UpdateTaskHandler, http.MethodPut, "/users", "", "1", task)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Should return 403 When user id is different than task's user ", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		repo.Tasks = map[uint]*model.Task{
			1: {
				Model:  gorm.Model{ID: 1},
				UserID: 3,
			},
		}
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 1, nil
			}}
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, mockAuth, nil)

		task := dto.Task{Summary: "Test task", UserID: 1}
		w := performRequest(taskController.UpdateTaskHandler, http.MethodPut, "/users", "", "1", task)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Should return 200 and updated task", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		repo.Tasks = map[uint]*model.Task{
			1: {
				Model:   gorm.Model{ID: 1},
				UserID:  1,
				Summary: "not edited summary",
			},
		}
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 1, nil
			}}
		mockedService := service.NewTaskService(repo)
		mockedNotification := service.NewManagerNotificationService("",
			&kafka.ProducerMock{PublishMessageMock: func(topic string, message []byte) error {
				return nil
			}})
		taskController := NewTaskController(mockedService, mockAuth, mockedNotification)

		performedTime := time.Now()
		task := dto.Task{Summary: "edited summary", UserID: 1, PerformedDate: &performedTime}
		w := performRequest(taskController.UpdateTaskHandler, http.MethodPut, "/users", "", "1", task)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "edited summary", response.Summary)
	})
}

func TestTaskController_DeleteTaskHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Should return 400 When path id is malformed", func(t *testing.T) {
		taskController := NewTaskController(nil, nil, nil)

		w := performRequest(taskController.DeleteTaskHandler, http.MethodDelete, "/users", "", "", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 200 When task is deleted", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		repo.Tasks = map[uint]*model.Task{
			1: {
				Model:   gorm.Model{ID: 1},
				UserID:  1,
				Summary: "summary",
			},
		}
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 1, nil
			}}
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, mockAuth, nil)

		w := performRequest(taskController.DeleteTaskHandler, http.MethodDelete, "/users", "", "1", nil)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestTaskController_FindByID(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("Should return 400 When path id is malformed", func(t *testing.T) {
		taskController := NewTaskController(nil, nil, nil)

		w := performRequest(taskController.FindByID, http.MethodGet, "/users", "", "", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return 404 When id Isn't from an existing task", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, nil, nil)

		w := performRequest(taskController.FindByID, http.MethodGet, "/users", "", "1", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Should return 500 When UserId is not set in claims", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		repo.Tasks = map[uint]*model.Task{
			1: {Model: gorm.Model{ID: 1}},
		}
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 0, errors.New("test error")
			}}
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, mockAuth, nil)

		w := performRequest(taskController.FindByID, http.MethodGet, "/users", "", "1", nil)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Should return 403 When user id is different than task's user ", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		repo.Tasks = map[uint]*model.Task{
			1: {
				Model:  gorm.Model{ID: 1},
				UserID: 3,
			},
		}
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 1, nil
			}}
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, mockAuth, nil)

		w := performRequest(taskController.FindByID, http.MethodGet, "/users", "", "1", nil)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Should return 200 and task", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		repo.Tasks = map[uint]*model.Task{
			1: {
				Model:   gorm.Model{ID: 1},
				UserID:  1,
				Summary: "not edited summary",
			},
		}
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 1, nil
			}}
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, mockAuth, nil)

		w := performRequest(taskController.FindByID, http.MethodGet, "/users", "", "1", nil)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dto.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), response.ID)
		assert.Equal(t, "not edited summary", response.Summary)
	})
}

func TestTaskController_FindByUserID(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("Should return 500 When UserId is not set in claims", func(t *testing.T) {
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 0, errors.New("test error")
			}}
		taskController := NewTaskController(nil, mockAuth, nil)

		w := performRequest(taskController.FindByUserID, http.MethodGet, "/users", "", "", nil)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Should return 404 When id Isn't from an existing task", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		mockedService := service.NewTaskService(repo)
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 1, nil
			}}
		taskController := NewTaskController(mockedService, mockAuth, nil)

		w := performRequest(taskController.FindByUserID, http.MethodGet, "/users", "", "", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("Should return 200 and  task", func(t *testing.T) {
		repo := repository.NewMockTaskRepository()
		repo.Tasks = map[uint]*model.Task{
			1: {
				Model:   gorm.Model{ID: 1},
				UserID:  1,
				Summary: "not edited summary",
			},
		}
		mockAuth := &auth.MockAuthService{
			ExtractUserIdFromContextMock: func(ctx *gin.Context) (uint, error) {
				return 1, nil
			}}
		mockedService := service.NewTaskService(repo)
		taskController := NewTaskController(mockedService, mockAuth, nil)

		w := performRequest(taskController.FindByUserID, http.MethodGet, "/users", "", "", nil)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*dto.Task
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 1)
	})
}
