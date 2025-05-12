package dto

import (
	"github.com/shopspring/decimal"
	"time"
)

type CreateCardRequest struct {
	AccountID uint `json:"account_id" binding:"required,gt=0"`
}

type NewCardResponse struct {
	Number string `json:"number"`
	CVV    string `json:"cvv"`
	Expiry string `json:"expiry"`
}

type CardResponse struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	AccountID uint            `json:"account_id" gorm:"not null"`
	Balance   decimal.Decimal `json:"balance" gorm:"column:balance;type:numeric(12,2);default:0.00"`
	Number    string          `json:"number" gorm:"column:card_number;type:text;unique"`
	Currency  string          `json:"currency" gorm:"column:currency;type:text;unique"`
	CVV       string          `json:"cvv" gorm:"type:text"`
	ExpiredAt time.Time       `json:"expired_at" gorm:"column:expired_at"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at"`
}

type CardPaymentRequest struct {
	CardNumber string    `json:"card_number" binding:"required"`
	Cvv        string    `json:"cvv" binding:"required,len=3"`
	Amount     float64   `json:"amount" binding:"required,gt=0"`
	ExpiredAt  time.Time `json:"expired_at"`
}

type CardPaymentResponse struct {
	Result  bool   `json:"result" binding:"required"`
	Message string `json:"message" binding:"required,len=3"`
}
