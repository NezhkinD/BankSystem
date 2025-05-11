package services

import (
	"BankSystem/internal/models"
	"BankSystem/internal/repositories"
	"BankSystem/internal/security"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

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

func (s *AuthService) GetCurrentUser(c *gin.Context) (*models.User, error) {
	email, exists := c.Get("email")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	user, err := s.userRepo.FindByEmail(email.(string))
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
	}

	return user, nil
}
