package repository

import "github.com/cristovaoolegario/tasks-api/internal/domain/model"

// UserRepository  interface for user persistence
type UserRepository interface {
	FindByID(id uint) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint) error
}
