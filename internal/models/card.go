package models

import "gorm.io/gorm"

type Card struct {
	gorm.Model
	AccountID uint
	Number    string
	CVV       string
	Expiry    string
}
