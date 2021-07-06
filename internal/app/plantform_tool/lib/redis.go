package lib

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"ls/internal/app/plantform_tool/clients"
)

func InitRedis() {
	clients.RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", "127.0.0.1", 6379),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}