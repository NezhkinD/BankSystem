package account

import (
	"BankSystem/internal/models"
	"BankSystem/internal/repositories"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type AccountService struct {
	accountRepo *repositories.AccountRepository
	userRepo    *repositories.UserRepository
	log         *logrus.Logger
}

func NewAccountService(repo *repositories.AccountRepository, log *logrus.Logger) *AccountService {
	return &AccountService{
		accountRepo: repo,
		log:         log,
	}
}

func (s *AccountService) CreateAccount(userID uint, balance decimal.Decimal) error {
	account := &models.Account{
		UserID:   userID,
		Balance:  balance,
		Currency: "RUB",
	}
	account.DeletedAt = gorm.DeletedAt{Valid: false, Time: time.Time{}}
	logrus.Info("created new account for user" + strconv.Itoa(int(userID)))
	return s.accountRepo.Create(account)
}

func (s *AccountService) Deposit(id uint, userID uint, amount decimal.Decimal) (decimal.Decimal, error) {
	account, err := s.accountRepo.FindByIdAndUserID(id, userID)
	if err != nil || account == nil {
		return decimal.Zero, errors.New("account not found")
	}

	account.Balance = account.Balance.Add(amount)
	if err := s.accountRepo.Update(account); err != nil {
		return decimal.Zero, err
	}

	logrus.Info("user " + strconv.Itoa(int(userID)) + " has been deposited successfully")
	return account.Balance, nil
}

func (s *AccountService) Withdraw(id uint, userID uint, amount decimal.Decimal) (decimal.Decimal, error) {
	account, err := s.accountRepo.FindByIdAndUserID(id, userID)
	if err != nil || account == nil {
		return decimal.Zero, errors.New("account not found")
	}

	if account.Balance.LessThan(amount) {
		return decimal.Zero, errors.New("insufficient funds")
	}

	account.Balance = account.Balance.Sub(amount)
	if err := s.accountRepo.Update(account); err != nil {
		return decimal.Zero, err
	}

	logrus.Info("user " + strconv.Itoa(int(userID)) + " has been withdraw successfully")
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

func (s *AccountService) IsUserAccount(id uint, userId uint) error {
	acc, err := s.accountRepo.FindByIdAndUserID(id, userId)
	if err != nil {
		return errors.New("account not found")
	}

	if acc != nil {
		return nil
	}

	return nil
}

func (s *AccountService) GetAccountsByUserID(userID uint) ([]*models.Account, error) {
	return s.accountRepo.FindAllByUserID(userID)
}

func (s *AccountService) Transfer(userID uint, fromAccID uint, toAccID uint, amount decimal.Decimal) error {
	account, err := s.accountRepo.FindByIdAndUserID(fromAccID, userID)
	if err != nil || account == nil {
		return errors.New("account not found")
	}

	return s.accountRepo.WithinTransaction(func(tx *gorm.DB) error {
		fromAccount, err := s.accountRepo.FindByIDWithLock(fromAccID)
		if err != nil {
			return errors.New("sender account not found")
		}

		toAccount, err := s.accountRepo.FindByIDWithLock(toAccID)
		if err != nil {
			return errors.New("recipient account not found")
		}

		if fromAccount.Balance.LessThan(amount) {
			return errors.New("insufficient funds")
		}

		fromAccount.Balance = fromAccount.Balance.Sub(amount)
		toAccount.Balance = toAccount.Balance.Add(amount)

		if err := s.accountRepo.UpdateWithTx(tx, fromAccount); err != nil {
			return err
		}

		if err := s.accountRepo.UpdateWithTx(tx, toAccount); err != nil {
			return err
		}

		tx.Create(&models.Transaction{
			FromAccountID:   fromAccID,
			ToAccountID:     toAccID,
			Amount:          amount,
			TransactionType: "transfer",
			Currency:        "RUB",
		})

		return nil
	})
}
