package repository

import "github.com/cristovaoolegario/tasks-api/internal/domain/model"

type TaskRepository interface {
	Create(task *model.Task) error
	Update(task *model.Task) error
	Delete(taskID uint) error
	FindByID(taskID uint) (*model.Task, error)
	FindByUser(userID uint) ([]*model.Task, error)
}
