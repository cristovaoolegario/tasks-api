package service

import (
	"errors"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/infra/repository"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserServiceImp(t *testing.T) {
	t.Run("Should create a valid user", func(t *testing.T) {
		repository := repository.NewMockUserRepository()
		service := NewUserService(repository)

		user := model.User{
			Username: "user1",
			Password: "password",
			Role:     model.Manager,
		}

		err := service.CreateUser(&user)

		assert.NoError(t, err)
		assert.Len(t, repository.Users, 1)

		foundUser, err := service.FindByUsername("user1")

		assert.NotNil(t, foundUser)
		assert.NoError(t, err)
	})

	t.Run("Should return an error When the username is empty", func(t *testing.T) {
		repository := repository.NewMockUserRepository()
		service := NewUserService(repository)

		user := model.User{
			Username: "",
			Password: "password",
			Role:     model.Manager,
		}

		err := service.CreateUser(&user)

		assert.Error(t, err)
		assert.Equal(t, errors.New("invalid user"), err)
		assert.NotContains(t, repository.Users, "")
	})

	t.Run("Should return an error When the password is empty", func(t *testing.T) {
		repository := repository.NewMockUserRepository()
		service := NewUserService(repository)

		user := model.User{
			Username: "user1",
			Password: "",
			Role:     model.Manager,
		}

		err := service.CreateUser(&user)

		assert.Error(t, err)
		assert.Equal(t, errors.New("invalid user"), err)
		assert.NotContains(t, repository.Users, "user1")
	})

	t.Run("Should return an error When the role is invalid", func(t *testing.T) {
		repository := repository.NewMockUserRepository()
		service := NewUserService(repository)

		user := model.User{
			Username: "user1",
			Password: "password",
			Role:     "invalidRole",
		}

		err := service.CreateUser(&user)

		assert.Error(t, err)
		assert.Equal(t, errors.New("invalid user"), err)
		assert.NotContains(t, repository.Users, "user1")
	})
}
