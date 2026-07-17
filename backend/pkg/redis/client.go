package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/devops-command-center/backend/config"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Client wraps go-redis.
type Client struct {
	rdb *goredis.Client
	log *zap.Logger
}

func New(cfg config.RedisConfig, log *zap.Logger) (*Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	log.Info("redis connected", zap.String("addr", cfg.Addr()))
	return &Client{rdb: rdb, log: log}, nil
}

func (c *Client) Raw() *goredis.Client { return c.rdb }

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *Client) Close() error {
	return c.rdb.Close()
}
