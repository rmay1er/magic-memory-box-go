package memorybox

import (
	"sync"
	"time"
)

// MapFields holds a cached value along with its expiration time.
type MapFields struct {
	Value      any       // The stored value in the cache.
	ExpireTime time.Time // The time when the value expires.
}

// InCache represents an in-memory cache using a map.
type MemoryCache struct {
	memory map[string]MapFields // The cache storage mapping keys to cached values.
	mu     sync.RWMutex         // Mutex to protect concurrent access to memory.
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
