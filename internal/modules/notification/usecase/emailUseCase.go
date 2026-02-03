package utils

import (
	"fmt"

	"github.com/NguyenNhatquang522004/BE-modular-golang/internal/common/configs"
	"gopkg.in/gomail.v2"
)

type IEmailService interface {
	SendEmail(to string, subject string, body string) error
	SendOTP(toEmail string, otpCode string) error
}

type EmailService struct {
	config configs.EmailSMTPConfig
	dialer *gomail.Dialer
}

// Constructor
func NewEmailService(cfg *configs.Config) *EmailService {
	// Khởi tạo Dialer 1 lần để tái sử dụng
	d := gomail.NewDialer(
		cfg.EmailSMTP.SMTP_HOST,
		cfg.EmailSMTP.SMTP_PORT,
		cfg.EmailSMTP.SMTP_EMAIL,
		cfg.EmailSMTP.SMTP_PASSWORD,
	)

	return &EmailService{
		config: cfg.EmailSMTP,
		dialer: d,
	}
}

// Hàm gửi Email cơ bản
func (s *EmailService) SendEmail(to string, subject string, body string) error {
	m := gomail.NewMessage()

	// Header: From, To, Subject
	// Format: "Tên Hiển Thị <email@gmail.com>"
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.config.SMTP_SENDER_NAME, s.config.SMTP_EMAIL))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Body: text/html để hỗ trợ định dạng đẹp, hoặc text/plain
	m.SetBody("text/html", body)

	// Thực hiện gửi
	if err := s.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// Hàm gửi OTP (Tiện ích cụ thể)
func (s *EmailService) SendOTP(toEmail string, otpCode string) error {
	subject := "Mã xác thực OTP của bạn"
	// Nội dung HTML đơn giản
	body := fmt.Sprintf(`
        <div style="font-family: Helvetica,Arial,sans-serif;min-width:1000px;overflow:auto;line-height:2">
            <div style="margin:50px auto;width:70%;padding:20px 0">
                <div style="border-bottom:1px solid #eee">
                    <a href="" style="font-size:1.4em;color: #00466a;text-decoration:none;font-weight:600">My App</a>
                </div>
                <p style="font-size:1.1em">Xin chào,</p>
                <p>Đây là mã xác thực OTP của bạn. Mã này sẽ hết hạn trong 15 phút.</p>
                <h2 style="background: #00466a;margin: 0 auto;width: max-content;padding: 0 10px;color: #fff;border-radius: 4px;">%s</h2>
                <p style="font-size:0.9em;">Xin cảm ơn,<br />My App Team</p>
            </div>
        </div>
    `, otpCode)

	return s.SendEmail(toEmail, subject, body)
}
