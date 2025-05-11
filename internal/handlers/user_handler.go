package handlers

import (
	_ "BankSystem/internal/models"
	"BankSystem/internal/repositories"
	"BankSystem/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userRepo    *repositories.UserRepository
	authService *services.AuthService
}

func NewUserHandler(repo *repositories.UserRepository, authService *services.AuthService) *UserHandler {
	return &UserHandler{
		userRepo:    repo,
		authService: authService,
	}
}

// GetCurrentUser godoc
// @Summary Получить данные текущего пользователя
// @Description Возвращает информацию о пользователе на основе токена
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /user/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}
