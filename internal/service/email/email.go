package email

import (
	"crypto/tls"
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
	msg := []byte(fmt.Sprintf("Subject: Your OTP Code\n\nYour OTP code is: %s", otpCode))
	auth := smtp.PlainAuth("", es.Sender, es.Password, es.SMTPHost)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         es.SMTPHost,
	}

	conn, err := tls.Dial("tcp", es.SMTPHost+":"+es.SMTPPort, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial SMTP server: %w", err)
	}

	client, err := smtp.NewClient(conn, es.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err := client.Mail(es.Sender); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data stream: %w", err)
	}
	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close data stream: %w", err)
	}

	// err := smtp.SendMail(es.SMTPHost+":"+es.SMTPPort, auth, es.Sender, []string{to}, msg)
	// if err != nil {
	// 	return err
	// }

	return nil
}