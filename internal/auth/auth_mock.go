package auth

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
)

// MockAuthService is a mock implementation of the AuthService interface.
type MockAuthService struct {
	GenerateJWTFunc  func(user, role string) (string, error)
	LoginFunc        func(username, password string) (string, error)
	FindByUsernameFn func(username string) (*model.User, error)
}

func (m *MockAuthService) GenerateJWT(user, role string) (string, error) {
	return m.GenerateJWTFunc(user, role)
}

func (m *MockAuthService) Login(username, password string) (string, error) {
	return m.LoginFunc(username, password)
}

func (m *MockAuthService) FindByUsername(username string) (*model.User, error) {
	return m.FindByUsernameFn(username)
}
