package no_sql_database

import (
	"context"
	"fmt"

	"github.com/ivan/storage-project-back/pkg/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	cfg    *config.Config
}

func NewRedisClient(cfg *config.Config) *RedisClient {
	addr := fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.RedisPassword,
	})

	return &RedisClient{client: client, cfg: cfg}
}

func (r *RedisClient) Ping(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()
	return err
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value any) error {
	return r.client.Set(ctx, key, value, r.cfg.CacheLifetime).Err()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}
