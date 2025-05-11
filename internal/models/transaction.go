package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	FromAccountID   uint            `db:"from_account_id"  json:"from_account_id"`
	ToAccountID     uint            `db:"to_account_id"  json:"to_account_id"`
	Amount          decimal.Decimal `db:"amount"  json:"amount"`
	TransactionType string          `db:"transaction_type"  json:"transaction_type"`
	Currency        string          `db:"currency"  json:"currency"`
}
