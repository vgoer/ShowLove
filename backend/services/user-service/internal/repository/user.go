package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"showlove/services/user-service/internal/model"
)

var (
	ErrNotFound      = errors.New("用户不存在")
	ErrAlreadyExists = errors.New("用户已存在")
)

// UserRepository defines the data access interface for users.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

// userRepository implements UserRepository with GORM.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	result := r.db.WithContext(ctx).Save(user)
	return result.Error
}
