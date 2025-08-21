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
		Port:          getEnv("PORT"),
		OtpCodeExpire: getEnv("OTP_CODE_EXPIRE"),

		DBWriterName:     getEnv("DB_WRITER_NAME"),
		DBWriterUser:     getEnv("DB_WRITER_USER"),
		DBWriterPassword: getEnv("DB_WRITER_PASSWORD"),
		DBWriterHost:     getEnv("DB_WRITER_HOST"),
		DBWriterPort:     getEnv("DB_WRITER_PORT"),

		DBReaderName:     getEnv("DB_READER_NAME"),
		DBReaderUser:     getEnv("DB_READER_USER"),
		DBReaderPassword: getEnv("DB_READER_PASSWORD"),
		DBReaderHost:     getEnv("DB_READER_HOST"),
		DBReaderPort:     getEnv("DB_READER_PORT"),

		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET"),
		AccessTokenExpire:  getEnv("ACCESS_TOKEN_EXPIRE"),
		RefreshTokenExpire: getEnv("REFRESH_TOKEN_EXPIRE"),

		RabbitMQUrl: getEnv("RABBITMQ_URL"),
		RedisAddr:   getEnv("REDIS_ADDR"),
		// RedisPassword: getEnv("REDIS_PASSWORD"),
		RedisDb: getEnv("REDIS_DB"),
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
