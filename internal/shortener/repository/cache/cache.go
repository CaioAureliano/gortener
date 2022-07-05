package cache

import (
	"context"
	"time"

	cacheRedis "github.com/CaioAureliano/gortener/pkg/cache"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Set(key, value string, duration time.Duration) error
	Get(key string) (string, error)
}

type cache struct {
	client *redis.Client
	ctx    context.Context
}

func New() Cache {
	return &cache{
		client: cacheRedis.Client(),
		ctx:    context.Background(),
	}
}

func (c cache) Set(key, value string, duration time.Duration) error {
	return c.client.Set(c.ctx, key, value, duration).Err()
}

func (c cache) Get(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}
