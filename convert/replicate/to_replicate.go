package convert

import "github.com/rmay1er/magic-memory-box-go/memorybox"

// ConvertMessagesForReplicate converts a slice of Message structs into a slice of maps with keys "role" and "content",
// suitable for use with the Replicate API.
func ToReplicate(msgs []memorybox.Message) []map[string]any {
	out := make([]map[string]any, len(msgs))
	for i, m := range msgs {
		newMap := make(map[string]any)
		newMap["role"] = string(m.Role)
		newMap["content"] = m.Content
		out[i] = newMap
	}
	return out
}
