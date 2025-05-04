package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseMarshaling(t *testing.T) {
	// Test that the Response struct marshals correctly
	resp := Response{
		JsonRpc: "2.0",
		ID:      123,
		Result: map[string]string{
			"key": "value",
		},
	}

	data, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"jsonrpc":"2.0"`)
	assert.Contains(t, string(data), `"id":123`)
	assert.Contains(t, string(data), `"result":{"key":"value"}`)

	// Test response with error
	respWithError := Response{
		JsonRpc: "2.0",
		ID:      456,
		Error: &errorObj{
			Code:    errCodeInvalidRequest,
			Message: "Invalid Request",
		},
	}

	data, err = json.Marshal(respWithError)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"jsonrpc":"2.0"`)
	assert.Contains(t, string(data), `"id":456`)
	assert.Contains(t, string(data), `"error":{"code":-32600,"message":"Invalid Request"}`)
}

func TestRequestUnmarshaling(t *testing.T) {
	// Test that the Request struct unmarshals correctly
	jsonStr := `{
		"jsonrpc": "2.0",
		"id": 789,
		"method": "test_method",
		"params": {"key": "value"}
	}`

	var req Request
	err := json.Unmarshal([]byte(jsonStr), &req)

	assert.NoError(t, err)
	assert.Equal(t, "2.0", req.JsonRpc)
	assert.Equal(t, int64(789), req.ID)
	assert.Equal(t, "test_method", req.Method)

	// Check params unmarshaled correctly
	var params map[string]string
	err = json.Unmarshal(req.Params, &params)
	assert.NoError(t, err)
	assert.Equal(t, "value", params["key"])
}

func TestToolStructs(t *testing.T) {
	// Test Tool struct
	tool := Tool{
		Name:        "test.tool",
		Description: "A test tool",
		InputSchema: InputSchema{
			Type: "object",
			Properties: map[string]any{
				"input": map[string]any{
					"type":        "string",
					"description": "Input parameter",
				},
			},
			Required: []string{"input"},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			return "result", nil
		},
	}

	// Verify fields are correct
	assert.Equal(t, "test.tool", tool.Name)
	assert.Equal(t, "A test tool", tool.Description)
	assert.Equal(t, "object", tool.InputSchema.Type)
	assert.Contains(t, tool.InputSchema.Properties, "input")
	propMap, ok := tool.InputSchema.Properties["input"].(map[string]any)
	assert.True(t, ok, "Property should be a map")
	assert.Equal(t, "string", propMap["type"])
	assert.NotNil(t, tool.Handler)

	// Verify JSON marshalling (which should exclude Handler function)
	data, err := json.Marshal(tool)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"test.tool"`)
	assert.Contains(t, string(data), `"description":"A test tool"`)
	assert.Contains(t, string(data), `"inputSchema":`)
	assert.NotContains(t, string(data), `"Handler":`)
}

func TestPromptStructs(t *testing.T) {
	// Test Prompt struct
	prompt := Prompt{
		Name:        "test.prompt",
		Description: "A test prompt description",
	}

	// Verify fields are correct
	assert.Equal(t, "test.prompt", prompt.Name)
	assert.Equal(t, "A test prompt description", prompt.Description)

	// Verify JSON marshalling
	data, err := json.Marshal(prompt)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"test.prompt"`)
	assert.Contains(t, string(data), `"description":"A test prompt description"`)
}

func TestResourceStructs(t *testing.T) {
	// Test Resource struct
	resource := Resource{
		Name:        "test.resource",
		URI:         "http://example.com/resource",
		Description: "A test resource",
	}

	// Verify fields are correct
	assert.Equal(t, "test.resource", resource.Name)
	assert.Equal(t, "http://example.com/resource", resource.URI)
	assert.Equal(t, "A test resource", resource.Description)

	// Verify JSON marshalling
	data, err := json.Marshal(resource)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"name":"test.resource"`)
	assert.Contains(t, string(data), `"uri":"http://example.com/resource"`)
	assert.Contains(t, string(data), `"description":"A test resource"`)
}

func TestContentTypes(t *testing.T) {
	// Test TextContent
	textContent := TextContent{
		Text: "Sample text",
		Annotations: &Annotations{
			Audience: []RoleType{RoleUser, RoleAssistant},
			Priority: ptr(1.0),
		},
	}

	data, err := json.Marshal(textContent)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"text":"Sample text"`)
	assert.Contains(t, string(data), `"audience":["user","assistant"]`)
	assert.Contains(t, string(data), `"priority":1`)

	// Test ImageContent
	imageContent := ImageContent{
		Data:     "base64data",
		MimeType: "image/png",
	}

	data, err = json.Marshal(imageContent)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"data":"base64data"`)
	assert.Contains(t, string(data), `"mimeType":"image/png"`)

	// Test AudioContent
	audioContent := AudioContent{
		Data:     "base64audio",
		MimeType: "audio/mp3",
	}

	data, err = json.Marshal(audioContent)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"data":"base64audio"`)
	assert.Contains(t, string(data), `"mimeType":"audio/mp3"`)
}

func TestCallToolResult(t *testing.T) {
	// Test CallToolResult
	result := CallToolResult{
		Result: Result{
			Meta: map[string]any{
				"progressToken": "token123",
			},
		},
		Content: []interface{}{
			TextContent{
				Text: "Sample result",
			},
		},
		IsError: false,
	}

	data, err := json.Marshal(result)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"_meta":{"progressToken":"token123"}`)
	assert.Contains(t, string(data), `"content":[{"text":"Sample result"}]`)
	assert.NotContains(t, string(data), `"isError":`)
}
