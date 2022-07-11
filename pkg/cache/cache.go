package cache

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	addr = os.Getenv("REDIS_URI")
)

func Client() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
