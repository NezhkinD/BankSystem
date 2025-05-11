// @title BankSystem API
// @version 1.0
// @description API Server for Bank System

// @host localhost:8080
// @BasePath /api
// @schemes http

package handlers

import (
	"BankSystem/internal/dto"
	"BankSystem/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET"))

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrEmailExists    = errors.New("user with this email already exists")
	ErrUsernameExists = errors.New("username is already taken")
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Параметры регистрации"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingByEmail, _ := h.userService.GetUserByEmail(req.Email)
	if existingByEmail != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrEmailExists.Error()})
		return
	}

	errUsernameExists, _ := h.userService.GetUserByUsername(req.Username)
	if errUsernameExists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrUsernameExists.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	err = h.userService.RegisterUser(req.Email, req.Username, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		logrus.Error(err.Error())
		return
	}

	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		logrus.Error(err.Error())
		return
	}

	tokenString, err := h.generateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}

// Login godoc
// @Summary Авторизация пользователя
// @Description Возвращает JWT токен после успешной авторизации
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Email и пароль"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := h.generateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"username": user.Username,
		},
	})
}

func (h *AuthHandler) generateJWT(id uint, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
