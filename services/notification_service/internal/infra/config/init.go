package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string

	EmailUsername string
	EmailPassword string
	SmtpHost      string
	SmtpPort      string

	RabbitMQUrl string
}

var AppConfig Config

func InitConfig() {
	_ = godotenv.Load()

	AppConfig = Config{
		AppName:       getEnv("APP_NAME"),
		EmailUsername: getEnv("EMAIL_ADDRESS"),
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
