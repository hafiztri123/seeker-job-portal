package redis

import (
	"github.com/hafiztri123/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Addr(),
		Password: cfg.Password,
		DB: cfg.DB,
	})
}