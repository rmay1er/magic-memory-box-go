# Magic Memory Box Go

A lightweight and efficient conversation context manager for AI assistants. It stores chat history in User/System/Assistant format, ready to be used with Replicate, OpenAI, and other AI APIs.

---

## 📦 Installation

```bash
go get github.com/rmay1er/magic-memory-box-go
```

---

## 🚀 Quick Start

```go
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rmay1er/magic-memory-box-go/incache"
	"github.com/rmay1er/magic-memory-box-go/memorybox"
)

func main() {
	ctx := context.Background()

	// Initialize in-memory cache (or replace with Redis adapter for production)
	cache := incache.NewCache()

	// Configure memory box: store up to 10 messages, expire after 2 hours
	mb := memorybox.NewMemoryBox(cache, memorybox.MemoryBoxConfig{
		ContextLenSize: 10,
		ExpireTime:     2 * time.Hour,
	})

	question := "Привет, как дела?"

	// Ensure system prompt is present
	memories, err := mb.GetMemories(ctx, "ruslan")
	if err != nil {
		fmt.Printf("MemoryBox error: %v\n", err)
	}
	if len(memories) < 1 {
		mb.AddRaw(ctx, "ruslan", memorybox.System, "Ты нео из матрицы отвечай весело но")
	}

	// Add user input to conversation context
	userMessages, err := mb.Talk(ctx, "ruslan", question)
	if err != nil {
		fmt.Printf("MemoryBox error: %v\n", err)
	}

	// Here you would send userMessages to your AI service and get a reply
	// For demonstration, we simply echo the last user message
	reply := "Привет! У меня всё отлично."

	fmt.Printf("Нео: %s\n", reply)

	// Save AI response into memory box
	if _, err := mb.Remember(ctx, "ruslan", reply); err != nil {
		fmt.Printf("MemoryBox error: %v\n", err)
	}
}
```

---

## 🎯 Key Features

- **Context Management**: Automatic message length limiting and configurable TTL
- **Role Support**: System, User, Assistant, Tool roles supported
- **Storage Flexibility**: Works with in-memory cache or Redis backend via a unified interface
- **API Ready**: Messages formatted for Replicate, OpenAI, Claude, DeepSeek, etc.

---

## 🔧 Usage Examples

### In-Memory Cache

```go
import (
	"github.com/rmay1er/magic-memory-box-go/incache"
	"github.com/rmay1er/magic-memory-box-go/memorybox"
)

cache := incache.NewCache()
config := memorybox.MemoryBoxConfig{
	ContextLenSize: 15,
	ExpireTime:     24 * time.Hour,
}

mb := memorybox.NewMemoryBox(cache, config)
```

### Redis Storage

```go
import (
	"github.com/go-redis/redis/v8"
	"github.com/rmay1er/magic-memory-box-go/rdb"
	"github.com/rmay1er/magic-memory-box-go/memorybox"
)

client := redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

// You can add a prefix (e.g., "chat:") and specify a boolean value to clear all messages
// with this prefix from the database upon program exit on Ctrl+C.
redisAdapter := rdb.NewRedisAdapter(client, "chat:", true)

mb := memorybox.NewMemoryBox(redisAdapter, memorybox.MemoryBoxConfig{
	ContextLenSize: 20,
	ExpireTime:     48 * time.Hour,
})
```

### Conversation Handling

```go
mb.Talk(ctx, "userID", "Hello!")
mb.Remember(ctx, "userID", "Hi, how can I assist you today?")
mb.AddRaw(ctx, "userID", memorybox.System, "You are a helpful assistant")

messages, _ := mb.GetMemories(ctx, "userID") // API-ready message slice
```

---

## 🔄 Integration with Reptiloid (Replicate API)

```go
import (
	"context"
	"os"

	"github.com/rmay1er/magic-memory-box-go/incache"
	"github.com/rmay1er/magic-memory-box-go/memorybox"
	"github.com/rmay1er/reptiloid-go/reptiloid"
	"github.com/rmay1er/reptiloid-go/reptiloid/models/text"
)

func reptiloidExample() {
	ctx := context.Background()

	cache := incache.NewCache()
	mb := memorybox.NewMemoryBox(cache, memorybox.MemoryBoxConfig{
		ContextLenSize: 10,
		ExpireTime:     time.Hour,
	})

	// Add user/system messages as needed here

	client := reptiloid.NewClient(text.GPT41mini, os.Getenv("REPLICATE_API_TOKEN"))
	messages, _ := mb.GetMemories(ctx, "user123")

	resp, err := client.Generate(text.GPT4SeriesInput{
		Messages: convertMessagesForReplicate(messages), // convertMessages as shown above
	})
	if err != nil {
		// Handle error
	}

	// Save AI response
	mb.Remember(ctx, "user123", strings.Join(resp.Output, ""))
}
```

---

## 📊 Message Format Example

Stored messages are JSON formatted to suit AI APIs:

```json
[
  {"role": "system", "content": "You are a helpful assistant"},
  {"role": "user", "content": "Hello!"},
  {"role": "assistant", "content": "Hi! How can I help you today?"}
]
```

---

## 📈 Performance and Design

- Minimal overhead with automatic context expiration
- Efficient memory usage with context length limiting
- Support for switching storage backend without code changes
- No external dependencies required for in-memory mode

---

## 🤝 Contributing

Your contributions are welcome! Please open issues and pull requests for improvements or new features.

---

## 📄 License

MIT License — see the LICENSE file for details.

---

**Magic Memory Box Go** — simple and effective conversation context management for AI applications, with easy integration and flexible storage options.
