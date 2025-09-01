package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	// "github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/rmay1er/magic-memory-box-go/memorybox"
	"github.com/rmay1er/magic-memory-box-go/incache"
	// "github.com/rmay1er/magic-memory-box-go/rdb"
	"github.com/rmay1er/reptiloid-go/reptiloid"
	"github.com/rmay1er/reptiloid-go/reptiloid/models/text"
)

var ctx = context.Background()

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cache := incache.NewCache()

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: os.Getenv("localhost:6379"), // e.g. "localhost:6379"
	// 	DB:   0,                           // use default DB
	// })

	// rdb := rdb.NewRedisAdapter(redisClient, "tests", true)

	box := memorybox.NewMemoryBox(cache, memorybox.MemoryBoxConfig{
		ContextLenSize: 10,
		ExpireTime:     2 * time.Hour,
	})

	client := reptiloid.NewClient(text.GPT41mini, os.Getenv("REPLICATE_API_TOKEN"))

	reader := bufio.NewReader(os.Stdin)

	// Start interactive loop
	for {
		fmt.Print("Вы: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Input error: %v\n", err)
			break
		}
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		mems, err := box.GetMemories(ctx, "ruslan")
		if err != nil {
			fmt.Printf("MemoryBox error: %v\n", err)
			break
		}
		if len(mems) < 1 {
			box.AddRaw(ctx, "ruslan", memorybox.System, "Ты нео из матрицы отвечай весело но ")
		}

		userMsgs, err := box.Talk(ctx, "ruslan", input)
		if err != nil {
			fmt.Printf("MemoryBox error: %v\n", err)
			break
		}

		resp, err := client.Generate(text.GPT4SeriesInput{
			Messages: convertMessages(userMsgs),
		})
		if err != nil {
			fmt.Printf("Generation error: %v\n", err)
			break
		}

		reply := strings.TrimSpace(strings.Join(resp.Output, ""))
		fmt.Printf("Нео: %s\n", reply)
		_, err = box.Remember(ctx, "ruslan", reply)
		if err != nil {
			fmt.Printf("MemoryBox error: %v\n", err)
			break
		}
	}
}

func convertMessages(msgs []memorybox.Message) []map[string]any {
	out := make([]map[string]any, len(msgs))
	for i, m := range msgs {
		newMap := make(map[string]any)
		newMap["role"] = string(m.Role)
		newMap["content"] = m.Content
		out[i] = newMap
	}
	return out
}
