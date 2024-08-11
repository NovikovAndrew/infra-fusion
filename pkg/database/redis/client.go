package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ClientConfig struct {
	DB                                              int
	Addr, Username, Password                        string
	ReadTimeout, WriteTimeout, DialTimeout, Timeout time.Duration
}

type Client struct {
	rCleint *redis.Client
	conf    *ClientConfig
}

func NewRedisClient(ctx context.Context, conf *ClientConfig) (IRedisClient, error) {
	rCleint := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Username:     conf.Username,
		Password:     conf.Password,
		DB:           conf.DB,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		DialTimeout:  conf.DialTimeout,
	})

	ctx, cancel := context.WithTimeout(ctx, conf.Timeout)
	defer cancel()

	if err := rCleint.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{
		rCleint: rCleint,
		conf:    conf,
	}, nil
}

func (c *Client) Close() error {
	return c.rCleint.Close()
}
