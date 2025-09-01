package memorybox

import (
	"context"
	"time"
)

type MemoryBoxConfig struct {
	ContextLenSize int
	ExpireTime     time.Duration
}

type MemoryBox struct {
	IMemorizer
	MemoryBoxConfig
}

// Redis — базовый интерфейс для работы с Redis.
type IMemorizer interface {
	// Set устанавливает значение с возможным временем жизни (TTL).
	Set(ctx context.Context, key string, value any, expiration ...time.Duration) error

	// Get возвращает значение по ключу.
	// Если ключа нет, возвращается redis.Nil.
	Get(ctx context.Context, key string) (string, error)

	// Del удаляет один или несколько ключей.
	// Del(ctx context.Context, keys ...string) (int64, error)

	// Exists проверяет существование ключей.
	// Exists(ctx context.Context, keys ...string) (int64, error)

	// Expire устанавливает TTL для ключа.
	// Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)

	// TTL возвращает оставшееся время жизни ключа.
	// TTL(ctx context.Context, key string) (time.Duration, error)

	// Keys возвращает список ключей по шаблону (например, "session:*").
	// Keys(ctx context.Context, pattern string) ([]string, error)
}

type Role string

const (
	System    Role = "system"
	Tool      Role = "tool"
	User      Role = "user"
	Assistant Role = "assistant"
)

type Message struct {
	Role    Role
	Content string
}

type IMemoryBox interface {
	AddRaw(ctx context.Context, userid string, msgType Role, value string) ([]Message, error)
	Talk(ctx context.Context, userid string, value string) ([]Message, error)
	Remember(ctx context.Context, userid string, value string) ([]Message, error)
	GetMemories(ctx context.Context, userid string) ([]Message, error)
}
