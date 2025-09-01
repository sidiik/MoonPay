package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	MongoDBURI string

	// RabbitMQ
	RabbitMQUrl string

	// REDIS
	RedisAddr     string
	RedisPassword string
	RedisDb       string
}

var AppConfig Config

func InitConfig() {
	_ = godotenv.Load()

	AppConfig = Config{
		Port:       getEnv("PORT"),
		MongoDBURI: getEnv("MONGODB_URI"),
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
