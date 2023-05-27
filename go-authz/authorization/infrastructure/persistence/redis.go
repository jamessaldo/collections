package persistence

import (
	"authorization/config"
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         context.Context
)

func ConnectRedis() {
	ctx = context.TODO()

	RedisClient = redis.NewClient(&redis.Options{
		Addr: config.StorageConfig.RedisHost + ":" + config.StorageConfig.RedisPort,
	})

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err := RedisClient.Set(ctx, "connection", "test connection", 0).Err()
	if err != nil {
		log.Fatal("❌ Redis client connection failed: ", err)
	}

	fmt.Println("✅ Redis client connected successfully...")
}
