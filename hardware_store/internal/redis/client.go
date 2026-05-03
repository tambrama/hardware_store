package redis

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
	logger *slog.Logger
}

type Config struct {
	Addr     string
	Password string
	DB       int
}

func NewClient(cfg Config, logger *slog.Logger) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Error("failed to connect to redis", slog.Any("error", err))
		return nil
	}
	logger.Info("connected to redis successfully", slog.String("addr", cfg.Addr))
	return &Client{
		Client: rdb,
		logger: logger,
	}
}
func (c *Client) Close() error {
	c.logger.Info("closing redis client")
	return c.Client.Close()
}
