package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/mjcode-max/TurboGin/config"
	"github.com/mjcode-max/TurboGin/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	cli *redis.Client
	cfg *config.RedisConfig
	log *logger.Logger
}

// New 创建Redis客户端
func New(cfg *config.Config, log *logger.Logger) (*Client, error) {
	if !cfg.Redis.Enabled {
		return nil, nil
	}
	cli := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读超时
		WriteTimeout: 3 * time.Second, // 写超时
	})

	// 健康检查
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if _, err := cli.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	log.Info("Redis connected successfully",
		logger.String("addr", cfg.Redis.Addr),
		logger.Int("db", cfg.Redis.DB))

	return &Client{
		cli: cli,
		cfg: &cfg.Redis,
		log: log,
	}, nil
}

// Close 安全关闭连接
func (c *Client) Close() error {
	if err := c.cli.Close(); err != nil {
		return fmt.Errorf("redis close error: %w", err)
	}
	c.log.Info("Redis connection closed")
	return nil
}

// GetClient 获取原生客户端（供特殊操作使用）
func (c *Client) GetClient() *redis.Client {
	return c.cli
}

// HealthCheck 健康检查（供/health端点使用）
func (c *Client) HealthCheck() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := c.cli.Ping(ctx).Result()
	return err == nil
}
