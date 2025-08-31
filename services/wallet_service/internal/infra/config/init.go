package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	// Write only Databases
	DBWriterName     string
	DBWriterUser     string
	DBWriterPassword string
	DBWriterHost     string
	DBWriterPort     string

	// Read only databases
	DBReaderName     string
	DBReaderUser     string
	DBReaderPassword string
	DBReaderHost     string
	DBReaderPort     string

	// JWT
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenExpire  string
	RefreshTokenExpire string

	// RabbitMQ
	RabbitMQUrl string

	// OTP
	OtpCodeExpire string

	// REDIS
	RedisAddr     string
	RedisPassword string
	RedisDb       string
}

var AppConfig Config

func InitConfig() {
	_ = godotenv.Load()

	AppConfig = Config{
		Port: getEnv("PORT"),
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
