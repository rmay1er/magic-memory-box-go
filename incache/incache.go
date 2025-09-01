package incache

import (
	"context"
	"fmt"
	"time"
)

// NewCache creates and returns a pointer to a new inCache instance.
// The cache uses an internal map to store keys and their associated values along with expiration metadata.
func NewCache() *inCache {
	return &inCache{
		memory: make(map[string]mapFields),
	}
}

// Set stores a value in the cache with an optional expiration time.
// If an expiration duration is provided and greater than zero, the key will expire after that duration.
// If no expiration is specified, the value is stored indefinitely without expiration.
// The value is converted to a string representation before being stored.
// Context parameter is accepted for future extensibility but currently not used.
func (c *inCache) Set(ctx context.Context, key string, value any, expiration ...time.Duration) error {
	var expireTime time.Time
	if len(expiration) > 0 && expiration[0] > 0 {
		expireTime = time.Now().Add(expiration[0])
	}

	c.memory[key] = mapFields{
		Value:      fmt.Sprintf("%v", value),
		ExpireTime: expireTime, // zero time means no TTL (time-to-live)
	}
	return nil
}

// Get retrieves a value from the cache by key.
// If the key is not found, it returns an error indicating so.
// If the key exists but has expired, the key is deleted from the cache and an expiration error is returned.
// Otherwise, the cached string value is returned.
// Context parameter is accepted for future extensibility but currently not used.
func (c *inCache) Get(ctx context.Context, key string) (string, error) {
	mf, ok := c.memory[key]
	if !ok {
		return "", fmt.Errorf("key not found")
	}

	// Check if the key has an expiration time and if it has passed
	if !mf.ExpireTime.IsZero() && time.Now().After(mf.ExpireTime) {
		delete(c.memory, key) // Remove expired key
		return "", fmt.Errorf("key expired")
	}

	return fmt.Sprintf("%v", mf.Value), nil
}
