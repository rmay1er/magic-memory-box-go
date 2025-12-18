package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	convert "github.com/rmay1er/magic-memory-box-go/convert/replicate"
	"github.com/rmay1er/magic-memory-box-go/memorybox"
	"github.com/rmay1er/magic-memory-box-go/rdb"
	"github.com/rmay1er/reptiloid-go/reptiloid"
	"github.com/rmay1er/reptiloid-go/reptiloid/models/text"
)

// RedisExample demonstrates how to use the Redis MemoryBox with the Reptiloid client.
func RedisExample() {
	// load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		// stop execution if .env could not be loaded
		panic(err)
	}

	// create a context to be used in Redis calls and requests
	var ctx = context.Background()

	// get Redis server address from environment variables (should be "localhost:6379")
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("localhost:6379"),
		DB:   0, // select Redis database 0 (default DB)
	})

	// wrap Redis client into adapter for MemoryBox usage, with namespace "tests"
	rdb := rdb.NewRedisAdapter(redisClient, "tests", true)

	// configure MemoryBox to keep last 10 messages for context and set expiration for memories after 2 hours
	box := memorybox.NewMemoryBox(rdb, memorybox.MemoryBoxConfig{
		ContextLenSize: 10,
		ExpireTime:     2 * time.Hour,
	})

	//Set Default Model
	model := &text.GPT41mini

	// create a new text generation client with model and API token
	client := reptiloid.NewClient(model, os.Getenv("REPLICATE_API_TOKEN"))

	// create a buffered reader to read user input from the console
	reader := bufio.NewReader(os.Stdin)

	// Start interactive loop to continuously read user input and generate responses
	for {
		// prompt user to enter text (in Russian)
		fmt.Print("Вы: ")

		// read until newline character
		input, err := reader.ReadString('\n')
		if err != nil {
			// print error message
			fmt.Printf("Input error: %v\n", err)
			// stop loop on error
			break
		}

		// remove leading and trailing white spaces including newline
		input = strings.TrimSpace(input)
		// if input is empty (user pressed enter without typing), skip this iteration and ask for input again
		if input == "" {
			continue
		}

		// retrieve existing memories for user "ruslan"
		mems, err := box.GetMemories(ctx, "ruslan")
		if err != nil {
			// print error message
			fmt.Printf("MemoryBox error: %v\n", err)
			// stop loop on error
			break
		}

		// if no memories found (first time interaction),
		// add an initial raw system message to set behavior (role-playing as Neo from Matrix)
		if len(mems) < 1 {
			_, _ = box.AddRaw(ctx, "ruslan", memorybox.SystemRole, "You are Neo from Matrix") // ignore error for AddRaw in this example
		}

		// update MemoryBox with user input and get conversation history
		userMsgs, err := box.Tell(ctx, "ruslan", input)
		if err != nil {
			// print error if Talk fails
			fmt.Printf("MemoryBox error: %v\n", err)
			break
		}

		// convert conversation messages to model input format and generate response
		resp, err := client.Generate(&text.GPT4SeriesInput{
			Messages: convert.ToReplicate(userMsgs), // or use inline convert.ToReplicate(box.TellUnsafe(ctx, "ruslan", input))
		})
		if err != nil {
			// print error message
			fmt.Printf("Generation error: %v\n", err)
			// stop loop on error
			break
		}

		// join output strings, trim whitespace
		reply := strings.TrimSpace(strings.Join(resp.Output, ""))
		// print generated response prefixed "Нео: "
		fmt.Printf("Нео: %s\n", reply)

		// save the generated response into memory
		_, err = box.Remember(ctx, "ruslan", reply)
		if err != nil {
			// print error if remembering failed
			fmt.Printf("MemoryBox error: %v\n", err)
			break
		}
	}
}
