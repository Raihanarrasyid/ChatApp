package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisDB(redisHost string, redisPass string) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost,
		Password: redisPass,
		DB: 	 0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return rdb, nil
}