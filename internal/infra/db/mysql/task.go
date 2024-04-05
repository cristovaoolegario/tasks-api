package mysql

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"gorm.io/gorm"
)

// TaskRepository is the GORM implementation of TaskRepository for MySQL.
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new instance of TaskRepository with a connected MySQL database.
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// FindByID finds a task by its ID.
func (repo *TaskRepository) FindByID(id uint) (*model.Task, error) {
	var task model.Task
	result := repo.db.First(&task, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &task, nil
}

// Create adds a new task to the database.
func (repo *TaskRepository) Create(task *model.Task) error {
	result := repo.db.Create(task)
	return result.Error
}

// Update modifies an existing task.
func (repo *TaskRepository) Update(task *model.Task) error {
	result := repo.db.Save(task)
	return result.Error
}

// Delete removes a task by its ID.
func (repo *TaskRepository) Delete(id uint) error {
	result := repo.db.Delete(&model.Task{}, id)
	return result.Error
}

// FindByUser retrieves tasks for a specific user ID.
func (repo *TaskRepository) FindByUser(userID uint) ([]*model.Task, error) {
	var tasks []*model.Task
	if err := repo.db.Where("user_id = ?", userID).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindAll retrieves all tasks from the database.
func (repo *TaskRepository) FindAll(page, pageSize int) ([]*model.Task, error) {
	var tasks []*model.Task
	if err := repo.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}
