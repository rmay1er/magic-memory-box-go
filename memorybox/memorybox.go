package memorybox

import (
	"context"
	"encoding/json"
	"fmt"
)

func NewMemoryBox(m IMemorizer, cfg MemoryBoxConfig) IMemoryBox {
	return &MemoryBox{
		IMemorizer:      m,
		MemoryBoxConfig: cfg,
	}
}

func (b *MemoryBox) AddRaw(ctx context.Context, userid string, role Role, value string) ([]Message, error) {
	lastMessages, err := b.Get(ctx, userid)
	if err != nil {
		fmt.Println(err)
	}
	// Инициализируем массив
	data := []Message{}

	// Распарсим, если есть старые сообщения
	if lastMessages != "" {
		if err := json.Unmarshal([]byte(lastMessages), &data); err != nil {
			return data, err
		}
	}

	// Добавляем новое сообщение
	data = append(data, Message{
		Role:    role,
		Content: value,
	})

	// Сохраняем обратно
	jsonMsgArr, err := json.Marshal(data)
	if err != nil {
		return data, err
	}
	if err := b.Set(ctx, userid, string(jsonMsgArr), b.ExpireTime); err != nil {
		return data, err
	}

	return data, nil
}

func (b *MemoryBox) Talk(ctx context.Context, userid string, value string) ([]Message, error) {
	resp, err := b.AddRaw(ctx, userid, User, value)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *MemoryBox) Remember(ctx context.Context, userid string, value string) ([]Message, error) {
	resp, err := b.AddRaw(ctx, userid, Assistant, value)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (b *MemoryBox) GetMemories(ctx context.Context, userid string) ([]Message, error) {
	lastMessages, err := b.Get(ctx, userid)
	if err != nil {
		fmt.Println(err)
	}
	// Инициализируем массив
	data := []Message{}

	// Распарсим, если есть старые сообщения
	if lastMessages != "" {
		if err := json.Unmarshal([]byte(lastMessages), &data); err != nil {
			return data, err
		}
	}

	return data, nil
}
