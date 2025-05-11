package account

import (
	"BankSystem/internal/models"
	"BankSystem/internal/repositories"
	"errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type AccountService struct {
	accountRepo *repositories.AccountRepository
	userRepo    *repositories.UserRepository
}

func NewAccountService(repo *repositories.AccountRepository) *AccountService {
	return &AccountService{accountRepo: repo}
}

func (s *AccountService) CreateAccount(userID uint, balance decimal.Decimal) error {
	account := &models.Account{
		UserID:   userID,
		Balance:  balance,
		Currency: "RUB",
	}
	account.DeletedAt = gorm.DeletedAt{Valid: false, Time: time.Time{}}
	return s.accountRepo.Create(account)
}

func (s *AccountService) Deposit(userID uint, amount decimal.Decimal) (decimal.Decimal, error) {
	account, err := s.accountRepo.FindByUserID(userID)
	if err != nil || account == nil {
		return decimal.Zero, errors.New("account not found")
	}

	account.Balance = account.Balance.Add(amount)
	if err := s.accountRepo.Update(account); err != nil {
		return decimal.Zero, err
	}

	return account.Balance, nil
}

func (s *AccountService) IsAccountExists(userID uint) (bool, error) {
	account, err := s.accountRepo.FindByUserID(userID)
	if err != nil {
		return false, err
	}
	return account != nil, nil
}

func (s *AccountService) GetByID(id uint, userId uint) (*models.Account, error) {
	acc, err := s.accountRepo.FindByIdAndUserID(id, userId)
	if err != nil {
		return nil, errors.New("account not found")
	}

	return acc, nil
}

func (s *AccountService) GetAccountsByUserID(userID uint) ([]*models.Account, error) {
	return s.accountRepo.FindAllByUserID(userID)
}
