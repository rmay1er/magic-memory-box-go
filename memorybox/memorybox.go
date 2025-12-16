package memorybox

import (
	"context"
	"encoding/json"
	"time"
)

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
func (b *MemoryBox) Talk(ctx context.Context, userid string, value string) ([]Message, error) {
	resp, err := b.AddRaw(ctx, userid, User, value)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Remember adds an assistant message to the memory for the specified user.
func (b *MemoryBox) Remember(ctx context.Context, userid string, value string) ([]Message, error) {
	resp, err := b.AddRaw(ctx, userid, Assistant, value)
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
