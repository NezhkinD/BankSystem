package models

import (
	"gorm.io/gorm"
	"time"
)

type Card struct {
	gorm.Model
	AccountId  uint      `db:"account_id" json:"account_id"`
	CardNumber string    `db:"card_number" json:"card_number"`
	Cvv        string    `db:"cvv" json:"cvv"`
	ExpiredAt  time.Time `db:"expired_at" json:"expired_at"`
}
