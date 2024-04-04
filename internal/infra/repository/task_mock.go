package repository

import (
	"errors"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
)

type MockTaskRepository struct {
	Tasks map[uint]*model.Task
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		Tasks: make(map[uint]*model.Task),
	}
}

func (m *MockTaskRepository) Create(task *model.Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}
	if _, exists := m.Tasks[task.ID]; exists {
		return errors.New("task already exists")
	}
	taskId := uint(len(m.Tasks) + 1)
	m.Tasks[taskId] = task
	return nil
}

func (m *MockTaskRepository) Update(task *model.Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}
	if _, exists := m.Tasks[task.ID]; !exists {
		return errors.New("task does not exist")
	}
	m.Tasks[task.ID] = task
	return nil
}

func (m *MockTaskRepository) Delete(taskID uint) error {
	if _, exists := m.Tasks[taskID]; !exists {
		return errors.New("task does not exist")
	}
	delete(m.Tasks, taskID)
	return nil
}

func (m *MockTaskRepository) FindByID(taskID uint) (*model.Task, error) {
	if task, exists := m.Tasks[taskID]; exists {
		return task, nil
	}
	return nil, errors.New("task does not exist")
}

func (m *MockTaskRepository) FindByUser(userID uint) ([]*model.Task, error) {
	var tasks []*model.Task
	for _, task := range m.Tasks {
		if task.UserID == userID {
			tasks = append(tasks, task)
		}
	}
	if len(tasks) == 0 {
		return nil, errors.New("no tasks found for user")
	}

	return tasks, nil
}
