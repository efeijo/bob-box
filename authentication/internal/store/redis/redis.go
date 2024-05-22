package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	db *redis.Client
}

type Config struct {
	Port              int           `json:"port,omitempty"`
	KeyExpirationTime time.Duration `json:"key_expiration_time,omitempty"`
}

func NewClient() *Redis {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	return &Redis{db: rdb}
}
