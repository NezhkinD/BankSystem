package services

import (
	"BankSystem/internal/models"
	"BankSystem/internal/repositories"
	"BankSystem/internal/services/account"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	userRepo       *repositories.UserRepository
	accountService *account.AccountService
	log            *logrus.Logger
}

func NewUserService(repo *repositories.UserRepository, accountService *account.AccountService, log *logrus.Logger) *UserService {
	return &UserService{
		userRepo:       repo,
		accountService: accountService,
		log:            log,
	}
}

func (s *UserService) RegisterUser(email, username, password string) error {
	user := &models.User{
		Email:    email,
		Username: username,
		Password: password,
	}
	user.DeletedAt = gorm.DeletedAt{Valid: false, Time: time.Time{}}
	logrus.Info("Create user " + username)
	return s.userRepo.Create(user)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if user == nil {
		logrus.Warn("User with email " + email + " not found")
		return nil, ErrUserNotFound
	}

	logrus.Info("Found user with email " + email + " not found")
	return user, err
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepo.FindByUserName(username)
	if user == nil {
		logrus.Warn("User with username " + username + " not found")
		return nil, ErrUserNotFound
	}
	logrus.Info("Found user with username " + username + " not found")
	return user, err
}
