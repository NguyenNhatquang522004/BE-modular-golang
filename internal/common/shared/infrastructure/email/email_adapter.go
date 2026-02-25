package email

import (
	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"gopkg.in/gomail.v2"
)

type EmailAdapter struct {
	config configs.EmailSMTPConfig
	dialer *gomail.Dialer
}

func NewEmailAdapter(cfg *configs.Config) *EmailAdapter {
	d := gomail.NewDialer(
		cfg.EmailSMTP.SMTP_HOST,
		cfg.EmailSMTP.SMTP_PORT,
		cfg.EmailSMTP.SMTP_EMAIL,
		cfg.EmailSMTP.SMTP_PASSWORD,
	)
	return &EmailAdapter{
		config: cfg.EmailSMTP,
		dialer: d,
	}
}

// SendEmail(to string, subject string, body string) error
func (e *EmailAdapter) SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.config.SMTP_EMAIL)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return e.dialer.DialAndSend(m)
}
