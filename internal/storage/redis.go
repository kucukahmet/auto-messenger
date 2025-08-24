package storage

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

func InitRedis(redisUri string) (*redis.Client, error) {
	redsidb := redis.NewClient(&redis.Options{
		Addr:     redisUri,
		Password: "",
		DB:       0,
	})
	defer redsidb.Close()

	ctx := context.Background()
	if err := redsidb.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		redsidb = nil
	}

	return redsidb, nil
}
