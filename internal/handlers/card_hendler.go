package handlers

import (
	"BankSystem/internal/dto"
	"BankSystem/internal/repositories"
	"BankSystem/internal/services"
	account_service "BankSystem/internal/services/account"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CardHandler struct {
	userService    *services.UserService
	accountService *account_service.AccountService
	cardService    *services.CardService
	authService    *services.AuthService
	accountRepo    *repositories.AccountRepository
	cardRepo       *repositories.CardRepository
}

func NewCardHandler(userService *services.UserService, accountService *account_service.AccountService, cardService *services.CardService, authService *services.AuthService, accountRepo *repositories.AccountRepository) *CardHandler {
	return &CardHandler{
		userService:    userService,
		accountService: accountService,
		cardService:    cardService,
		authService:    authService,
		accountRepo:    accountRepo,
	}
}

// CreateCard godoc
// @Summary Создание новой карты
// @Description Создаёт новую карту для текущего пользователя
// @Tags card
// @Security BearerAuth
// @Produce json
// @Param request body dto.CreateCardRequest true "Данные для выпуска карты"
// @Success 201 {object} dto.NewCardResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /card/create [post]
func (h *CardHandler) CreateCard(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req dto.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountRepo.FindByIdAndUserID(req.AccountID, user.ID)
	if err != nil || account == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "account not found or user does not have access to the account"})
		return
	}

	card, err := h.cardService.GenerateCard(req.AccountID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not generate card: " + err.Error()})
		return
	}

	response := dto.NewCardResponse{
		Number: card.CardNumber,
		CVV:    card.Cvv,
		Expiry: card.ExpiredAt.Format("01/06"),
	}

	c.JSON(http.StatusCreated, response)
}

// GetCards godoc
// @Summary Получить все карты пользователя
// @Description Возвращает список всех карт с балансом по аккаунтам
// @Tags card
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.CardWithBalance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cards [get]
func (h *CardHandler) GetCards(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	cards, err := h.cardService.GetCardsByUserID(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Could not load cards"})
		return
	}

	for i := range cards {
		cards[i].CVV = "***"
	}

	c.JSON(http.StatusOK, cards)
}

// PayWithCard godoc
// @Summary Оплата с помощью карты
// @Description Выполняет оплату, проверяя данные карты и списывая средства
// @Tags payment
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CardPaymentRequest true "Данные карты и сумма"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 402 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payment/card [post]
func (h *CardHandler) PayWithCard(c *gin.Context) {
	user, err := h.authService.GetCurrentUser(c)
	if err != nil || user == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req dto.CardPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBalance, err := h.cardService.PayWithCard(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "payment failed: " + err.Error()})
		return
	}

	response := dto.CardPaymentResponse{
		Result:  true,
		Message: "payment done, new balance: " + newBalance.String(),
	}
	c.JSON(http.StatusOK, response)
}
