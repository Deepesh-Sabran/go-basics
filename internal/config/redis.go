package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err:= RedisClient.Ping(Ctx).Err(); err != nil {
		log.Println("Failed to connect to Redis")
		panic(err)
	}
}