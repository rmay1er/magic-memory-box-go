package fantasy

import (
	"encoding/json"

	origfantasy "charm.land/fantasy"
	"github.com/rmay1er/magic-memory-box-go/memorybox"
)

func ToFantasy(msgs []memorybox.Message) []origfantasy.Message {
	messages := make([]origfantasy.Message, 0, len(msgs))
	for _, m := range msgs {
		var content []origfantasy.MessagePart
		switch m.Role {
		case "user", "system", "assistant":
			content = []origfantasy.MessagePart{origfantasy.TextPart{Text: m.Content}}
		case "tool":
			var data map[string]string
			if err := json.Unmarshal([]byte(m.Content), &data); err == nil {
				content = []origfantasy.MessagePart{origfantasy.ToolResultPart{
					ToolCallID: data["tool_call_id"],
					Output:     origfantasy.ToolResultOutputContentText{Text: data["content"]},
				}}
			} else {
				// Fallback to text if JSON parsing fails
				content = []origfantasy.MessagePart{origfantasy.TextPart{Text: m.Content}}
			}
		default:
			// Default to text part if role is unknown
			content = []origfantasy.MessagePart{origfantasy.TextPart{Text: m.Content}}
		}
		newmsg := origfantasy.Message{
			Role:    origfantasy.MessageRole(m.Role),
			Content: content,
		}
		messages = append(messages, newmsg)
	}
	return messages
}
