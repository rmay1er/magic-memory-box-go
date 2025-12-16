package memorybox

import (
	"sync"
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
	mu     sync.RWMutex         // Mutex to protect concurrent access to memory.
}

type MemoryBoxConfig struct {
	// ContextLenSize defines the size of the context length.
	ContextLenSize int

	// ExpireTime defines the expiration duration for stored memories.
	ExpireTime time.Duration
}

type Role string

const (
	SystemRole    Role = "system"
	ToolRole      Role = "tool"
	UserRole      Role = "user"
	AssistantRole Role = "assistant"
)

// Message represents a chat message with a role and content.
type Message struct {
	Role    Role   // Role of the message sender
	Content string // Content of the message
}
