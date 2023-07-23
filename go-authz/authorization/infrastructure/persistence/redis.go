package persistence

import (
	"authorization/config"
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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
		log.Fatal().Caller().Err(err).Msg("❌ Redis client is not connected")
	}

	err := RedisClient.Set(ctx, "connection", "test connection", 0).Err()
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("❌ Redis client connection failed")
	}

	log.Info().Caller().Msg("✅ Redis client connected successfully...")
}
