package rdb

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

func NewRedisAdapter(client *redis.Client, prefix string, flushSessionOnSIGINT ...bool) *RedisAdapter {
	adapter := &RedisAdapter{client: client, prefix: prefix}

	if len(flushSessionOnSIGINT) > 0 && flushSessionOnSIGINT[0] {
		adapter.flushSessionOnSIGINT = true

		// Ловим Ctrl+C
		go func() {
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			<-c
			fmt.Println("Clearing redis session before exit...")

			if err := adapter.ClearPrefix(context.Background()); err != nil {
				fmt.Printf("failed to clear redis: %v\n", err)
			}
			os.Exit(0)
		}()
	}

	return adapter
}

// Set сохраняет значение с TTL (если указан)
func (r *RedisAdapter) Set(ctx context.Context, key string, value any, expiration ...time.Duration) error {
	var exp time.Duration
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	return r.client.Set(ctx, r.prefix+key, value, exp).Err()
}

// Get получает значение по ключу
func (r *RedisAdapter) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, r.prefix+key).Result()
}

// ClearPrefix удаляет все ключи по текущему префиксу
func (r *RedisAdapter) ClearPrefix(ctx context.Context) error {
	iter := r.client.Scan(ctx, 0, r.prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}
	return iter.Err()
}
