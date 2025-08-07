package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Event struct {
	Type string
	Data map[string]any
}

func parseEvent(input string) (*Event, error) {
	var evt Event
	var dataStr string

	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "event:") {
			evt.Type = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
		} else if strings.HasPrefix(line, "data:") {
			dataStr = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(dataStr) > 0 {
		if err := json.Unmarshal([]byte(dataStr), &evt.Data); err != nil {
			return nil, fmt.Errorf("failed to parse data: %w", err)
		}
	}

	return &evt, nil
}

// TestToTypedPromptMessages tests the toTypedPromptMessages function
func TestToTypedPromptMessages(t *testing.T) {
	// Test with multiple message types in one test
	t.Run("MixedContentTypes", func(t *testing.T) {
		// Create test data with different content types
		messages := []PromptMessage{
			{
				Role: RoleUser,
				Content: TextContent{
					Text: "Hello, this is a text message",
					Annotations: &Annotations{
						Audience: []RoleType{RoleUser, RoleAssistant},
						Priority: ptr(0.8),
					},
				},
			},
			{
				Role: RoleAssistant,
				Content: ImageContent{
					Data:     "base64ImageData",
					MimeType: "image/jpeg",
				},
			},
			{
				Role: RoleUser,
				Content: AudioContent{
					Data:     "base64AudioData",
					MimeType: "audio/mp3",
				},
			},
			{
				Role:    "system",
				Content: "This is a simple string that should be handled as unknown type",
			},
		}

		// Call the function
		result := toTypedPromptMessages(messages)

		// Validate results
		require.Len(t, result, 4, "Should return the same number of messages")

		// Validate first message (TextContent)
		msg := result[0]
		assert.Equal(t, RoleUser, msg.Role, "Role should be preserved")

		// Type assertion using reflection since Content is an interface
		typed, ok := msg.Content.(typedTextContent)
		require.True(t, ok, "Should be typedTextContent")
		assert.Equal(t, ContentTypeText, typed.Type, "Type should be text")
		assert.Equal(t, "Hello, this is a text message", typed.Text, "Text content should be preserved")
		require.NotNil(t, typed.Annotations, "Annotations should be preserved")
		assert.Equal(t, []RoleType{RoleUser, RoleAssistant}, typed.Annotations.Audience, "Audience should be preserved")
		require.NotNil(t, typed.Annotations.Priority, "Priority should be preserved")
		assert.Equal(t, 0.8, *typed.Annotations.Priority, "Priority value should be preserved")

		// Validate second message (ImageContent)
		msg = result[1]
		assert.Equal(t, RoleAssistant, msg.Role, "Role should be preserved")

		// Type assertion for image content
		typedImg, ok := msg.Content.(typedImageContent)
		require.True(t, ok, "Should be typedImageContent")
		assert.Equal(t, ContentTypeImage, typedImg.Type, "Type should be image")
		assert.Equal(t, "base64ImageData", typedImg.Data, "Image data should be preserved")
		assert.Equal(t, "image/jpeg", typedImg.MimeType, "MimeType should be preserved")

		// Validate third message (AudioContent)
		msg = result[2]
		assert.Equal(t, RoleUser, msg.Role, "Role should be preserved")

		// Type assertion for audio content
		typedAudio, ok := msg.Content.(typedAudioContent)
		require.True(t, ok, "Should be typedAudioContent")
		assert.Equal(t, ContentTypeAudio, typedAudio.Type, "Type should be audio")
		assert.Equal(t, "base64AudioData", typedAudio.Data, "Audio data should be preserved")
		assert.Equal(t, "audio/mp3", typedAudio.MimeType, "MimeType should be preserved")

		// Validate fourth message (unknown type converted to TextContent)
		msg = result[3]
		assert.Equal(t, RoleType("system"), msg.Role, "Role should be preserved")

		// Should be converted to a typedTextContent with error message
		typedUnknown, ok := msg.Content.(typedTextContent)
		require.True(t, ok, "Unknown content should be converted to typedTextContent")
		assert.Equal(t, ContentTypeText, typedUnknown.Type, "Type should be text")
		assert.Contains(t, typedUnknown.Text, "Unknown content type:", "Should contain error about unknown type")
		assert.Contains(t, typedUnknown.Text, "string", "Should mention the actual type")
	})

	// Test empty input
	t.Run("EmptyInput", func(t *testing.T) {
		messages := []PromptMessage{}
		result := toTypedPromptMessages(messages)
		assert.Empty(t, result, "Should return empty slice for empty input")
	})

	// Test with nil annotations
	t.Run("NilAnnotations", func(t *testing.T) {
		messages := []PromptMessage{
			{
				Role: RoleUser,
				Content: TextContent{
					Text:        "Text with nil annotations",
					Annotations: nil,
				},
			},
		}

		result := toTypedPromptMessages(messages)
		require.Len(t, result, 1, "Should return one message")

		typed, ok := result[0].Content.(typedTextContent)
		require.True(t, ok, "Should be typedTextContent")
		assert.Equal(t, ContentTypeText, typed.Type, "Type should be text")
		assert.Equal(t, "Text with nil annotations", typed.Text, "Text content should be preserved")
		assert.Nil(t, typed.Annotations, "Nil annotations should be preserved as nil")
	})
}

// TestToTypedContents tests the toTypedContents function
func TestToTypedContents(t *testing.T) {
	// Test with multiple content types in one test
	t.Run("MixedContentTypes", func(t *testing.T) {
		// Create test data with different content types
		contents := []any{
			TextContent{
				Text: "Hello, this is a text content",
				Annotations: &Annotations{
					Audience: []RoleType{RoleUser, RoleAssistant},
					Priority: ptr(0.7),
				},
			},
			ImageContent{
				Data:     "base64ImageData",
				MimeType: "image/png",
			},
			AudioContent{
				Data:     "base64AudioData",
				MimeType: "audio/wav",
			},
			"This is a simple string that should be handled as unknown type",
		}

		// Call the function
		result := toTypedContents(contents)

		// Validate results
		require.Len(t, result, 4, "Should return the same number of contents")

		// Validate first content (TextContent)
		typed, ok := result[0].(typedTextContent)
		require.True(t, ok, "Should be typedTextContent")
		assert.Equal(t, ContentTypeText, typed.Type, "Type should be text")
		assert.Equal(t, "Hello, this is a text content", typed.Text, "Text content should be preserved")
		require.NotNil(t, typed.Annotations, "Annotations should be preserved")
		assert.Equal(t, []RoleType{RoleUser, RoleAssistant}, typed.Annotations.Audience, "Audience should be preserved")
		require.NotNil(t, typed.Annotations.Priority, "Priority should be preserved")
		assert.Equal(t, 0.7, *typed.Annotations.Priority, "Priority value should be preserved")

		// Validate second content (ImageContent)
		typedImg, ok := result[1].(typedImageContent)
		require.True(t, ok, "Should be typedImageContent")
		assert.Equal(t, ContentTypeImage, typedImg.Type, "Type should be image")
		assert.Equal(t, "base64ImageData", typedImg.Data, "Image data should be preserved")
		assert.Equal(t, "image/png", typedImg.MimeType, "MimeType should be preserved")

		// Validate third content (AudioContent)
		typedAudio, ok := result[2].(typedAudioContent)
		require.True(t, ok, "Should be typedAudioContent")
		assert.Equal(t, ContentTypeAudio, typedAudio.Type, "Type should be audio")
		assert.Equal(t, "base64AudioData", typedAudio.Data, "Audio data should be preserved")
		assert.Equal(t, "audio/wav", typedAudio.MimeType, "MimeType should be preserved")

		// Validate fourth content (unknown type converted to TextContent)
		typedUnknown, ok := result[3].(typedTextContent)
		require.True(t, ok, "Unknown content should be converted to typedTextContent")
		assert.Equal(t, ContentTypeText, typedUnknown.Type, "Type should be text")
		assert.Contains(t, typedUnknown.Text, "Unknown content type:", "Should contain error about unknown type")
		assert.Contains(t, typedUnknown.Text, "string", "Should mention the actual type")
	})

	// Test empty input
	t.Run("EmptyInput", func(t *testing.T) {
		contents := []any{}
		result := toTypedContents(contents)
		assert.Empty(t, result, "Should return empty slice for empty input")
	})

	// Test with nil annotations
	t.Run("NilAnnotations", func(t *testing.T) {
		contents := []any{
			TextContent{
				Text:        "Text with nil annotations",
				Annotations: nil,
			},
		}

		result := toTypedContents(contents)
		require.Len(t, result, 1, "Should return one content")

		typed, ok := result[0].(typedTextContent)
		require.True(t, ok, "Should be typedTextContent")
		assert.Equal(t, ContentTypeText, typed.Type, "Type should be text")
		assert.Equal(t, "Text with nil annotations", typed.Text, "Text content should be preserved")
		assert.Nil(t, typed.Annotations, "Nil annotations should be preserved as nil")
	})

	// Test with custom struct (should be handled as unknown type)
	t.Run("CustomStruct", func(t *testing.T) {
		type CustomContent struct {
			Data string
		}

		contents := []any{
			CustomContent{
				Data: "custom data",
			},
		}

		result := toTypedContents(contents)
		require.Len(t, result, 1, "Should return one content")

		typed, ok := result[0].(typedTextContent)
		require.True(t, ok, "Custom struct should be converted to typedTextContent")
		assert.Equal(t, ContentTypeText, typed.Type, "Type should be text")
		assert.Contains(t, typed.Text, "Unknown content type:", "Should contain error about unknown type")
		assert.Contains(t, typed.Text, "CustomContent", "Should mention the actual type")
	})
}
