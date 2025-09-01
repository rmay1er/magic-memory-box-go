package memorybox

import (
	"context"
	"time"
)

type MemoryBoxConfig struct {
	// ContextLenSize defines the size of the context length.
	ContextLenSize int

	// ExpireTime defines the expiration duration for stored memories.
	ExpireTime time.Duration
}

type MemoryBox struct {
	IMemorizer
	MemoryBoxConfig
}

// IMemorizer is the base interface for working with Redis.
type IMemorizer interface {
	// Set sets a value with an optional expiration time (TTL).
	Set(ctx context.Context, key string, value any, expiration ...time.Duration) error

	// Get returns the value of the given key.
	// If the key does not exist, redis.Nil is returned.
	Get(ctx context.Context, key string) (string, error)
}

type Role string

const (
	System    Role = "system"
	Tool      Role = "tool"
	User      Role = "user"
	Assistant Role = "assistant"
)

// Message represents a chat message with a role and content.
type Message struct {
	Role    Role   // Role of the message sender
	Content string // Content of the message
}

type IMemoryBox interface {
	AddRaw(ctx context.Context, userid string, msgType Role, value string) ([]Message, error)
	Talk(ctx context.Context, userid string, value string) ([]Message, error)
	Remember(ctx context.Context, userid string, value string) ([]Message, error)
	GetMemories(ctx context.Context, userid string) ([]Message, error)
}
