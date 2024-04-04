package auth

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/gin-gonic/gin"
)

// MockAuthService is a mock implementation of the AuthService interface.
type MockAuthService struct {
	GenerateJWTMock              func(user, role, id string) (string, error)
	LoginMock                    func(username, password string) (string, error)
	FindByUsernameMock           func(username string) (*model.User, error)
	ExtractUserIdFromContextMock func(ctx *gin.Context) (uint, error)
}

func (m *MockAuthService) GenerateJWT(user, role, id string) (string, error) {
	return m.GenerateJWTMock(user, role, id)
}

func (m *MockAuthService) Login(username, password string) (string, error) {
	return m.LoginMock(username, password)
}

func (m *MockAuthService) FindByUsername(username string) (*model.User, error) {
	return m.FindByUsernameMock(username)
}

func (m *MockAuthService) ExtractUserIdFromContext(ctx *gin.Context) (uint, error) {
	return m.ExtractUserIdFromContextMock(ctx)
}
