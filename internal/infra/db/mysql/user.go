package mysql

import (
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository is the GORM implementation of UserRepository for MySQL.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository with a connected MySQL database.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByID finds a user by their ID.
func (repo *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	result := repo.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindByUsername finds a user by their username.
func (repo *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	result := repo.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Create adds a new user to the database.
func (repo *UserRepository) Create(user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	result := repo.db.Create(user)
	return result.Error
}

// Update modifies an existing user.
func (repo *UserRepository) Update(user *model.User) error {
	result := repo.db.Save(user)
	return result.Error
}

// Delete removes a user by their ID.
func (repo *UserRepository) Delete(id uint) error {
	result := repo.db.Delete(&model.User{}, id)
	return result.Error
}
