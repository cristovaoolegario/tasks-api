package repository

import (
	"fmt"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
)

type MockUserRepository struct {
	Users map[uint]*model.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[uint]*model.User),
	}
}

func (m *MockUserRepository) FindByID(id uint) (*model.User, error) {
	if user, exists := m.Users[id]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) FindByUsername(username string) (*model.User, error) {
	for _, user := range m.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MockUserRepository) Create(user *model.User) error {
	m.Users[user.ID] = user
	return nil
}

func (m *MockUserRepository) Update(user *model.User) error {
	if _, exists := m.Users[user.ID]; exists {
		m.Users[user.ID] = user
		return nil
	}
	return fmt.Errorf("user not found")
}

func (m *MockUserRepository) Delete(id uint) error {
	if _, exists := m.Users[id]; exists {
		delete(m.Users, id)
		return nil
	}
	return fmt.Errorf("user not found")
}
