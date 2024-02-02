package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rohanraj7316/logger"
)

type Handler struct {
	prefix string

	Client *redis.Client
}

// New responsible for creating new redis client
func New(config ...Config) (handler Handler, err error) {
	err = logger.Configure()
	if err != nil {
		return handler, err
	}

	cfg, err := configDefault(config...)
	if err != nil {
		return handler, err
	}

	handler = Handler{
		Client: redis.NewClient(cfg.Redis),
		prefix: cfg.Prefix,
	}

	return handler, nil
}

func (h Handler) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	key = fmt.Sprintf("%s%s", h.prefix, key)
	return h.Client.Set(ctx, key, value, expiration)
}

func (h Handler) Get(ctx context.Context, key string) *redis.StringCmd {
	key = fmt.Sprintf("%s%s", h.prefix, key)
	return h.Client.Get(ctx, key)
}

func (h Handler) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	for index, val := range keys {
		keys[index] = fmt.Sprintf("%s%s", h.prefix, val)
	}

	return h.Client.Del(ctx, keys...)
}

func (h Handler) ExpireAt(ctx context.Context, key string, expiration time.Time) *redis.BoolCmd {
	key = fmt.Sprintf("%s%s", h.prefix, key)
	return h.Client.ExpireAt(ctx, key, expiration)
}
