package services

import (
	"log/slog"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	smtpHost string
	smtpPort int
	username string
	password string
}

func NewEmailService(smtpHost, smtpPort, username, password string) *EmailService {
	smtpPortInt, err := strconv.ParseInt(smtpPort, 10, 64)
	if err != nil {
		slog.Error("failed to parse smtpPort")
		os.Exit(1)
	}
	return &EmailService{
		smtpHost: smtpHost,
		smtpPort: int(smtpPortInt),
		username: username,
		password: password,
	}
}

func (s *EmailService) Send(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.username, s.password)

	if err := d.DialAndSend(m); err != nil {
		slog.Error("failed to send email: %v", "error", err)
		return err
	}

	slog.Info("Message sent successfully")

	return nil

}
