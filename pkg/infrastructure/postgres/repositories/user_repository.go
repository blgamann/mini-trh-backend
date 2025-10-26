package repositories

import (
	"errors"
	"mini-trh-backend/pkg/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindById(id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

// FindAll retrieves all users (with pagination support)
func (r *UserRepository) FindAll(limit, offset int) ([]entities.User, error) {
	var users []entities.User
	err := r.db.Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Update updates a user's information
func (r *UserRepository) Update(user *entities.User) error {
	return r.db.Save(user).Error
}

// Delete soft-deletes a user
func (r *UserRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entities.User{}, id).Error
}

// Count returns the total number of users
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&entities.User{}).Count(&count).Error
	return count, err
}

// ExistsByEmail checks if a user with the given email exists
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
