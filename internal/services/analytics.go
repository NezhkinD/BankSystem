package services

import (
	"BankSystem/internal/models"
	"gorm.io/gorm"
)

func CalculateUserBalance(db *gorm.DB, userID uint) (float64, error) {
	var total float64
	err := db.Model(&models.Account{}).Where("user_id = ?", userID).Select("SUM(balance)").Scan(&total).Error
	return total, err
}
