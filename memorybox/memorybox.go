package memorybox

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
)

type IMemoryBox interface {
	AddRaw(ctx context.Context, userid string, msgType Role, value string) ([]Message, error)
	Tell(ctx context.Context, userid string, value string) ([]Message, error)
	TellUnsafe(ctx context.Context, userid string, value string) []Message
	Remember(ctx context.Context, userid string, value string) ([]Message, error)
	GetMemories(ctx context.Context, userid string) ([]Message, error)
}

// IMemorizer is the base interface for working with Redis.
type IMemorizer interface {
	// Set sets a value with an optional expiration time (TTL).
	Set(ctx context.Context, key string, value any, expiration ...time.Duration) error

	// Get returns the value of the given key.
	// If the key does not exist, redis.Nil is returned.
	Get(ctx context.Context, key string) (string, error)
}

type MemoryBox struct {
	IMemorizer
	MemoryBoxConfig
}

// NewMemoryBox creates a new MemoryBox instance with the given IMemorizer and configuration.
func NewMemoryBox(m IMemorizer, cfg MemoryBoxConfig) IMemoryBox {
	return &MemoryBox{
		IMemorizer:      m,
		MemoryBoxConfig: cfg,
	}
}

// WithDefault creates a new MemoryBox instance with default.
func NewMemoryBoxDefault() IMemoryBox {
	cache := NewCache()
	return &MemoryBox{
		IMemorizer: cache,
		MemoryBoxConfig: MemoryBoxConfig{
			ContextLenSize: 20,
			ExpireTime:     1 * time.Hour,
		},
	}
}

// AddRaw retrieves the existing messages for a user, appends a new message with the specified role and content,
// and saves the updated list back to the memory store.
func (b *MemoryBox) AddRaw(ctx context.Context, userid string, role Role, value string) ([]Message, error) {
	// Initialize the array
	lastMessages, _ := b.Get(ctx, userid)

	data := []Message{}

	// Parse existing messages if any
	if lastMessages != "" {
		if err := json.Unmarshal([]byte(lastMessages), &data); err != nil {
			return data, err
		}
	}

	if len(data) > b.ContextLenSize {
		if data[0].Role == "system" {
			// Удаляем второй элемент, сдвигая срез так, чтобы исключить data[1]
			data = append(data[:1], data[2:]...)
		} else {
			// Если первый элемент не system, просто удаляем первый
			data = data[1:]
		}
	}

	// Add the new message
	data = append(data, Message{
		Role:    role,
		Content: value,
	})

	// Save the updated messages back
	jsonMsgArr, err := json.Marshal(data)
	if err != nil {
		return data, err
	}
	if err := b.Set(ctx, userid, string(jsonMsgArr), b.ExpireTime); err != nil {
		return data, err
	}

	return data, nil
}

// Talk adds a user message to the memory for the specified user.
func (b *MemoryBox) Tell(ctx context.Context, userid string, value string) ([]Message, error) {
	resp, err := b.AddRaw(ctx, userid, UserRole, value)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// TalkUnsafe adds a user message to the memory for the specified user, without ERROR
// for inline pattern design.
func (b *MemoryBox) TellUnsafe(ctx context.Context, userid string, value string) []Message {
	msgs, err := b.Tell(ctx, userid, value)
	if err != nil {
		// Логируй или panic (но panic — плохо для продакшена)
		slog.Error("Ошибка в Tell", "err", err)
		return []Message{} // Или верни nil
	}
	return msgs
}

// Remember adds an assistant message to the memory for the specified user.
func (b *MemoryBox) Remember(ctx context.Context, userid string, value string) ([]Message, error) {
	resp, err := b.AddRaw(ctx, userid, AssistantRole, value)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMemories retrieves all stored messages for the specified user.
func (b *MemoryBox) GetMemories(ctx context.Context, userid string) ([]Message, error) {
	lastMessages, err := b.Get(ctx, userid)
	if err != nil {
		return []Message{}, err
	}
	// Initialize the array
	data := []Message{}

	// Parse existing messages if any
	if lastMessages != "" {
		if err := json.Unmarshal([]byte(lastMessages), &data); err != nil {
			return data, err
		}
	}

	return data, nil
}

// ConvertMessagesForReplicate converts a slice of Message structs into a slice of maps with keys "role" and "content",
// suitable for use with the Replicate API.
func ConvertMessagesForReplicate(msgs []Message) []map[string]any {
	out := make([]map[string]any, len(msgs))
	for i, m := range msgs {
		newMap := make(map[string]any)
		newMap["role"] = string(m.Role)
		newMap["content"] = m.Content
		out[i] = newMap
	}
	return out
}
