package main

import (
	"log/slog"
	"os"

	"github.com/sidiik/notification-service/internal/infra/config"
	"github.com/sidiik/notification-service/internal/infra/rabbitmq"
	"github.com/sidiik/notification-service/internal/services"
)

func main() {
	slog.Info("Initializing app config")
	config.InitConfig()
	appConfig := config.AppConfig

	slog.Info("Initializing email service")
	emailService := services.NewEmailService(appConfig.SmtpHost, appConfig.SmtpPort, appConfig.Email, appConfig.EmailPassword)

	slog.Info("initializing auth_consumer")
	authConsumer, err := rabbitmq.NewConsumer(appConfig.RabbitMQUrl, "notifications", "auth_exchange")
	if err != nil {
		slog.Error("failed to initialize auth consumer", "error", err)
		os.Exit(1)
	}

	defer authConsumer.Conn.Close()
	defer authConsumer.Channel.Close()

	authNotificationService := services.NewAuthNotificationService(emailService)

	if err := authConsumer.Start(authNotificationService.HandleEvent); err != nil {
		slog.Error("failed to start auth notification service", "error", err)
		os.Exit(1)
	}

}
