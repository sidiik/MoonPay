package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {

	// JWT
	AccessTokenSecret  string
	RefreshTokenSecret string

	// REDIS
	RedisAddr     string
	RedisPassword string
	RedisDb       string
}

var AppConfig Config

func InitConfig() {
	_ = godotenv.Load()

	AppConfig = Config{

		AccessTokenSecret:  getEnv("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: getEnv("REFRESH_TOKEN_SECRET"),

		RedisAddr:     getEnv("REDIS_ADDR"),
		RedisPassword: getEnv("REDIS_PASSWORD"),
		RedisDb:       getEnv("REDIS_DB"),
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
