package repository

import (
	"context"
	"errors"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/database"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetByID(id string, ctx context.Context) (*model.User, error)
	GetByAccountID(accountID string, ctx context.Context) (*model.User, error)
	Create(user *model.User, ctx context.Context) error
	WithTx(tx *gorm.DB) IUserRepository
}

type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) WithTx(tx *gorm.DB) IUserRepository {
	return &userRepository{db: tx}
}

// NewUserRepository will a user repository
func NewUserRepository() IUserRepository {
	return &userRepository{db: database.DB()}
}

func (u *userRepository) GetByID(id string, ctx context.Context) (*model.User, error) {
	var user model.User
	if err := u.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, errors.New("an error occurred")
	}

	return &user, nil
}

func (u *userRepository) Create(user *model.User, ctx context.Context) error {
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		return errors.New("an error occurred")
	}

	return nil
}

func (u *userRepository) GetByAccountID(accountID string, ctx context.Context) (*model.User, error) {
	var user model.User
	if err := u.db.WithContext(ctx).Where("account_id = ?", accountID).First(&user).Error; err != nil {
		return nil, errors.New("an error occurred")
	}

	return &user, nil
}
