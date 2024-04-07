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
	taskService         service.TaskService
	authService         auth.Service
	notificationService service.ManagerNotificationService
}

// NewTaskController creates a new instance of TaskController.
func NewTaskController(taskService service.TaskService, authService auth.Service, notificationService service.ManagerNotificationService) *TaskController {
	return &TaskController{
		taskService:         taskService,
		authService:         authService,
		notificationService: notificationService,
	}
}

// CreateTaskHandler handles the creation of a new task.
// @Summary Create a new task
// @Description Add a new task for the authenticated user
// @Tags task
// @Accept json
// @Produce json
// @Param task body dto.Task true "Create Task"
// @Success 201 {object} dto.Task
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/tasks [post]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (tc *TaskController) CreateTaskHandler(ctx *gin.Context) {
	var task dto.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _, err := tc.authService.ExtractUserFromContext(ctx)
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
// @Summary Update an existing task
// @Description Update task details for the authenticated user
// @Tags task
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body dto.Task true "Update Task"
// @Success 200 {object} dto.Task
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 403 {object} map[string]interface{} "Forbidden"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/tasks/{id} [put]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
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

	userId, _, err := tc.authService.ExtractUserFromContext(ctx)
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

	convertedTask := &dto.Task{
		ID:            updatedTask.ID,
		Summary:       updatedTask.Summary,
		PerformedDate: updatedTask.PerformedDate,
		UserID:        updatedTask.UserID,
	}

	go tc.notificationService.Notification(userId, convertedTask)
	ctx.JSON(http.StatusOK, &convertedTask)
}

// DeleteTaskHandler handles the deletion of a task.
// @Summary Delete a task
// @Description Delete a task for the authenticated user
// @Tags task
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string "Task deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/tasks/{id} [delete]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
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

	ctx.JSON(http.StatusOK, &dto.Task{
		ID:            task.ID,
		Summary:       task.Summary,
		PerformedDate: task.PerformedDate,
		UserID:        task.UserID,
	})
}

// FindByUserID handles GET requests for finding tasks by a user's ID
// @Summary Get a task by ID
// @Description Get a task by its ID for the authenticated user
// @Tags task
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} dto.Task "Task found"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 404 {object} map[string]interface{} "Not Found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/tasks/{id} [get]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (tc *TaskController) FindByUserID(ctx *gin.Context) {
	userId, role, err := tc.authService.ExtractUserFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if role == string(model.Manager) {
		tc.GetAllTasksHandler(ctx)
		return
	}

	tasks, err := tc.taskService.FindTasksByUser(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "tasks for the given user not found"})
		return
	}

	convertedTasks := toDtoTasks(tasks)
	ctx.JSON(http.StatusOK, convertedTasks)
}

// GetAllTasksHandler godoc
// @Summary Get all tasks with pagination
// @Description Retrieve all tasks with pagination
// @Tags task
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(10)
// @Success 200 {array} dto.Task "List of tasks"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /api/tasks [get]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (tc *TaskController) GetAllTasksHandler(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	tasks, err := tc.taskService.FindPaginatedTasks(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	convertedTasks := toDtoTasks(tasks)
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

func toDtoTasks(tasks []*model.Task) []*dto.Task {
	var convertedTasks []*dto.Task
	for _, task := range tasks {
		convertedTasks = append(convertedTasks, &dto.Task{
			ID:            task.ID,
			Summary:       task.Summary,
			PerformedDate: task.PerformedDate,
			UserID:        task.UserID,
		})
	}
	return convertedTasks
}
