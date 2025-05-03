package mcp

import "fmt"

// formatSSEMessage formats a Server-Sent Event message with proper CRLF line endings
func formatSSEMessage(event string, data []byte) string {
	return fmt.Sprintf("event: %s\r\ndata: %s\r\n\r\n", event, string(data))
}

// ptr is a helper function to get a pointer to a value
func ptr[T any](v T) *T {
	return &v
}

func toTypedContents(contents []any) []any {
	typedContents := make([]any, len(contents))

	for i, content := range contents {
		switch v := content.(type) {
		case TextContent:
			typedContents[i] = typedTextContent{
				Type:        ContentTypeText,
				TextContent: v,
			}
		case ImageContent:
			typedContents[i] = typedImageContent{
				Type:         ContentTypeImage,
				ImageContent: v,
			}
		case AudioContent:
			typedContents[i] = typedAudioContent{
				Type:         ContentTypeAudio,
				AudioContent: v,
			}
		default:
			typedContents[i] = typedTextContent{
				Type: ContentTypeText,
				TextContent: TextContent{
					Text: fmt.Sprintf("Unknown content type: %T", v),
				},
			}
		}
	}

	return typedContents
}

func toTypedPromptMessages(messages []PromptMessage) []PromptMessage {
	typedMessages := make([]PromptMessage, len(messages))

	for i, msg := range messages {
		switch v := msg.Content.(type) {
		case TextContent:
			typedMessages[i] = PromptMessage{
				Role: msg.Role,
				Content: typedTextContent{
					Type:        ContentTypeText,
					TextContent: v,
				},
			}
		case ImageContent:
			typedMessages[i] = PromptMessage{
				Role: msg.Role,
				Content: typedImageContent{
					Type:         ContentTypeImage,
					ImageContent: v,
				},
			}
		case AudioContent:
			typedMessages[i] = PromptMessage{
				Role: msg.Role,
				Content: typedAudioContent{
					Type:         ContentTypeAudio,
					AudioContent: v,
				},
			}
		default:
			typedMessages[i] = PromptMessage{
				Role: msg.Role,
				Content: typedTextContent{
					Type: ContentTypeText,
					TextContent: TextContent{
						Text: fmt.Sprintf("Unknown content type: %T", v),
					},
				},
			}
		}
	}

	return typedMessages
}

// validatePromptArguments checks if all required arguments are provided
// Returns a list of missing required arguments
func validatePromptArguments(prompt Prompt, providedArgs map[string]string) []string {
	var missingArgs []string

	for _, arg := range prompt.Arguments {
		if arg.Required {
			if value, exists := providedArgs[arg.Name]; !exists || len(value) == 0 {
				missingArgs = append(missingArgs, arg.Name)
			}
		}
	}

	return missingArgs
}
