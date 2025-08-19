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
	case "user.otp.requestPasswordReset":
		return s.HandlePasswordResetRequest(event)
	case "user.otp.passwordReset":
		return s.HandlePasswordReset(event)
	default:
		fmt.Printf("⚠️ Unknown event: %+v\n", event)
	}

	return nil
}

func (s *AuthNotificationService) HandleUserRegisteredEvent(event map[string]any) error {
	slog.Info("User registered event triggered")
	slog.Info("Sending welcome message")

	subject := "🎉 Welcome to MoonPay, " + event["fullName"].(string) + "!"

	body := "Hi " + event["fullName"].(string) + ",<br><br>" +
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

func (s *AuthNotificationService) HandlePasswordResetRequest(event map[string]any) error {
	slog.Info("Password reset request event triggered")
	slog.Info("Sending otp message")

	subject := "Password Reset Request #" + event["code"].(string)

	body := "Hi " + event["fullName"].(string) + ",<br><br>" +
		"We received a request to reset your password for your <b>MoonPay</b> account.<br>" +
		"Your One-Time Password (OTP) is:<br><br>" +
		"<b>" + event["code"].(string) + "</b><br><br>" +
		"This OTP is valid for " + event["expiresIn"].(string) + " Minutes, so please use it soon.<br><br>" +
		"If you did not request a password reset, please ignore this email.<br><br>" +
		"Cheers,<br>" +
		"The MoonPay Team"

	if err := s.emailService.Send(event["email"].(string), subject, body); err != nil {
		slog.Info("failed to send email", "error", err)
		return err
	}

	return nil

}

func (s *AuthNotificationService) HandlePasswordReset(event map[string]any) error {
	slog.Info("Password reset event triggered")
	slog.Info("Sending password reset notification message")

	subject := "Password Reset - MoonPay"

	body := "Hi " + event["fullName"].(string) + ",<br><br>" +
		"Your password for your <b>MoonPay</b> account has been successfully reset.<br>" +
		"If you initiated this change, you can safely continue using your account.<br><br>" +
		"Location of the reset attempt: <b>" + event["location"].(string) + "</b><br><br>" +
		"If you did not reset your password, please contact our support immediately.<br><br>" +
		"Cheers,<br>" +
		"The MoonPay Team"

	if err := s.emailService.Send(event["email"].(string), subject, body); err != nil {
		slog.Info("failed to send email", "error", err)
		return err
	}

	return nil

}
