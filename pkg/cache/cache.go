package cache

import (
	"github.com/go-redis/redis/v9"
)

var (
	addr = "localhost:6379"
)

func Client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}
