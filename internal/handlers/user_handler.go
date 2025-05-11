package handlers

import (
	_ "BankSystem/internal/models"
	"BankSystem/internal/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userRepo *repositories.UserRepository
}

func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{userRepo: repo}
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
	email, exists := c.Get("email")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := h.userRepo.FindByEmail(email.(string))
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
