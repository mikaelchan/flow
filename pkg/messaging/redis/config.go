package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Client        *redis.Client
	HandleTimeout time.Duration
}
