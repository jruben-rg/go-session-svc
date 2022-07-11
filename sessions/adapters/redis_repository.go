package adapters

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jruben-rg/go-session-svc/sessions/domain/session"
)

type redisCache struct {
	host     string
	db       int
	password string
	expires  time.Duration
	client   *redis.Client
}

func NewRedisCache(host string, db int, password string, expires time.Duration) session.Repository {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})
	return &redisCache{host, db, password, expires, client}
}

func (c *redisCache) Set(ctx context.Context, key string, value interface{}) error {

	_, err := c.client.Set(ctx, key, value, c.expires).Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *redisCache) Get(ctx context.Context, key string) (interface{}, error) {

	val, err := c.client.Get(ctx, key).Result()
	return val, err
}

func (c redisCache) Delete(ctx context.Context, key string) (int64, error) {

	val, err := c.client.Del(ctx, key).Result()
	return val, err
}
