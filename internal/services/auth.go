package services

import (
	"BankSystem/internal/models"
	"BankSystem/internal/security"
	"gorm.io/gorm"
)

func RegisterUser(db *gorm.DB, username, email, password string) (*models.User, error) {
	hashedPass, _ := security.HashPassword(password)
	user := &models.User{
		Username: username,
		Email:    email,
		Password: hashedPass,
	}
	result := db.Create(user)
	return user, result.Error
}

func AuthenticateUser(db *gorm.DB, email, password string) (*models.User, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if !security.CheckPasswordHash(password, user.Password) {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}
