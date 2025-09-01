package incache

import (
	"time"
)

// mapFields holds a cached value along with its expiration time.
type mapFields struct {
	Value      any       // The stored value in the cache.
	ExpireTime time.Time // The time when the value expires.
}

// inCache represents an in-memory cache using a map.
type inCache struct {
	memory map[string]mapFields // The cache storage mapping keys to cached values.
}
