// internal/repositories/user_repository.go

package repositories

import (
	"BankSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *models.Account) error {
	result := r.db.Create(account)
	return result.Error
}

func (r *AccountRepository) FindByID(id uint) (*models.Account, error) {
	var account models.Account
	result := r.db.Where("id = ?", id).First(&account)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &account, result.Error
}

func (r *AccountRepository) FindByUserID(id uint) (*models.Account, error) {
	var account models.Account
	result := r.db.Where("user_id = ?", id).First(&account)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &account, result.Error
}

func (r *AccountRepository) FindByIdAndUserID(id uint, userId uint) (*models.Account, error) {
	var account models.Account
	result := r.db.Where("id = ? AND user_id = ?", id, userId).First(&account)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &account, result.Error
}

func (r *AccountRepository) FindAllByUserID(userID uint) ([]*models.Account, error) {
	var accounts []*models.Account
	result := r.db.Where("user_id = ?", userID).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}

func (r *AccountRepository) Update(model *models.Account) error {
	return r.db.Save(model).Error
}

// FindByIDWithLock — получение аккаунта с блокировкой строки
func (r *AccountRepository) FindByIDWithLock(id uint) (*models.Account, error) {
	var account models.Account
	result := r.db.Where("id = ?", id).Set("gorm:query_option", "FOR UPDATE").First(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

// UpdateWithTx — обновление в рамках транзакции
func (r *AccountRepository) UpdateWithTx(tx *gorm.DB, account *models.Account) error {
	return tx.Save(account).Error
}

// WithinTransaction — обёртка для выполнения в транзакции
func (r *AccountRepository) WithinTransaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}
