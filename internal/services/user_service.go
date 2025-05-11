package services

import (
	"BankSystem/internal/models"
	"BankSystem/internal/repositories"
	"errors"
	"gorm.io/gorm"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) RegisterUser(email, username, password string) error {
	user := &models.User{
		Email:    email,
		Username: username,
		Password: password,
	}
	user.DeletedAt = gorm.DeletedAt{Valid: false, Time: time.Time{}}
	return s.userRepo.Create(user)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, err
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepo.FindByUserName(username)
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, err
}
