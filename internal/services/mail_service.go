package services

import (
	"BankSystem/internal/dto"
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sirupsen/logrus"
	"time"
)

type MailService struct {
	mg     mailgun.Mailgun
	domain string
	log    *logrus.Logger
}

func NewMailService(apiKey, domain string, log *logrus.Logger) *MailService {
	return &MailService{
		mg:     mailgun.NewMailgun(domain, apiKey),
		domain: domain,
		log:    log,
	}
}

func (s *MailService) SendPaymentSuccess(data dto.PaymentNotification) error {
	message := s.mg.NewMessage(
		"noreply@yourbank.com",
		"Платёж успешно выполнен",
		s.buildHTML(data),
		data.To,
	)

	logrus.Info("Send email: ", s.buildHTML(data))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := s.mg.Send(ctx, message)
	if err != nil {
		logrus.WithError(err).WithField("email", data.To).Error("Failed to send email")
		return err
	}

	return nil
}

func (s *MailService) buildHTML(data dto.PaymentNotification) string {
	return `
        <h2>Платёж успешно выполнен</h2>
        <p>Здравствуйте, ` + data.Name + `</p>
        <p>С карты <strong>**** **** **** ` + data.CardLast4 + `</strong> списано <strong>` + data.Amount.StringFixed(2) + ` RUB</strong></p>
        <p>Новый баланс: ` + data.Balance.StringFixed(2) + ` RUB</p>
        <p>Дата: ` + data.Date.Format("02.01.2006 15:04") + `</p>
        <hr/>
        <p><small>© BankSystem - Ваш банк доверяет Go</small></p>
    `
}
