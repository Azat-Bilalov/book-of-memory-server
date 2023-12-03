package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Azat-Bilalov/book-of-memory-server/internal/app/config"
	"github.com/go-redis/redis/v8"
)

const servicePrefix = "book-of-memory."

type Client struct {
	cfg    config.RedisConfig
	client *redis.Client
}

func New(ctx context.Context, cfg config.RedisConfig) (*Client, error) {
	client := &Client{}

	client.cfg = cfg

	redisClient := redis.NewClient(&redis.Options{
		Password:    cfg.Password,
		Username:    cfg.User,
		Addr:        cfg.Host + ":" + strconv.Itoa(cfg.Port),
		DB:          0,
		DialTimeout: time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Second,
	})

	client.client = redisClient

	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("cant ping redis: %w", err)
	}

	return client, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}
