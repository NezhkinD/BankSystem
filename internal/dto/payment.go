package dto

import (
	"github.com/shopspring/decimal"
	"time"
)

type PaymentNotification struct {
	To        string
	Name      string
	CardLast4 string
	Amount    decimal.Decimal
	Balance   decimal.Decimal
	Date      time.Time
}
