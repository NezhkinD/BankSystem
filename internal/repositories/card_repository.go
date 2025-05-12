// internal/repositories/user_repository.go

package repositories

import (
	"BankSystem/internal/dto"
	"BankSystem/internal/models"
	"errors"
	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(entity *models.Card) error {
	result := r.db.Create(entity)
	return result.Error
}

func (r *CardRepository) FindByAccountId(accountId uint) (*models.Card, error) {
	var card models.Card
	result := r.db.Where("account_id = ?", accountId).First(&card)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &card, result.Error
}

func (r *CardRepository) FindByCardNumber(cardNumber string) (*models.Card, error) {
	var card models.Card
	result := r.db.Where("card_number = ?", cardNumber).First(&card)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &card, result.Error
}

func (r *CardRepository) GetCardsByUserID(userID uint) ([]dto.CardResponse, error) {
	var result []dto.CardResponse
	err := r.db.Table("cards").
		Select("cards.id, cards.card_number, accounts.id as account_id, accounts.balance, accounts.currency, cards.expired_at,  cards.created_at").
		Joins("JOIN accounts ON cards.account_id = accounts.id").
		Where("accounts.user_id = ?", userID).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *CardRepository) FindByPlainCardNumber(number string, key string) (*models.Card, error) {
	var card models.Card

	query := `
        SELECT * FROM cards 
        WHERE pgp_sym_decrypt(decode(card_number, 'hex'), ?) = ?
    `
	result := r.db.Raw(query, key, number).Scan(&card)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("card not found")
	}
	return &card, nil
}

func (r *CardRepository) Update(card *models.Card) error {
	return r.db.Save(card).Error
}

func (r *CardRepository) Delete(id uint) error {
	return r.db.Delete(&models.Card{}, id).Error
}
