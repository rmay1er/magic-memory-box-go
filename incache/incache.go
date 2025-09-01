package incache

import (
	"context"
	"fmt"
	"time"
)

func NewCache() *inCache {
	return &inCache{
		memory: make(map[string]mapFields),
	}
}

// Set stores a value in the cache with an optional expiration time.
// If an expiration duration is provided and greater than zero, the key will expire after that duration.
// If no expiration is specified, the value is stored indefinitely.
func (c *inCache) Set(ctx context.Context, key string, value any, expiration ...time.Duration) error {
	var expireTime time.Time
	if len(expiration) > 0 && expiration[0] > 0 {
		expireTime = time.Now().Add(expiration[0])
	}

	c.memory[key] = mapFields{
		Value:      fmt.Sprintf("%v", value),
		ExpireTime: expireTime, // zero time means no TTL
	}
	return nil
}

// Get retrieves a value from the cache by key.
// It returns an error if the key does not exist or if the key has expired,
// in which case the key is also removed from the cache.
func (c *inCache) Get(ctx context.Context, key string) (string, error) {
	mf, ok := c.memory[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}

	// Check if the key has a TTL and if it has expired
	if !mf.ExpireTime.IsZero() && time.Now().After(mf.ExpireTime) {
		delete(c.memory, key)
		return "", fmt.Errorf("key expired")
	}

	return fmt.Sprintf("%v", mf.Value), nil
}
