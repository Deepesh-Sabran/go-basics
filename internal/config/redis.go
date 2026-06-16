package config

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         GetRedisAddr(),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		log.Fatal("❌ Failed to connect to Redis: ", err)
	}

	log.Println("✅ Successfully connected to Redis !!")
}