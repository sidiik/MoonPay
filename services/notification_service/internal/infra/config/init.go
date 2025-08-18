package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Email         string
	EmailPassword string
	SmtpHost      string
	SmtpPort      string

	// RabbitMQ
	RabbitMQUrl string
}

var AppConfig Config

func InitConfig() {
	_ = godotenv.Load()

	AppConfig = Config{
		Email:         getEnv("EMAIL_ADDRESS"),
		EmailPassword: getEnv("EMAIL_PASSWORD"),
		RabbitMQUrl:   getEnv("RABBITMQ_URL"),
		SmtpHost:      getEnv("SMTP_HOST"),
		SmtpPort:      getEnv("SMTP_PORT"),
	}

}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		slog.Error("missing required environment variable", slog.String("key", key))
		os.Exit(1)
	}
	return val
}
