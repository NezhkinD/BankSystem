package services

import (
	"BankSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

func TransferFunds(db *gorm.DB, fromID, toID uint, amount float64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var from, to models.Account
		if err := tx.First(&from, fromID).Error; err != nil {
			return err
		}
		if from.Balance < amount {
			return errors.New("insufficient funds")
		}
		from.Balance -= amount
		to.Balance += amount
		tx.Save(&from)
		tx.Save(&to)
		return nil
	})
}
