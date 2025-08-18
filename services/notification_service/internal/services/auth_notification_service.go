package services

import (
	"fmt"
	"log/slog"
)

type AuthNotificationService struct {
	emailService EmailService
}

func NewAuthNotificationService(emailService *EmailService) *AuthNotificationService {
	return &AuthNotificationService{
		emailService: *emailService,
	}
}

func (s *AuthNotificationService) HandleEvent(event map[string]any) error {
	switch event["event"] {
	case "user.login":
		slog.Info("Login triggered by user", "email", event["email"])
	case "user.registered":
		return s.HandleUserRegisteredEvent(event)
	default:
		fmt.Printf("⚠️ Unknown event: %+v\n", event)
	}

	return nil
}

func (s *AuthNotificationService) HandleUserRegisteredEvent(event map[string]any) error {
	slog.Info("User registered event triggered")
	slog.Info("Sending welcome message")

	subject := "🎉 Welcome to MoonPay, " + event["Name"].(string) + "!"

	body := "Hi " + event["Name"].(string) + ",<br><br>" +
		"Welcome to <b>MoonPay</b>! 🚀<br>" +
		"We're excited to have you join us and can’t wait for you to start exploring.<br><br>" +
		"If you have any questions, our team is always here to help.<br><br>" +
		"Cheers,<br>" +
		"The MoonPay Team"

	if err := s.emailService.Send(event["email"].(string), subject, body); err != nil {
		slog.Info("failed to send email", "error", err)
		return err
	}

	return nil

}
