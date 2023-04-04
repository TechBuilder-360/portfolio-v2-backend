package repository

import (
	"context"
	"errors"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/database"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/model"
	"gorm.io/gorm"
)

type IAccountRepository interface {
	GetById(ctx context.Context, id string) (*model.Account, error)
	GetByEmail(email string, ctx context.Context) (*model.Account, error)
	Create(account *model.Account, ctx context.Context) error
	Update(account *model.Account, ctx context.Context) error
	WithTx(tx *gorm.DB) IAccountRepository
}

type accountRepository struct {
	db *gorm.DB
}

func (a *accountRepository) WithTx(tx *gorm.DB) IAccountRepository {
	return &accountRepository{db: tx}
}

// NewAccountRepository will a account repository
func NewAccountRepository() IAccountRepository {
	return &accountRepository{db: database.DB()}
}

func (a *accountRepository) GetById(ctx context.Context, id string) (*model.Account, error) {
	var account model.Account
	if err := a.db.WithContext(ctx).First(&account, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, errors.New("an error occurred")
	}

	return &account, nil
}

func (a *accountRepository) GetByEmail(email string, ctx context.Context) (*model.Account, error) {
	var user model.Account
	if err := a.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, errors.New("an error occurred")
	}

	return &user, nil
}

func (a *accountRepository) Create(account *model.Account, ctx context.Context) error {
	if err := a.db.WithContext(ctx).Create(account).Error; err != nil {
		return err
	}

	return nil
}

func (a *accountRepository) Update(account *model.Account, ctx context.Context) error {
	if err := a.db.WithContext(ctx).Save(account).Error; err != nil {
		return err
	}

	return nil
}
