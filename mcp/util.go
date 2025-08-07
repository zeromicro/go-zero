package mcp

import (
	"fmt"
)

// ptr is a helper function to get a pointer to a value
func ptr[T any](v T) *T {
	return &v
}

// formatSSEMessage formats a Server-Sent Event message with proper CRLF line endings
func formatSSEMessage(event string, data []byte) string {
	return fmt.Sprintf("event: %s\r\ndata: %s\r\n\r\n", event, string(data))
}
