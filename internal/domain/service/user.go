package service

import (
	"errors"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
)

type UserService interface {
	CreateUser(user *model.User) error
	FindByUsername(userName string) (*model.User, error)
}

type UserServiceImp struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserServiceImp {
	return &UserServiceImp{
		repository: repository,
	}
}

func (s *UserServiceImp) CreateUser(user *model.User) error {
	if user.Username == "" || user.Password == "" ||
		(user.Role != model.Manager && user.Role != model.Technician) {
		return errors.New("invalid user")
	}

	return s.repository.Create(user)
}

func (s *UserServiceImp) FindByUsername(userName string) (*model.User, error) {
	return s.repository.FindByUsername(userName)
}
