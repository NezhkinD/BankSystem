package services

import (
	"BankSystem/internal/dto"
	"BankSystem/internal/models"
	"BankSystem/internal/repositories"
	"BankSystem/internal/utils"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"

	accountservice "BankSystem/internal/services/account"
)

type CardService struct {
	db             *gorm.DB
	cardRepo       *repositories.CardRepository
	accountRepo    *repositories.AccountRepository
	userRepo       *repositories.UserRepository
	accountService *accountservice.AccountService
	mailService    *MailService
	encryptKey     string
	log            *logrus.Logger
}

func NewCardService(
	db *gorm.DB,
	cardRepo *repositories.CardRepository,
	accountRepo *repositories.AccountRepository,
	userRepo *repositories.UserRepository,
	accountService *accountservice.AccountService,
	mailService *MailService,
	encryptKey string,
	log *logrus.Logger) *CardService {
	return &CardService{
		db:             db,
		cardRepo:       cardRepo,
		accountRepo:    accountRepo,
		userRepo:       userRepo,
		accountService: accountService,
		mailService:    mailService,
		encryptKey:     encryptKey,
		log:            log,
	}
}

func (s *CardService) GetCardsByUserID(userID uint) ([]dto.CardResponse, error) {
	dtos, err := s.cardRepo.GetCardsByUserID(userID)
	if err != nil {
		return nil, errors.New("error getting card by user id")
	}

	for i := range dtos {
		decryptedNumber, err := s.decryptWithPGP(dtos[i].Number, s.encryptKey)
		if err != nil {
			accountId := strconv.Itoa(int(dtos[i].AccountID))
			decryptedNumber = "**** **** **** "
			logrus.Error("failed to decrypt card number for account: " + accountId + ", error: " + err.Error())
		}

		dtos[i].Number = decryptedNumber
	}
	return dtos, nil
}

func (s *CardService) GenerateCard(accountID uint) (*models.Card, error) {
	id, err := s.cardRepo.FindByAccountId(accountID)
	if err != nil || id != nil {
		logrus.Warn("a card has already been created for this accountId " + strconv.Itoa(int(accountID)))
		return nil, errors.New("a card has already been created for this account")
	}

	cardNumber := utils.GenerateLuhnNumber(16)
	cvv := fmt.Sprintf("%03d", rand.Intn(999)+1)
	expiry := time.Now().AddDate(0, 5, 0).AddDate(0, rand.Intn(24), 0)

	// Шифруем номер карты
	encryptedNumber, err := s.encryptWithPGP(cardNumber, s.encryptKey)
	if err != nil {
		return nil, err
	}

	// Хэшируем CVV
	hashedCVV, err := s.hashCVV(cvv)
	if err != nil {
		return nil, err
	}

	card := &models.Card{
		AccountId:  accountID,
		CardNumber: encryptedNumber,
		Cvv:        hashedCVV,
		ExpiredAt:  expiry,
	}

	err = s.cardRepo.Create(card)
	if err != nil {
		return nil, err
	}
	card.CardNumber = cardNumber
	card.Cvv = cvv

	logrus.Info("created new card for account " + strconv.Itoa(int(accountID)))
	return card, nil
}

// encryptWithPGP шифрует данные с помощью PGP через SQL-функцию
func (s *CardService) encryptWithPGP(data string, key string) (string, error) {
	var encrypted string
	query := `SELECT encode(pgp_sym_encrypt($1::text, $2::text), 'hex')`
	err := s.db.Raw(query, data, key).Scan(&encrypted).Error
	if err != nil {
		return "", err
	}
	return encrypted, nil
}

// decryptWithPGP расшифровывает данные с помощью PGP через SQL-функцию
func (s *CardService) decryptWithPGP(encryptedHex string, key string) (string, error) {
	var decrypted string
	query := `SELECT pgp_sym_decrypt(decode($1::text, 'hex'), $2::text)`
	err := s.db.Raw(query, encryptedHex, key).Scan(&decrypted).Error
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

// hashCVV создает bcrypt-хеш для CVV
func (s *CardService) hashCVV(cvv string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// validateCVV проверяет CVV против хеша
func (s *CardService) validateCVV(cvv string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(cvv))
	return err == nil
}

func (s *CardService) PayWithCard(req dto.CardPaymentRequest) (decimal.Decimal, error) {
	card, err := s.cardRepo.FindByPlainCardNumber(req.CardNumber, s.encryptKey)
	if err != nil || card == nil {
		return decimal.Zero, errors.New("card not found")
	}

	isValidCVV := s.validateCVV(req.Cvv, card.Cvv)
	if !isValidCVV {
		return decimal.Zero, errors.New("CVV is not valid")
	}

	account, err := s.accountRepo.FindByID(card.AccountId)
	if err != nil {
		return decimal.Zero, errors.New("card account not found")
	}

	withdraw, err := s.accountService.Withdraw(account.ID, account.UserID, decimal.NewFromFloat(req.Amount))
	if err != nil {
		return decimal.Zero, err
	}

	user, err := s.userRepo.FindByID(account.UserID)
	if err != nil {
		return decimal.Zero, err
	}

	transaction := &models.Transaction{
		FromAccountID:   account.ID,
		ToAccountID:     account.ID,
		Amount:          decimal.NewFromFloat(req.Amount),
		TransactionType: "payment",
		Currency:        "RUB",
	}
	err = s.db.Create(transaction).Error

	if err != nil {
		return decimal.Zero, err
	}

	notification := dto.PaymentNotification{
		To:        user.Email,
		Name:      user.Username,
		CardLast4: getLast4Digits(card.CardNumber),
		Amount:    transaction.Amount,
		Balance:   withdraw,
		Date:      transaction.CreatedAt,
	}
	err = s.mailService.SendPaymentSuccess(notification)
	if err != nil {
		logrus.Warningf("Mail not found: " + err.Error())
	}

	return withdraw, nil
}

func getLast4Digits(s string) string {
	if len(s) < 4 {
		return ""
	}
	return s[len(s)-4:]
}
