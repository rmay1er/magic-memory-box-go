package memorybox

import (
	"context"
	"encoding/json"
	"fmt"
)

// NewMemoryBox creates a new MemoryBox instance with the given IMemorizer and configuration.
func NewMemoryBox(m IMemorizer, cfg MemoryBoxConfig) IMemoryBox {
	return &MemoryBox{
		IMemorizer:      m,
		MemoryBoxConfig: cfg,
	}
}

// AddRaw retrieves the existing messages for a user, appends a new message with the specified role and content,
// and saves the updated list back to the memory store.
func (b *MemoryBox) AddRaw(ctx context.Context, userid string, role Role, value string) ([]Message, error) {
	lastMessages, err := b.Get(ctx, userid)
	if err != nil {
		fmt.Println(err)
	}
	// Initialize the array
	data := []Message{}

	// Parse existing messages if any
	if lastMessages != "" {
		if err := json.Unmarshal([]byte(lastMessages), &data); err != nil {
			return data, err
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
		fmt.Println(err)
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
