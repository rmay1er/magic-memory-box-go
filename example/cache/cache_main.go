package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/rmay1er/magic-memory-box-go/convert"
	"github.com/rmay1er/magic-memory-box-go/memorybox"
	"github.com/rmay1er/reptiloid-go/reptiloid"
	"github.com/rmay1er/reptiloid-go/reptiloid/models/text"
)

// CacheExample demonstrates usage of the memory box and AI client for a simple interactive conversation.
// It includes detailed comments to help even beginners understand each step.
func CacheExample() {

	// Load environment variables from a .env file. This is where secrets like API tokens are stored.
	err := godotenv.Load()
	if err != nil {
		// If loading the .env file fails, stop execution and show the error.
		panic(err)
	}

	// Create a context to manage request lifetimes and cancellation signals.
	var ctx = context.Background()

	// Initialize an in-memory cache which will store conversation memories.
	cache := memorybox.NewCache()

	// Create a new MemoryBox instance using the cache and configuration settings.
	// ContextLenSize defines how many recent messages to keep.
	// ExpireTime sets how long memories are kept before they are dropped.
	box := memorybox.NewMemoryBox(cache, memorybox.MemoryBoxConfig{
		ContextLenSize: 5,
		ExpireTime:     time.Hour,
	})

	// Create an AI client instance using the GPT-4 1.1 mini model and the API token from environment variables.
	client := reptiloid.NewClient(&text.GPT41mini, os.Getenv("REPLICATE_API_TOKEN"))

	// Set up a reader to read user input from the standard input (keyboard).
	reader := bufio.NewReader(os.Stdin)

	// Begin an infinite loop to continuously prompt the user and respond.
	for {
		// Prompt printed in Russian meaning "You:"
		fmt.Print("Вы: ")

		// Read the user's input until they press Enter.
		input, err := reader.ReadString('\n')
		if err != nil {
			// If an error occurs while reading input, print the error and exit the loop.
			fmt.Printf("Input error: %v\n", err)
			break
		}

		// Trim whitespace and newlines from the user's input.
		input = strings.TrimSpace(input)

		// If the input is empty (just Enter pressed), skip this loop iteration.
		if input == "" {
			continue
		}

		// Retrieve existing memories for the user "ruslan" from the MemoryBox.
		mems, err := box.GetMemories(ctx, "ruslan")
		if err != nil {
			// On error fetching memories, print it and exit the loop.
			fmt.Printf("MemoryBox error: %v\n", err)
		}

		// If no memories exist yet, add a system message to guide the conversation style.
		if len(mems) < 1 {
			// The message roughly means: "You are Neo from the Matrix, respond cheerfully but ..."
			box.AddRaw(ctx, "ruslan", memorybox.SystemRole, "You are Neo from Matrix")
		}

		// Pass the user's input to the MemoryBox.Talk method,
		// which prepares the full conversation context for the AI.
		userMsgs, err := box.Tell(ctx, "ruslan", input)
		if err != nil {
			// If something goes wrong, print the error and exit the loop.
			fmt.Printf("MemoryBox error: %v\n", err)

		}

		// Call the AI client to generate a response using the conversation messages.
		// Convert memorybox messages into the format expected by the AI client.
		resp, err := client.Generate(&text.GPT4SeriesInput{
			Messages: convert.ToReplicate(userMsgs),
		})
		if err != nil {
			// Print AI generation errors and stop.
			fmt.Printf("Generation error: %v\n", err)
			break
		}

		// Join the AI output pieces into a single string and trim spaces.
		reply := strings.TrimSpace(strings.Join(resp.Output, ""))

		// Print the AI's response labeled as "Нео:" (Neo).
		fmt.Printf("Нео: %#v\n", resp.Output)

		// Remember the AI's reply to append to the conversation history.
		_, err = box.Remember(ctx, "ruslan", reply)
		if err != nil {
			// Handle errors while saving memory.
			fmt.Printf("MemoryBox error: %v\n", err)
		}
	}
}

func main() {
	CacheExample()
}
