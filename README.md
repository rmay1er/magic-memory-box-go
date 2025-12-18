# Magic Memory Box Go

**Universal Memory for AI Models** — A conversation context management library for AI assistants and chatbots.

Magic Memory Box Go provides a simple and flexible way to store and manage conversation history with AI models. The library works with any AI service (OpenAI, Replicate, Claude, DeepSeek, etc.) and supports multiple storage backends.

## 🖼️ Need AI Image or Text Generation?

Want to quickly generate AI images or text with Replicate? Try [Reptiloid Go](https://github.com/rmay1er/reptiloid-go) — a simple and fast Go library that integrates perfectly with Magic Memory Box for memory-enabled AI conversations!

---

## 🎯 Who is this library for?

- **AI application developers** who need a ready-made conversation context management system
- **Chatbot creators** wanting to add memory of previous messages
- **Researchers** working with different AI models
- **Anyone who wants** to easily integrate "memory" into their AI projects

---

## ✨ Key Features

### 🧠 Smart Context Management
- Automatic history length limiting (keep only the last N messages)
- Configurable message lifetime (TTL)
- Role support: System, User, Assistant, Tool

### 🔄 Flexible Storage Options
- **Built-in memory** — Fast in-memory operation
- **Redis** — Distributed storage for production
- **Unified interface** — Easy switching between different storage backends

### 🤝 Compatibility with All AI Services
- Ready-made message format for OpenAI, Replicate, Claude, DeepSeek
- Easy integration with any AI service
- Multi-language support (English, Russian, etc.)

---

## 🚀 Quick Start

### Installation
```bash
go get github.com/rmay1er/magic-memory-box-go
```

### Simple Example
```go
package main

import (
    "context"
    "time"
    
    "github.com/rmay1er/magic-memory-box-go/memorybox"
)

func main() {
    ctx := context.Background()
    
    // Create in-memory storage
    cache := memorybox.NewCache()
    
    // Configure the memory box
    mb := memorybox.NewMemoryBox(cache, memorybox.MemoryBoxConfig{
        ContextLenSize: 10,     // Keep last 10 messages
        ExpireTime:     2 * time.Hour, // Auto-expire after 2 hours
    })
    
    // Add system message (assistant personality)
    mb.AddRaw(ctx, "user123", memorybox.SystemRole, "You are a helpful assistant")
    
    // User says something
    mb.Tell(ctx, "user123", "Hello! How are you?")
    
    // Get full history to send to AI model
    messages, _ := mb.GetMemories(ctx, "user123")
    
    // Send messages to any AI service...
    // Get response...
    
    // Save AI response to memory
    mb.Remember(ctx, "user123", "Hello! I'm doing great, thank you!")
}
```

---

## 📦 Storage Options

### 1. Built-in Memory (for development and testing)
```go
import "github.com/rmay1er/magic-memory-box-go/memorybox"

cache := memorybox.NewCache()
// Fast, simple, no external dependencies
```

### 2. Redis (for production)
```go
import (
    "github.com/go-redis/redis/v8"
    "github.com/rmay1er/magic-memory-box-go/memorybox"
    "github.com/rmay1er/magic-memory-box-go/rdb"
)

client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

redisAdapter := rdb.NewRedisAdapter(client, "chat:", true)
// Reliable, distributed, with persistence
```

---

## 🔗 AI Service Integration

### With Replicate (via Reptiloid)
```go
import "github.com/rmay1er/magic-memory-box-go/convert"

// Get conversation history
messages, _ := mb.GetMemories(ctx, "user123")

// Convert to Replicate format
replicateMessages := convert.ToReplicate(messages)

// Send to AI model...
```

### With OpenAI API
```go
// Ready format for OpenAI Chat Completion
openaiMessages := []map[string]string{
    {"role": "system", "content": "You are a helpful assistant"},
    {"role": "user", "content": "Hello!"},
    {"role": "assistant", "content": "Hi there!"},
}
// messages from GetMemories are already in this format!
```

### Integrations
- **[Reptiloid](https://github.com/rmay1er/reptiloid-go)**: Seamless integration with Replicate for AI image and text generation.
- **[Fantasy](https://github.com/charmbracelet/fantasy)**: Support for multi-provider AI agents with tool calling.
- **Planned**: Integration with popular [go-openai](https://github.com/sashabaranov/go-openai) library for ChatGPT models.

---

## 🎮 Use Cases

### Chatbot with Memory
Create a bot that remembers the entire conversation history with a user and can reference previous messages.

### Multi-user System
Store separate histories for each user with automatic cleanup of old conversations.

### A/B Testing Prompts
Easily switch between different system prompts for different user groups.

### Long Conversations
Manage lengthy dialogues by automatically trimming oldest messages while preserving context.

---

## 📊 Why Magic Memory Box?

| Feature | Magic Memory Box | DIY Implementation |
|---------|------------------|---------------------|
| **Ready Context** | ✅ Automatic | ❌ Need to write code |
| **TTL Support** | ✅ Built-in | ❌ Complex to implement |
| **Storage Switching** | ✅ 2 lines of code | ❌ Rewrite logic |
| **AI-ready Format** | ✅ Ready-to-use | ❌ Manual conversion |
| **Role Support** | ✅ System/User/Assistant | ❌ Need to design |

---

## 🛠️ Ready Examples

The repository includes complete working examples:

- **`example/cache/`** — Example with built-in memory
- **`example/redis/`** — Example with Redis for production

Run them to see the library in action immediately!

---

## 📈 Performance

- **Minimal latency** — Optimized for real-time use
- **Efficient memory usage** — Automatic cleanup of old messages
- **Scalability** — From one user to millions

---

## 🤝 Community & Support

Found a bug? Have an improvement idea? Want to add a new storage backend?

- **Issues** — Report problems
- **Pull Requests** — Suggest improvements
- **Discussions** — Share usage ideas

---

## 📄 License

MIT License — Free to use for any purpose.

---

Made with ❤️ for Go beginners exploring AI generation
