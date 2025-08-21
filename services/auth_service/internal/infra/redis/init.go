package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sidiik/moonpay/auth_service/internal/infra/config"
)

type RedisClient struct {
	rdb *redis.Client
}

func InitClient(config config.Config) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
		// Password: config.RedisPassword,
		DB: 0, // FIXME: Use the db from config
	})

	return &RedisClient{
		rdb: rdb,
	}
}

func (r *RedisClient) Set(ctx context.Context, exp time.Duration, key, value string) error {
	return r.rdb.Set(ctx, key, value, exp).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
