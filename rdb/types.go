package rdb

import (
	"github.com/go-redis/redis/v8"
)

type RedisAdapter struct {
	client               *redis.Client
	flushSessionOnSIGINT bool
	prefix               string
}
