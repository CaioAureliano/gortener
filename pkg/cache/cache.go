package cache

import (
	"github.com/go-redis/redis/v9"
)

type Cache struct {
}

func New() *Cache {
	return &Cache{}
}

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
