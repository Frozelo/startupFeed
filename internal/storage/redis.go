package storage

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

var ctx = context.Background()

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal("failed to connect redis", err)
	}

	log.Printf("Connected to redis with %s", rdb.Options().Addr)

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Set(
	key string,
	value any,
	expiration time.Duration,
) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

// Get возвращает значение по ключу
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Delete удаляет ключ
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(ctx, key).Err()
}
