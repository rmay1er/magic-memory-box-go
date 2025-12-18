package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/fantasy"
	"charm.land/fantasy/providers/openrouter"
	convert "github.com/rmay1er/magic-memory-box-go/convert/fantasy"
	"github.com/rmay1er/magic-memory-box-go/memorybox"
)

// FantasyExample demonstrates usage of the memory box and AI client for a simple interactive conversation.
// It includes detailed comments to help even beginners understand each step.
func FantasyExample() {

	// Load environment variables from a .env file. This is where secrets like API tokens are stored.
	// err := godotenv.Load()
	// if err != nil {
	// 	// If loading the .env file fails, stop execution and show the error.
	// 	panic(err)
	// }

	// Create a context to manage request lifetimes and cancellation signals.
	var ctx = context.Background()

	// Initialize an in-memory cache which will store conversation memories.
	cache := memorybox.NewCache()

	// Create a new MemoryBox instance using the cache and configuration settings.
	// ContextLenSize defines how many recent messages to keep.
	// ExpireTime sets how long memories are kept before they are dropped.
	box := memorybox.NewMemoryBox(cache, memorybox.MemoryBoxConfig{
		ContextLenSize: 20,
		ExpireTime:     time.Hour,
	})

	// Choose your fave provider.
	provider, err := openrouter.New(openrouter.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Whoops:", err)
		os.Exit(1)
	}

	// Pick your fave model.
	model, err := provider.LanguageModel(ctx, "google/gemini-3-flash-preview")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Dang:", err)
		os.Exit(1)
	}

	weatherTool := fantasy.NewAgentTool("get_weather", "Get weather by city", WeatherTool)

	agent := fantasy.NewAgent(model, fantasy.WithTools(weatherTool))

	// Set up a reader to read user input from the standard input (keyboard).
	reader := bufio.NewReader(os.Stdin)

	// Begin an infinite loop to continuously prompt the user and respond.
	for {
		// Prompt the user for input.
		fmt.Print("You: ")

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

		// Retrieve existing memories for the user "user" from the MemoryBox.
		mems, err := box.GetMemories(ctx, "user")
		if err != nil {
			// On error fetching memories, print it and exit the loop.
			fmt.Printf("MemoryBox have no data: %v\n", err)
		}

		// If no memories exist yet, add a system message to guide the conversation style.
		if len(mems) < 1 {
			// The message means: "You are an AI - weather analysis tool"
			box.AddRaw(ctx, "user", memorybox.SystemRole, "You are an AI - weather analysis tool")
		}

		// Pass the user's input to the MemoryBox.Tell method,
		// which prepares the full conversation context for the AI.
		userMsgs, err := box.Tell(ctx, "user", input)
		if err != nil {
			// If something goes wrong, print the error and exit the loop.
			fmt.Printf("MemoryBox error: %v\n", err)

		}

		// Call the AI client to generate a response using the conversation messages.
		// Convert memorybox messages into the format expected by the AI client.
		resp, err := agent.Generate(ctx, fantasy.AgentCall{
			// Prompt:   input, // Was Error because empty.
			Messages: convert.ToFantasy(userMsgs),
		})
		if err != nil {
			// Print AI generation errors and stop.
			fmt.Printf("Generation error: %v\n", err)
			break
		}

		// Remember tool results if any tools were used
		for _, step := range resp.Steps {
			for _, tr := range step.Content.ToolResults() {
				if text, ok := fantasy.AsToolResultOutputType[fantasy.ToolResultOutputContentText](tr.Result); ok {
					data := map[string]string{"tool_call_id": tr.ToolCallID, "content": text.Text}
					jsonData, _ := json.Marshal(data)
					box.AddRaw(ctx, "user", memorybox.ToolRole, string(jsonData))
				}
			}
		}

		// Print the AI's response labeled as "AI:".
		fmt.Printf("AI: %#v\n", resp.Response.Content.Text())

		// Remember the AI's reply to append to the conversation history.
		_, err = box.Remember(ctx, "user", resp.Response.Content.Text())
		if err != nil {
			// Handle errors while saving memory.
			fmt.Printf("MemoryBox error: %v\n", err)
		}

		// Colored output of the messages array for testing
		fmt.Printf("\033[32mMessages: %+v\033[0m\n", userMsgs)
	}
}

type WeatherReq struct {
	Location string `json:"location" description:"City to get the weather for. In English."`
}

func WeatherTool(ctx context.Context, input WeatherReq, _ fantasy.ToolCall) (fantasy.ToolResponse, error) {
	if input.Location == "Moscow" {
		return fantasy.ToolResponse{Content: "Sunny in Moscow"}, nil
	} else if input.Location == "St. Petersburg" {
		return fantasy.ToolResponse{Content: "Rainy in St. Petersburg"}, nil
	} else if input.Location == "Helsinki" {
		return fantasy.ToolResponse{Content: "Cold in Helsinki"}, nil
	} else {
		return fantasy.ToolResponse{Content: "Sorry, I don't know the weather in this city."}, nil
	}
}

func main() {
	FantasyExample()
}
