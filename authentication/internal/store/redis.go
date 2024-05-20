package store

import (
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	db *redis.Client
}

func NewRedisClient() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	return &Redis{db: rdb}
}
