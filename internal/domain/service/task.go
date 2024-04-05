package service

import (
	"errors"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
)

type TaskService interface {
	CreateTask(task *model.Task) error
	UpdateTask(task *model.Task) error
	DeleteTask(taskID uint) error
	FindTaskByID(taskID uint) (*model.Task, error)
	FindTasksByUser(userID uint) ([]*model.Task, error)
	FindPaginatedTasks(page, pageSize int) ([]*model.Task, error)
}

type TaskServiceImp struct {
	repository repository.TaskRepository
}

func NewTaskService(repository repository.TaskRepository) *TaskServiceImp {
	return &TaskServiceImp{
		repository: repository,
	}
}

func (s *TaskServiceImp) CreateTask(task *model.Task) error {
	if task.Summary == "" || len(task.Summary) > 2500 {
		return errors.New("invalid task summary")
	}
	if task.UserID == 0 {
		return errors.New("invalid user ID")
	}
	return s.repository.Create(task)
}

func (s *TaskServiceImp) UpdateTask(task *model.Task) error {
	if task.ID == 0 {
		return errors.New("invalid task ID")
	}
	return s.repository.Update(task)
}

func (s *TaskServiceImp) DeleteTask(taskID uint) error {
	if taskID == 0 {
		return errors.New("invalid task ID")
	}
	return s.repository.Delete(taskID)
}

func (s *TaskServiceImp) FindTaskByID(taskID uint) (*model.Task, error) {
	return s.repository.FindByID(taskID)
}

func (s *TaskServiceImp) FindTasksByUser(userID uint) ([]*model.Task, error) {
	return s.repository.FindByUser(userID)
}

func (s *TaskServiceImp) FindPaginatedTasks(page, pageSize int) ([]*model.Task, error) {
	return s.repository.FindAll(page, pageSize)
}
