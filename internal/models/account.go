package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	UserID   uint            `db:"user_id"  json:"user_id"`
	Balance  decimal.Decimal `db:"balance"  json:"balance"`
	Currency string          `db:"currency"  json:"currency"`
}
