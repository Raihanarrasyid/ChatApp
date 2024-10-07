package email

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	SMTPHost string
	SMTPPort string
	Sender   string
	Password string
}

func NewEmailService(SMTPHost, SMTPPort, sender, password string) *EmailService {
	return &EmailService{SMTPHost, SMTPPort, sender, password}
}

func (es *EmailService) SendOTP(to string, otpCode string) error {
	auth := smtp.PlainAuth("", es.Sender, es.Password, es.SMTPHost)

	msg := []byte(fmt.Sprintf("Subject: Your OTP Code\n\nYour OTP code is: %s", otpCode))

	err := smtp.SendMail(es.SMTPHost+":"+es.SMTPPort, auth, es.Sender, []string{to}, msg)
	if err != nil {
		return err
	}

	return nil
}