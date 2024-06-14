package connection

import (
	"backend/server/env"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

type Redis struct{}

var RedisClient *redis.Client

var redisConnected = false

func (ref Redis) Connect() (*redis.Client, context.Context, error) {
	ctx := context.Background()

	if RedisClient == nil {
		uri := env.GetRedisUrl()
		options, err := redis.ParseURL(uri)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse Redis URL: %v", err)
		}

		client := redis.NewClient(options)

		status := client.Ping(ctx)
		if status.Err() != nil {
			log.Fatalf("❌ Redis error: %s", status.Err())
			return nil, nil, fmt.Errorf("failed to connect to Redis: %v", status.Err())
		}

		if !redisConnected {
			log.Println("✅ Redis Connected")
			redisConnected = true
		}

		RedisClient = client
		return client, ctx, nil
	}

	return RedisClient, ctx, nil
}
