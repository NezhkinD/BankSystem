package handlers

import (
	"BankSystem/internal/dto"
	"BankSystem/internal/services"
	accountService "BankSystem/internal/services/account"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
)

type AccountHandler struct {
	accountService *accountService.AccountService
	userService    *services.UserService
	authService    *services.AuthService
}

func NewAccountHandler(accountService *accountService.AccountService, userService *services.UserService, authService *services.AuthService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		userService:    userService,
		authService:    authService,
	}
}

// CreateAccount godoc
// @Summary Создание аккаунта
// @Description Создаёт новый аккаунт для текущего пользователя
// @Tags account
// @Security BearerAuth
// @Produce json
// @Success 201 {object} models.Account
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account/create [post]
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	initialBalance := decimal.NewFromInt(0)
	err = h.accountService.CreateAccount(user.ID, initialBalance)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not create account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
}

// Deposit godoc
// @Summary Пополнение баланса
// @Description Пополняет баланс аккаунта текущего пользователя
// @Tags account
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.DepositRequest true "Сумма пополнения"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /account/deposit [post]
func (h *AccountHandler) Deposit(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req dto.DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBalance, err := h.accountService.Deposit(user.ID, decimal.NewFromFloat(req.Amount))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not deposit funds"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Deposit successful",
		"new_balance": newBalance.String(),
	})
}
