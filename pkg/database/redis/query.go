package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type IRedisHash interface {
	HGet(ctx context.Context, key, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd
	HExists(ctx context.Context, key, field string) *redis.BoolCmd
}

type IRedisGeneric interface {
	Keys(ctx context.Context, pattern string) *redis.StringSliceCmd
}

type IPubSub interface {
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
}

type IRedisClient interface {
	IRedisHash
	IRedisGeneric
	IPubSub
}

func (c *Client) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	return c.rCleint.HGet(ctx, key, field)
}

func (c *Client) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	return c.rCleint.HGetAll(ctx, key)
}

func (c *Client) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	return c.rCleint.HSet(ctx, key, values)
}

func (c *Client) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	return c.rCleint.HMSet(ctx, key, values)
}

func (c *Client) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	return c.rCleint.HExists(ctx, key, field)
}

func (c *Client) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	return c.rCleint.Keys(ctx, pattern)
}

func (c *Client) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	return c.rCleint.Publish(ctx, channel, message)
}
