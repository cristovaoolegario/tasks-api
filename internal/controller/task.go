package controller

import (
	"github.com/cristovaoolegario/tasks-api/internal/auth"
	"github.com/cristovaoolegario/tasks-api/internal/domain/dto"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"net/http"
	"strconv"

	"github.com/cristovaoolegario/tasks-api/internal/domain/service"
	"github.com/gin-gonic/gin"
)

// TaskController handles HTTP requests related to tasks.
type TaskController struct {
	taskService service.TaskService
	authService auth.Service
}

// NewTaskController creates a new instance of TaskController.
func NewTaskController(taskService service.TaskService, authService auth.Service) *TaskController {
	return &TaskController{
		taskService: taskService,
		authService: authService,
	}
}

// CreateTaskHandler handles the creation of a new task.
func (tc *TaskController) CreateTaskHandler(ctx *gin.Context) {
	var task dto.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := tc.authService.ExtractUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	task.UserID = userId
	convertedTask := task.ToModel()

	if err := tc.taskService.CreateTask(convertedTask); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	task.ID = convertedTask.ID
	ctx.JSON(http.StatusCreated, task)
}

// UpdateTaskHandler handles the update of an existing task.
func (tc *TaskController) UpdateTaskHandler(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var updatedTaskDto dto.Task
	if err := ctx.ShouldBindJSON(&updatedTaskDto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	oldTask, err := tc.taskService.FindTaskByID(uint(taskID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userId, err := tc.authService.ExtractUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if oldTask.UserID != userId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	updatedTask := extractTaskDataToUpdate(oldTask, updatedTaskDto)

	if err := tc.taskService.UpdateTask(updatedTask); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	ctx.JSON(http.StatusOK, &dto.Task{
		ID:            updatedTask.ID,
		Summary:       updatedTask.Summary,
		PerformedDate: updatedTask.PerformedDate,
		UserID:        updatedTask.UserID,
	})
}

// DeleteTaskHandler handles the deletion of a task.
func (tc *TaskController) DeleteTaskHandler(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := tc.taskService.DeleteTask(uint(taskID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// FindByID handles GET requests for finding a task by its ID
func (tc *TaskController) FindByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	task, err := tc.taskService.FindTaskByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	userId, err := tc.authService.ExtractUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if task.UserID != userId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	ctx.JSON(http.StatusOK, &dto.Task{
		ID:            task.ID,
		Summary:       task.Summary,
		PerformedDate: task.PerformedDate,
		UserID:        task.UserID,
	})
}

// FindByUserID handles GET requests for finding tasks by a user's ID
func (tc *TaskController) FindByUserID(ctx *gin.Context) {
	userId, err := tc.authService.ExtractUserIdFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	tasks, err := tc.taskService.FindTasksByUser(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "tasks for the given user not found"})
		return
	}

	var convertedTasks []*dto.Task
	for _, task := range tasks {
		convertedTasks = append(convertedTasks, &dto.Task{
			ID:            task.ID,
			Summary:       task.Summary,
			PerformedDate: task.PerformedDate,
			UserID:        task.UserID,
		})
	}

	ctx.JSON(http.StatusOK, convertedTasks)
}

func extractTaskDataToUpdate(oldTask *model.Task, updatedTaskDto dto.Task) *model.Task {
	updatedTask := oldTask
	if updatedTaskDto.Summary != "" {
		updatedTask.Summary = updatedTaskDto.Summary
	}
	if updatedTaskDto.PerformedDate != nil {
		updatedTask.PerformedDate = updatedTaskDto.PerformedDate
	}
	return updatedTask
}
