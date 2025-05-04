package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx/logtest"
)

// mockMcpServer is a helper for testing the MCP server
// It encapsulates the server and test server setup and teardown logic
type mockMcpServer struct {
	server     *sseMcpServer
	testServer *httptest.Server
	requestId  int64
}

// newMockMcpServer initializes a mock MCP server for testing
func newMockMcpServer(t *testing.T) *mockMcpServer {
	const yamlConf = `name: test-server
host: localhost
port: 8080
mcp:
  name: mcp-test-server
  messageTimeout: 5s
`

	var c McpConf
	assert.NoError(t, conf.LoadFromYamlBytes([]byte(yamlConf), &c))

	server := NewMcpServer(c).(*sseMcpServer)
	mux := http.NewServeMux()
	mux.HandleFunc(c.Mcp.SseEndpoint, server.handleSSE)
	mux.HandleFunc(c.Mcp.MessageEndpoint, server.handleRequest)
	testServer := httptest.NewServer(mux)
	server.conf.Mcp.BaseUrl = testServer.URL

	return &mockMcpServer{
		server:     server,
		testServer: testServer,
		requestId:  1,
	}
}

// shutdown closes the test server
func (m *mockMcpServer) shutdown() {
	m.testServer.Close()
}

// registerExamplePrompt registers a test prompt
func (m *mockMcpServer) registerExamplePrompt() {
	m.server.RegisterPrompt(Prompt{
		Name:        "test.prompt",
		Description: "A test prompt",
	})
}

// registerExampleResource registers a test resource
func (m *mockMcpServer) registerExampleResource() {
	m.server.RegisterResource(Resource{
		Name:        "test.resource",
		URI:         "file:///test.file",
		Description: "A test resource",
	})
}

// registerExampleTool registers a test tool
func (m *mockMcpServer) registerExampleTool() {
	_ = m.server.RegisterTool(Tool{
		Name:        "test.tool",
		Description: "A test tool",
		InputSchema: InputSchema{
			Properties: map[string]any{
				"input": map[string]any{
					"type":        "string",
					"description": "Input parameter",
				},
			},
			Required: []string{"input"},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			input, ok := params["input"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid input parameter")
			}
			return fmt.Sprintf("Processed: %s", input), nil
		},
	})
}

// Helper function to create and add a test client
func addTestClient(server *sseMcpServer, clientID string, initialized bool) *mcpClient {
	client := &mcpClient{
		id:          clientID,
		channel:     make(chan string, eventChanSize),
		initialized: initialized,
	}
	server.clientsLock.Lock()
	server.clients[clientID] = client
	server.clientsLock.Unlock()
	return client
}

func TestNewMcpServer(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	mock.registerExamplePrompt()
	mock.registerExampleResource()
	mock.registerExampleTool()

	require.NotNil(t, mock.server, "Server should be created")
	assert.NotEmpty(t, mock.server.tools, "Tools map should be initialized")
	assert.NotEmpty(t, mock.server.prompts, "Prompts map should be initialized")
	assert.NotEmpty(t, mock.server.resources, "Resources map should be initialized")
}

func TestNewMcpServer_WithCors(t *testing.T) {
	const yamlConf = `name: test-server
host: localhost
port: 8080
mcp:
  cors:
    - http://localhost:3000
  messageTimeout: 5s
`

	var c McpConf
	assert.NoError(t, conf.LoadFromYamlBytes([]byte(yamlConf), &c))

	server := NewMcpServer(c).(*sseMcpServer)
	assert.Equal(t, "test-server", server.conf.Name, "Server name should be set")
}

func TestHandleRequest_badRequest(t *testing.T) {
	t.Run("empty session ID", func(t *testing.T) {
		mock := newMockMcpServer(t)
		defer mock.shutdown()

		// Create a request with an invalid session ID
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      1,
			Method:  methodToolsCall,
			Params:  []byte(`{"sessionId": "invalid-session"}`),
		}

		jsonBody, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "/?session_id=", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		mock.server.handleRequest(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("bad body", func(t *testing.T) {
		mock := newMockMcpServer(t)
		defer mock.shutdown()

		addTestClient(mock.server, "test-session", true)

		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-session", bytes.NewReader([]byte(`{`)))
		w := httptest.NewRecorder()
		mock.server.handleRequest(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRegisterTool(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	tool := Tool{
		Name:        "example.tool",
		Description: "An example tool",
		InputSchema: InputSchema{
			Properties: map[string]any{
				"input": map[string]any{
					"type":        "string",
					"description": "Input parameter",
				},
			},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			return "result", nil
		},
	}

	// Test with valid tool
	err := mock.server.RegisterTool(tool)
	assert.NoError(t, err, "Should not error with valid tool")

	// Check tool was registered
	_, exists := mock.server.tools["example.tool"]
	assert.True(t, exists, "Tool should be registered")

	// Test with missing handler
	invalidTool := tool
	invalidTool.Name = "invalid.tool"
	invalidTool.Handler = nil
	err = mock.server.RegisterTool(invalidTool)
	assert.Error(t, err, "Should error with missing handler")
}

func TestRegisterPrompt(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	prompt := Prompt{
		Name:        "example.prompt",
		Description: "An example prompt",
	}

	// Test registering prompt
	mock.server.RegisterPrompt(prompt)

	// Check prompt was registered
	_, exists := mock.server.prompts["example.prompt"]
	assert.True(t, exists, "Prompt should be registered")
}

func TestRegisterResource(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	resource := Resource{
		Name:        "example.resource",
		URI:         "http://example.com/resource",
		Description: "An example resource",
	}

	// Test registering resource
	mock.server.RegisterResource(resource)

	// Check resource was registered
	_, exists := mock.server.resources["http://example.com/resource"]
	assert.True(t, exists, "Resource should be registered")
}

// TestToolCallBasic tests the basic functionality of a tool call
func TestToolsList(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Register a test tool
	mock.registerExampleTool()

	// Simulate a client to test tool call
	client := addTestClient(mock.server, "test-client", true)

	// Create a tool call request
	params := struct {
		Cursor string `json:"cursor"`
		Meta   struct {
			ProgressToken any `json:"progressToken"`
		} `json:"_meta"`
	}{
		Cursor: "my-cursor",
	}

	paramBytes, _ := json.Marshal(params)
	req := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  methodToolsList,
		Params:  paramBytes,
	}

	// Process the tool call
	mock.server.processListTools(context.Background(), client, req)

	// Get the response from the client's channel
	select {
	case response := <-client.channel:
		evt, err := parseEvent(response)
		assert.NoError(t, err)

		assert.Equal(t, eventMessage, evt.Type, "Event type should be message")
		result, ok := evt.Data["result"].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "my-cursor", result["nextCursor"], "Cursor should match")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tool call response")
	}
}

// TestToolCallBasic tests the basic functionality of a tool call
func TestToolCallBasic(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Register a test tool
	mock.registerExampleTool()

	// Simulate a client to test tool call
	client := addTestClient(mock.server, "test-client", true)

	// Create a tool call request
	params := struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}{
		Name: "test.tool",
		Arguments: map[string]any{
			"input": "test-input",
		},
	}

	paramBytes, _ := json.Marshal(params)
	req := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  methodToolsCall,
		Params:  paramBytes,
	}

	// Process the tool call
	mock.server.processToolCall(context.Background(), client, req)

	// Get the response from the client's channel
	select {
	case response := <-client.channel:
		// Check response format
		assert.Contains(t, response, "event: message", "Response should have message event")
		assert.Contains(t, response, "data:", "Response should have data")

		// Extract JSON from the SSE response
		jsonStart := strings.Index(response, "{")
		jsonEnd := strings.LastIndex(response, "}")
		jsonStr := response[jsonStart : jsonEnd+1]

		// Parse the JSON
		var parsed struct {
			Result struct {
				Content []map[string]any `json:"content"`
				IsError bool             `json:"isError"`
			} `json:"result"`
		}

		err := json.Unmarshal([]byte(jsonStr), &parsed)
		require.NoError(t, err, "Response should be valid JSON")

		// Verify the response content
		assert.Len(t, parsed.Result.Content, 1, "Response should contain one content item")
		assert.Equal(t, "Processed: test-input", parsed.Result.Content[0][ContentTypeText], "Tool result incorrect")
		assert.False(t, parsed.Result.IsError, "Response should not be an error")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tool call response")
	}
}

// TestToolCallMapResult tests a tool that returns a map[string]any result
func TestToolCallMapResult(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Register a tool that returns a map
	mapTool := Tool{
		Name:        "map.tool",
		Description: "A tool that returns a map result",
		InputSchema: InputSchema{
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			// Return a complex nested map structure
			return map[string]any{
				"string":  "value",
				"number":  42,
				"boolean": true,
				"nested": map[string]any{
					"array": []string{"item1", "item2"},
					"obj": map[string]any{
						"key": "value",
					},
				},
				"nullValue": nil,
			}, nil
		},
	}

	err := mock.server.RegisterTool(mapTool)
	require.NoError(t, err)

	// Create a client
	client := addTestClient(mock.server, "test-client", true)

	// Create a tool call request
	params := struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}{
		Name:      "map.tool",
		Arguments: map[string]any{},
	}

	paramBytes, _ := json.Marshal(params)
	req := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  methodToolsCall,
		Params:  paramBytes,
	}

	// Process the tool call
	mock.server.processToolCall(context.Background(), client, req)

	// Get the response from the client's channel
	select {
	case response := <-client.channel:
		// Parse the response
		jsonStart := strings.Index(response, "{")
		jsonEnd := strings.LastIndex(response, "}")
		require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")

		jsonStr := response[jsonStart : jsonEnd+1]
		var parsed map[string]any
		err := json.Unmarshal([]byte(jsonStr), &parsed)
		require.NoError(t, err, "Response should be valid JSON")

		// Get the result
		result, ok := parsed["result"].(map[string]any)
		require.True(t, ok, "Response should have a result object")

		// Get the content array
		content, ok := result["content"].([]any)
		require.True(t, ok, "Result should have a content array")
		require.NotEmpty(t, content, "Content should not be empty")

		// The first content item should be our map result converted to JSON text
		firstItem, ok := content[0].(map[string]any)
		require.True(t, ok, "First content item should be an object")

		// Get the text content which should be our JSON
		text, ok := firstItem[ContentTypeText].(string)
		require.True(t, ok, "Content should have text")

		// Verify the text is valid JSON and contains our data
		var mapResult map[string]any
		err = json.Unmarshal([]byte(text), &mapResult)
		require.NoError(t, err, "Text should be valid JSON")

		// Verify the content of our map
		assert.Equal(t, "value", mapResult["string"], "String value should match")
		assert.Equal(t, float64(42), mapResult["number"], "Number value should match")
		assert.Equal(t, true, mapResult["boolean"], "Boolean value should match")

		// Check nested structure
		nested, ok := mapResult["nested"].(map[string]any)
		require.True(t, ok, "Should have nested map")

		array, ok := nested["array"].([]any)
		require.True(t, ok, "Should have array in nested map")
		assert.Len(t, array, 2, "Array should have 2 items")
		assert.Equal(t, "item1", array[0], "First array item should match")

		obj, ok := nested["obj"].(map[string]any)
		require.True(t, ok, "Should have obj in nested map")
		assert.Equal(t, "value", obj["key"], "Nested object key should match")

		// Check null value
		_, exists := mapResult["nullValue"]
		assert.True(t, exists, "Null value key should exist")

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tool call response")
	}
}

// TestToolCallArrayResult tests a tool that returns an array result
func TestToolCallArrayResult(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Register a tool that returns an array
	arrayTool := Tool{
		Name:        "array.tool",
		Description: "A tool that returns an array result",
		InputSchema: InputSchema{
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			// Return an array of mixed content types
			return []any{
				"string item",
				42,
				true,
				map[string]any{"key": "value"},
				[]string{"nested", "array"},
				nil,
			}, nil
		},
	}

	err := mock.server.RegisterTool(arrayTool)
	require.NoError(t, err)

	// Create a client
	client := addTestClient(mock.server, "test-client", true)

	// Create a tool call request
	params := struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}{
		Name:      "array.tool",
		Arguments: map[string]any{},
	}

	paramBytes, _ := json.Marshal(params)
	req := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  methodToolsCall,
		Params:  paramBytes,
	}

	// Process the tool call
	mock.server.processToolCall(context.Background(), client, req)

	// Get the response from the client's channel
	select {
	case response := <-client.channel:
		// Parse the response
		jsonStart := strings.Index(response, "{")
		jsonEnd := strings.LastIndex(response, "}")
		require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")

		jsonStr := response[jsonStart : jsonEnd+1]
		var parsed map[string]any
		err := json.Unmarshal([]byte(jsonStr), &parsed)
		require.NoError(t, err, "Response should be valid JSON")

		// Get the result
		result, ok := parsed["result"].(map[string]any)
		require.True(t, ok, "Response should have a result object")

		// Get the content array
		content, ok := result["content"].([]any)
		require.True(t, ok, "Result should have a content array")
		require.Equal(t, 6, len(content), "Content should have 6 items, one for each array item")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tool call response")
	}
}

// TestToolCallTextContentResult tests a tool that returns a TextContent result
func TestToolCallTextContentResult(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Register a tool that returns a TextContent
	textContentTool := Tool{
		Name:        "text.content.tool",
		Description: "A tool that returns a TextContent result",
		InputSchema: InputSchema{
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			// Return a TextContent object directly
			return TextContent{
				Text: "This is a direct TextContent result",
				Annotations: &Annotations{
					Audience: []RoleType{RoleUser, RoleAssistant},
					Priority: func() *float64 { p := 0.9; return &p }(),
				},
			}, nil
		},
	}

	err := mock.server.RegisterTool(textContentTool)
	require.NoError(t, err)

	// Create a client
	client := addTestClient(mock.server, "test-client", true)

	// Create a tool call request
	params := struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}{
		Name:      "text.content.tool",
		Arguments: map[string]any{},
	}

	paramBytes, _ := json.Marshal(params)
	req := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  methodToolsCall,
		Params:  paramBytes,
	}

	// Process the tool call
	mock.server.processToolCall(context.Background(), client, req)

	// Get the response from the client's channel
	select {
	case response := <-client.channel:
		// Parse the response
		jsonStart := strings.Index(response, "{")
		jsonEnd := strings.LastIndex(response, "}")
		require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")

		jsonStr := response[jsonStart : jsonEnd+1]
		var parsed map[string]any
		err := json.Unmarshal([]byte(jsonStr), &parsed)
		require.NoError(t, err, "Response should be valid JSON")

		// Get the result
		result, ok := parsed["result"].(map[string]any)
		require.True(t, ok, "Response should have a result object")

		// Get the content array
		content, ok := result["content"].([]any)
		require.True(t, ok, "Result should have a content array")
		require.NotEmpty(t, content, "Content should not be empty")

		// The first content item should be our TextContent
		firstItem, ok := content[0].(map[string]any)
		require.True(t, ok, "First content item should be an object")

		// Check annotations
		annotations, ok := firstItem["annotations"].(map[string]any)
		require.True(t, ok, "Should have annotations")

		audience, ok := annotations["audience"].([]any)
		require.True(t, ok, "Should have audience in annotations")
		assert.Len(t, audience, 2, "Audience should have 2 items")

		priority, ok := annotations["priority"].(float64)
		require.True(t, ok, "Should have priority in annotations")
		assert.Equal(t, 0.9, priority, "Priority should match")

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tool call response")
	}
}

// TestToolCallImageContentResult tests a tool that returns an ImageContent result
func TestToolCallImageContentResult(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Register a tool that returns an ImageContent
	imageContentTool := Tool{
		Name:        "image.content.tool",
		Description: "A tool that returns an ImageContent result",
		InputSchema: InputSchema{
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			// Return an ImageContent object directly
			return ImageContent{
				Data:     "dGVzdCBpbWFnZSBkYXRhIChiYXNlNjQgZW5jb2RlZCk=", // "test image data (base64 encoded)" in base64
				MimeType: "image/png",
			}, nil
		},
	}

	err := mock.server.RegisterTool(imageContentTool)
	require.NoError(t, err)

	// Create a client
	client := addTestClient(mock.server, "test-client", true)

	// Create a tool call request
	params := struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}{
		Name:      "image.content.tool",
		Arguments: map[string]any{},
	}

	paramBytes, _ := json.Marshal(params)
	req := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  methodToolsCall,
		Params:  paramBytes,
	}

	// Process the tool call
	mock.server.processToolCall(context.Background(), client, req)

	// Get the response from the client's channel
	select {
	case response := <-client.channel:
		// Parse the response
		jsonStart := strings.Index(response, "{")
		jsonEnd := strings.LastIndex(response, "}")
		require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")

		jsonStr := response[jsonStart : jsonEnd+1]
		var parsed map[string]any
		err := json.Unmarshal([]byte(jsonStr), &parsed)
		require.NoError(t, err, "Response should be valid JSON")

		// Get the result
		result, ok := parsed["result"].(map[string]any)
		require.True(t, ok, "Response should have a result object")

		// Get the content array
		content, ok := result["content"].([]any)
		require.True(t, ok, "Result should have a content array")
		require.NotEmpty(t, content, "Content should not be empty")

		// The first content item should be our ImageContent
		firstItem, ok := content[0].(map[string]any)
		require.True(t, ok, "First content item should be an object")

		// Check image data
		data, ok := firstItem["data"].(string)
		require.True(t, ok, "Content should have data")
		assert.Equal(t, "dGVzdCBpbWFnZSBkYXRhIChiYXNlNjQgZW5jb2RlZCk=", data, "Image data should match")

		// Check mime type
		mimeType, ok := firstItem["mimeType"].(string)
		require.True(t, ok, "Content should have mimeType")
		assert.Equal(t, "image/png", mimeType, "MimeType should match")

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tool call response")
	}
}

// TestToolCallToolResultType tests a tool that returns a ToolResult type
func TestToolCallToolResultType(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	toolResultTool := Tool{
		Name:        "toolresult.tool",
		Description: "A tool that returns a ToolResult object",
		InputSchema: InputSchema{
			Type:       ContentTypeObject,
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			return ToolResult{
				Type:    ContentTypeText,
				Content: "This is a ToolResult with text content type",
			}, nil
		},
	}
	err := mock.server.RegisterTool(toolResultTool)
	require.NoError(t, err)

	toolResultImageTool := Tool{
		Name:        "toolresult.image.tool",
		Description: "A tool that returns a ToolResult with image content",
		InputSchema: InputSchema{
			Type:       ContentTypeObject,
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			return ToolResult{
				Type: "image",
				Content: map[string]any{
					"data":     "dGVzdCBpbWFnZSBkYXRhIGZvciB0b29sIHJlc3VsdA==", // "test image data for tool result" in base64
					"mimeType": "image/jpeg",
				},
			}, nil
		},
	}
	err = mock.server.RegisterTool(toolResultImageTool)
	require.NoError(t, err)

	toolResultAudioTool := Tool{
		Name:        "toolresult.audio.tool",
		Description: "A tool that returns a ToolResult with audio content",
		InputSchema: InputSchema{
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			// Test with image type
			return ToolResult{
				Type: "audio",
				Content: map[string]any{
					"data":     "dGVzdCBpbWFnZSBkYXRhIGZvciB0b29sIHJlc3VsdA==", // "test image data for tool result" in base64
					"mimeType": "audio",
				},
			}, nil
		},
	}
	err = mock.server.RegisterTool(toolResultAudioTool)
	require.NoError(t, err)

	toolResultIntType := Tool{
		Name:        "toolresult.int.tool",
		Description: "A tool that returns a ToolResult with int content",
		InputSchema: InputSchema{
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			return 2, nil
		},
	}
	err = mock.server.RegisterTool(toolResultIntType)
	require.NoError(t, err)

	toolResultBadType := Tool{
		Name:        "toolresult.bad.tool",
		Description: "A tool that returns a ToolResult with bad content",
		InputSchema: InputSchema{
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			return map[string]any{
				"type": "custom",
				"data": make(chan int),
			}, nil
		},
	}
	err = mock.server.RegisterTool(toolResultBadType)
	require.NoError(t, err)

	// Test text ToolResult
	t.Run("textToolResult", func(t *testing.T) {
		// Create a client
		client := addTestClient(mock.server, "test-client-text", true)

		// Create a tool call request
		params := struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}{
			Name:      "toolresult.tool",
			Arguments: map[string]any{},
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      1,
			Method:  methodToolsCall,
			Params:  paramBytes,
		}

		// Process the tool call
		mock.server.processToolCall(context.Background(), client, req)

		// Get the response from the client's channel
		select {
		case response := <-client.channel:
			// Parse the response
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")

			jsonStr := response[jsonStart : jsonEnd+1]
			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Response should be valid JSON")

			// Get the result
			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			// Get the content array
			content, ok := result["content"].([]any)
			require.True(t, ok, "Result should have a content array")
			require.NotEmpty(t, content, "Content should not be empty")

			// The first content item should be converted from ToolResult to TextContent
			firstItem, ok := content[0].(map[string]any)
			require.True(t, ok, "First content item should be an object")

			// Check text content
			text, ok := firstItem[ContentTypeText].(string)
			require.True(t, ok, "Content should have text")
			assert.Equal(t, "This is a ToolResult with text content type", text, "Text content should match")

		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for tool call response")
		}
	})

	// Test image ToolResult
	t.Run("imageToolResult", func(t *testing.T) {
		// Create a client
		client := addTestClient(mock.server, "test-client-image", true)

		// Create a tool call request
		params := struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}{
			Name:      "toolresult.image.tool",
			Arguments: map[string]any{},
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      2,
			Method:  methodToolsCall,
			Params:  paramBytes,
		}

		// Process the tool call
		mock.server.processToolCall(context.Background(), client, req)

		// Get the response from the client's channel
		select {
		case response := <-client.channel:
			// Parse the response
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")

			jsonStr := response[jsonStart : jsonEnd+1]
			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Response should be valid JSON")

			// Get the result
			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			// Get the content array
			content, ok := result["content"].([]any)
			require.True(t, ok, "Result should have a content array")
			require.NotEmpty(t, content, "Content should not be empty")

			// The first content item should be converted from ToolResult to ImageContent
			firstItem, ok := content[0].(map[string]any)
			require.True(t, ok, "First content item should be an object")

			// Check image data and mime type
			data, ok := firstItem["data"].(string)
			require.True(t, ok, "Content should have data")
			assert.Equal(t, "dGVzdCBpbWFnZSBkYXRhIGZvciB0b29sIHJlc3VsdA==", data, "Image data should match")

			mimeType, ok := firstItem["mimeType"].(string)
			require.True(t, ok, "Content should have mimeType")
			assert.Equal(t, "image/jpeg", mimeType, "MimeType should match")

		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for tool call response")
		}
	})

	// Test image ToolResult
	t.Run("audioToolResult", func(t *testing.T) {
		// Create a client
		client := addTestClient(mock.server, "test-client-image", true)

		// Create a tool call request
		params := struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}{
			Name:      "toolresult.audio.tool",
			Arguments: map[string]any{},
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      2,
			Method:  methodToolsCall,
			Params:  paramBytes,
		}

		// Process the tool call
		mock.server.processToolCall(context.Background(), client, req)

		// Get the response from the client's channel
		select {
		case response := <-client.channel:
			// Parse the response
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")

			jsonStr := response[jsonStart : jsonEnd+1]
			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Response should be valid JSON")

			// Get the result
			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			// Get the content array
			content, ok := result["content"].([]any)
			require.True(t, ok, "Result should have a content array")
			require.NotEmpty(t, content, "Content should not be empty")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for tool call response")
		}
	})

	// Test text ToolResult
	t.Run("ToolResult with int type", func(t *testing.T) {
		// Create a client
		client := addTestClient(mock.server, "test-client-text", true)

		// Create a tool call request
		params := struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}{
			Name:      "toolresult.int.tool",
			Arguments: map[string]any{},
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      1,
			Method:  methodToolsCall,
			Params:  paramBytes,
		}

		// Process the tool call
		mock.server.processToolCall(context.Background(), client, req)

		// Get the response from the client's channel
		select {
		case response := <-client.channel:
			// Parse the response
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")

			jsonStr := response[jsonStart : jsonEnd+1]
			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Response should be valid JSON")

			// Get the result
			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			// Get the content array
			content, ok := result["content"].([]any)
			require.True(t, ok, "Result should have a content array")
			require.NotEmpty(t, content, "Content should not be empty")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for tool call response")
		}
	})

	// Test text ToolResult
	t.Run("ToolResult with bad type", func(t *testing.T) {
		// Create a client
		client := addTestClient(mock.server, "test-client-text", true)

		// Create a tool call request
		params := struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}{
			Name:      "toolresult.bad.tool",
			Arguments: map[string]any{},
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      1,
			Method:  methodToolsCall,
			Params:  paramBytes,
		}

		// Process the tool call
		mock.server.processToolCall(context.Background(), client, req)

		// Get the response from the client's channel
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "json: unsupported type")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for tool call response")
		}
	})
}

// TestToolCallError tests that tool errors are properly handled
func TestToolCallError(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Register a tool that returns an error
	err := mock.server.RegisterTool(Tool{
		Name:        "error.tool",
		Description: "A tool that returns an error",
		InputSchema: InputSchema{
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			return nil, fmt.Errorf("tool execution failed")
		},
	})
	require.NoError(t, err)

	// Simulate a client
	client := addTestClient(mock.server, "test-client", true)

	// Create a tool call request
	params := struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}{
		Name:      "error.tool",
		Arguments: map[string]any{},
	}

	paramBytes, _ := json.Marshal(params)
	req := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  methodToolsCall,
		Params:  paramBytes,
	}

	// Process the tool call
	mock.server.processToolCall(context.Background(), client, req)

	// Check the response
	select {
	case response := <-client.channel:
		assert.Contains(t, response, "event: message", "Response should have message event")
		assert.Contains(t, response, "Error:", "Response should contain the error message")
		assert.Contains(t, response, "isError", "Response should indicate it's an error")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tool call response")
	}
}

// TestToolCallTimeout tests that tool timeouts are properly handled
func TestToolCallTimeout(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Register a tool that times out
	err := mock.server.RegisterTool(Tool{
		Name:        "timeout.tool",
		Description: "A tool that times out",
		InputSchema: InputSchema{
			Properties: map[string]any{},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			<-ctx.Done()
			return nil, fmt.Errorf("tool execution timed out")
		},
	})
	require.NoError(t, err)

	// Simulate a client
	client := addTestClient(mock.server, "test-client", true)

	// Create a tool call request
	params := struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}{
		Name:      "timeout.tool",
		Arguments: map[string]any{},
	}

	paramBytes, _ := json.Marshal(params)
	req := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  methodToolsCall,
		Params:  paramBytes,
	}
	jsonBody, _ := json.Marshal(req)

	// Create HTTP request
	r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client", bytes.NewReader(jsonBody))
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancel()
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	// Process through handleRequest
	go mock.server.handleRequest(w, r)

	// Check the response
	select {
	case response := <-client.channel:
		assert.Contains(t, response, "event: message", "Response should have message event")
		assert.Contains(t, response, `-32001`, "Response should contain a timeout error code")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tool call response")
	}
}

// TestInitializeAndNotifications tests the client initialization flow
func TestInitializeAndNotifications(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Create a test client
	client := addTestClient(mock.server, "test-client", false)

	// Test initialize request
	initReq := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  "initialize",
		Params:  json.RawMessage(`{}`),
	}

	mock.server.processInitialize(context.Background(), client, initReq)

	// Check that client is initialized after initialize request
	assert.True(t, client.initialized, "Client should be marked as initialized after initialize request")

	// Check the response format
	select {
	case response := <-client.channel:
		// Check response contains required initialization fields
		assert.Contains(t, response, "protocolVersion", "Response should include protocol version")
		assert.Contains(t, response, "capabilities", "Response should include capabilities")
		assert.Contains(t, response, "serverInfo", "Response should include server info")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for initialize response")
	}

	// Test notification initialized
	mock.server.processNotificationInitialized(client)
	assert.True(t, client.initialized, "Client should remain initialized after notification")
}

// TestRequestHandlingWithoutInitialization tests that requests are properly rejected when client is not initialized
func TestRequestHandlingWithoutInitialization(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	mock.registerExampleTool()

	// Create an uninitialized test client
	client := addTestClient(mock.server, "test-client", false)

	// Attempt a tool call before initialization
	params := struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}{
		Name:      "test.tool",
		Arguments: map[string]any{"input": "foo"},
	}

	paramBytes, _ := json.Marshal(params)
	req := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  methodToolsCall,
		Params:  paramBytes,
	}
	jsonBody, _ := json.Marshal(req)

	// Handle the request - this should fail because client is not initialized
	r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client", bytes.NewReader(jsonBody))
	mock.server.handleRequest(httptest.NewRecorder(), r)

	// Check error response
	select {
	case response := <-client.channel:
		assert.Contains(t, strings.ToLower(response), "error", "Response should contain an error")
		assert.Contains(t, strings.ToLower(response), "not fully initialized",
			"Response should mention client not being initialized")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for error response")
	}
}

// TestPing tests the ping request handling
func TestPing(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Create a test client
	client := addTestClient(mock.server, "test-client", true)

	// Create a ping request
	pingReq := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  "ping",
		Params:  json.RawMessage(`{}`),
	}

	jsonBody, _ := json.Marshal(pingReq)

	// Handle the request - this should fail because client is not initialized
	r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client", bytes.NewReader(jsonBody))
	mock.server.handleRequest(httptest.NewRecorder(), r)

	// Check response
	select {
	case response := <-client.channel:
		assert.Contains(t, response, `"result":`, "Response should contain a result field")
		assert.Contains(t, response, `"id":1`, "Response should have the same ID as the request")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for ping response")
	}
}

// TestNotificationCancelled tests the notification cancelled handling
func TestNotificationCancelled(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Create a test client
	client := addTestClient(mock.server, "test-client", true)

	// Create a cancellation request
	paramBytes, _ := json.Marshal(map[string]any{
		"requestId": 123,
		"reason":    "user_cancelled",
	})

	cancelReq := Request{
		JsonRpc: jsonRpcVersion,
		Method:  "notifications/cancelled",
		Params:  paramBytes,
	}

	jsonBody, _ := json.Marshal(cancelReq)

	// Handle the request - this should fail because client is not initialized
	r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client", bytes.NewReader(jsonBody))
	mock.server.handleRequest(httptest.NewRecorder(), r)

	// No response expected for notifications
	select {
	case <-client.channel:
		t.Fatal("No response expected for notifications")
	case <-time.After(50 * time.Millisecond):
		// This is the expected outcome - no response
	}
}

func TestNotificationCancelled_badParams(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	client := addTestClient(mock.server, "test-client", true)

	cancelReq := Request{
		JsonRpc: jsonRpcVersion,
		Method:  "notifications/cancelled",
		Params:  []byte(`invalid json`),
	}

	buf := logtest.NewCollector(t)
	mock.server.processNotificationCancelled(context.Background(), client, cancelReq)

	select {
	case <-client.channel:
		t.Fatal("No response expected for notifications")
	case <-time.After(50 * time.Millisecond):
		assert.Contains(t, buf.String(), "Failed to parse cancellation params")
	}
}

func TestUnknownRequest(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Create a test client
	client := addTestClient(mock.server, "test-client", true)
	req := Request{
		JsonRpc: jsonRpcVersion,
		Method:  "unknown",
	}

	jsonBody, _ := json.Marshal(req)

	// Handle the request - this should fail because client is not initialized
	r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client", bytes.NewReader(jsonBody))
	mock.server.handleRequest(httptest.NewRecorder(), r)

	// No response expected for notifications
	select {
	case message := <-client.channel:
		evt, err := parseEvent(message)
		require.NoError(t, err, "Should parse event without error")
		errCode := evt.Data["error"].(map[string]any)["code"]
		// because error code will be converted into float64
		assert.Equal(t, float64(errCodeMethodNotFound), errCode)
	case <-time.After(50 * time.Millisecond):
		// This is the expected outcome - no response
	}
}

func TestResponseWriter_notFlusher(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Handle the request - this should fail because client is not initialized
	r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client", http.NoBody)
	var w notFlusherResponseWriter
	mock.server.handleSSE(&w, r)
	assert.Equal(t, http.StatusInternalServerError, w.code)
}

func TestResponseWriter_cantWrite(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Handle the request - this should fail because client is not initialized
	r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client", http.NoBody)
	var w cantWriteResponseWriter
	mock.server.handleSSE(&w, r)
	assert.Equal(t, 0, w.code)
}

func TestHandleSSE_channelClosed(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Handle the request - this should fail because client is not initialized
	r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client", http.NoBody)
	w := httptest.NewRecorder()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		mock.server.handleSSE(w, r)
		wg.Done()
	}()

	buf := logtest.NewCollector(t)
	for {
		time.Sleep(time.Millisecond)
		mock.server.clientsLock.Lock()
		if len(mock.server.clients) > 0 {
			for _, client := range mock.server.clients {
				close(client.channel)
				delete(mock.server.clients, client.id)
			}
			mock.server.clientsLock.Unlock()
			break
		}
		mock.server.clientsLock.Unlock()
	}
	wg.Wait()
	assert.Contains(t, "channel was closed", buf.Content(), "Should log channel closed error")
}

func TestHandleSSE_badResponseWriter(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Handle the request - this should fail because client is not initialized
	r := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		var w writeOnceResponseWriter
		mock.server.handleSSE(&w, r)
		wg.Done()
	}()

	var session string
	for {
		time.Sleep(time.Millisecond)
		mock.server.clientsLock.Lock()
		if len(mock.server.clients) > 0 {
			for _, client := range mock.server.clients {
				session = client.id
			}
			mock.server.clientsLock.Unlock()
			break
		}
		mock.server.clientsLock.Unlock()
	}

	time.Sleep(100 * time.Millisecond)
	// Create a ping request
	pingReq := Request{
		JsonRpc: jsonRpcVersion,
		ID:      1,
		Method:  "ping",
		Params:  json.RawMessage(`{}`),
	}

	jsonBody, _ := json.Marshal(pingReq)
	buf := logtest.NewCollector(t)

	// Handle the request - this should fail because client is not initialized
	r = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/?session_id=%s", session),
		bytes.NewReader(jsonBody))
	mock.server.handleRequest(httptest.NewRecorder(), r)

	wg.Wait()
	assert.Contains(t, "Failed to write", buf.Content())
}

// TestGetPrompt tests the prompts/get endpoint
func TestGetPrompt(t *testing.T) {
	t.Run("test prompt", func(t *testing.T) {
		mock := newMockMcpServer(t)
		defer mock.shutdown()

		// Create a test client
		client := addTestClient(mock.server, "test-client", true)

		// Register a test prompt
		testPrompt := Prompt{
			Name:        "test.prompt",
			Description: "A test prompt",
		}
		mock.server.RegisterPrompt(testPrompt)

		// Create a get prompt request
		paramBytes, _ := json.Marshal(map[string]any{
			"name": "test.prompt",
			"arguments": map[string]string{
				"topic": "test topic",
			},
		})

		promptReq := Request{
			JsonRpc: jsonRpcVersion,
			ID:      1,
			Method:  "prompts/get",
			Params:  paramBytes,
		}

		// Process the request
		mock.server.processGetPrompt(context.Background(), client, promptReq)

		// Check response
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "description", "Response should include prompt description")
			assert.Contains(t, response, "prompts", "Response should include prompts array")
			assert.Contains(t, response, "A test prompt", "Response should include the topic argument")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for prompt response")
		}
	})

	t.Run("test prompt with invalid params", func(t *testing.T) {
		mock := newMockMcpServer(t)
		defer mock.shutdown()

		// Create a test client
		client := addTestClient(mock.server, "test-client", true)

		paramBytes := []byte("invalid json")
		promptReq := Request{
			JsonRpc: jsonRpcVersion,
			ID:      1,
			Method:  "prompts/get",
			Params:  paramBytes,
		}

		// Process the request
		mock.server.processGetPrompt(context.Background(), client, promptReq)

		// Check response
		select {
		case response := <-client.channel:
			evt, err := parseEvent(response)
			assert.NoError(t, err, "Should be able to parse event")
			errMsg, ok := evt.Data["error"].(map[string]any)
			assert.True(t, ok, "Should have error in response")
			assert.Equal(t, "Invalid parameters", errMsg["message"], "Error message should match")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for prompt response")
		}
	})

	t.Run("test prompt with nil params", func(t *testing.T) {
		mock := newMockMcpServer(t)
		defer mock.shutdown()

		// Create a test client
		client := addTestClient(mock.server, "test-client", true)

		// Register a test prompt
		testPrompt := Prompt{
			Name:        "test.prompt",
			Description: "A test prompt",
		}
		mock.server.RegisterPrompt(testPrompt)

		// Create a get prompt request
		paramBytes, _ := json.Marshal(map[string]any{
			"name": "test.prompt",
		})
		promptReq := Request{
			JsonRpc: jsonRpcVersion,
			ID:      1,
			Method:  "prompts/get",
			Params:  paramBytes,
		}

		// Process the request
		mock.server.processGetPrompt(context.Background(), client, promptReq)

		// Check response
		select {
		case response := <-client.channel:
			_, err := parseEvent(response)
			assert.NoError(t, err, "Should be able to parse event")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for prompt response")
		}
	})
}

// TestBroadcast tests the broadcast functionality
func TestBroadcast(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Create two test clients
	client1 := addTestClient(mock.server, "client-1", true)
	client2 := addTestClient(mock.server, "client-2", true)

	// Broadcast a test message
	testData := map[string]string{"key": "value"}
	mock.server.broadcast("test_event", testData)

	// Check both clients received the broadcast
	for i, client := range []*mcpClient{client1, client2} {
		select {
		case response := <-client.channel:
			assert.Contains(t, response, `event: test_event`, "Response should have the correct event")
			assert.Contains(t, response, `"key":"value"`, "Response should contain the broadcast data")
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Timed out waiting for broadcast on client %d", i+1)
		}
	}

	buf := logtest.NewCollector(t)
	mock.server.broadcast("test_event", make(chan string))
	// Check that the broadcast was logged
	content := buf.Content()
	assert.Contains(t, content, "Failed", "Broadcast should be logged")

	for i := 0; i < eventChanSize; i++ {
		mock.server.broadcast("test_event", "test")
	}

	done := make(chan struct{})
	go func() {
		mock.server.broadcast("test_event", "ignored")
		close(done)
	}()

	select {
	case <-time.After(100 * time.Millisecond):
		assert.Fail(t, "broadcast should not block")
	case <-done:
	}
}

// TestHandleSSEPing tests the automatic ping functionality in the SSE handler
func TestHandleSSEPing(t *testing.T) {
	originalInterval := pingInterval.Load()
	pingInterval.Set(50 * time.Millisecond)
	defer func() {
		pingInterval.Set(originalInterval)
	}()

	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Create a request context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a test ResponseRecorder and Request
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", mock.server.conf.Mcp.SseEndpoint, nil).WithContext(ctx)

	// Create a channel to coordinate the test
	done := make(chan struct{})

	// Set up a custom ResponseRecorder that captures writes and signals the test
	customResponseWriter := &testResponseWriter{
		ResponseRecorder: w,
		writes:           make([]string, 0),
		done:             done,
		pingDetected:     false,
	}

	// Start the SSE handler in a goroutine
	go func() {
		mock.server.handleSSE(customResponseWriter, r)
	}()

	// Wait for ping or timeout
	select {
	case <-done:
		// A ping was detected
		assert.True(t, customResponseWriter.pingDetected, "Ping message should have been sent")
	case <-time.After(pingInterval.Load() + 100*time.Millisecond):
		t.Fatal("Timed out waiting for ping message")
	}

	// Verify that the client was added and cleaned up
	mock.server.clientsLock.Lock()
	clientCount := len(mock.server.clients)
	mock.server.clientsLock.Unlock()

	// Clean up by cancelling the context
	cancel()

	// Wait for cleanup to complete
	time.Sleep(50 * time.Millisecond)

	// Verify client was removed
	mock.server.clientsLock.Lock()
	finalClientCount := len(mock.server.clients)
	mock.server.clientsLock.Unlock()

	assert.Equal(t, 1, clientCount, "One client should be added during the test")
	assert.Equal(t, 0, finalClientCount, "Client should be cleaned up after context cancelation")
}

// TestHandleSSEPing tests the automatic ping functionality in the SSE handler
func TestHandleSSEPing_writeOnce(t *testing.T) {
	originalInterval := pingInterval.Load()
	pingInterval.Set(50 * time.Millisecond)
	defer func() {
		pingInterval.Set(originalInterval)
	}()

	buf := logtest.NewCollector(t)
	var bufLock sync.Mutex
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Start the SSE handler in a goroutine
	go func() {
		var w writeOnceResponseWriter
		r := httptest.NewRequest(http.MethodGet, mock.server.conf.Mcp.SseEndpoint, http.NoBody)
		bufLock.Lock()
		defer bufLock.Unlock()
		mock.server.handleSSE(&w, r)
	}()

	// Wait for ping or timeout
	time.Sleep(100 * time.Millisecond)
	bufLock.Lock()
	assert.Contains(t, "Failed to send ping", buf.Content())
	bufLock.Unlock()
}

func TestServerStartStop(t *testing.T) {
	// Create a simple configuration for testing
	const yamlConf = `name: test-server
host: localhost
port: 0
timeout: 1000
mcp:
  name: mcp-test-server
`
	var c McpConf
	assert.NoError(t, conf.LoadFromYamlBytes([]byte(yamlConf), &c))

	// Create the server
	s := NewMcpServer(c)

	// Start and stop in goroutine to avoid blocking
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go func() {
		s.Start()
	}()

	// Allow a brief moment for startup
	time.Sleep(50 * time.Millisecond)

	// Stop the server
	s.Stop()

	// Wait for context to ensure we properly stopped or timed out
	<-ctx.Done()
}

// TestNotificationInitialized tests the notifications/initialized handling in detail
func TestNotificationInitialized(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	t.Run("uninitializedClient", func(t *testing.T) {
		// Create an uninitialized test client
		client := addTestClient(mock.server, "test-client-uninitialized", false)
		assert.False(t, client.initialized, "Client should start as uninitialized")

		// Create a notification request
		req := Request{
			JsonRpc: jsonRpcVersion,
			Method:  methodNotificationsInitialized,
			// No ID for notifications
			Params: json.RawMessage(`{}`), // Empty params acceptable for this notification
		}

		// Process through the request handler
		jsonBody, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "/?session_id="+client.id, bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		mock.server.handleRequest(w, r)

		// Verify client is now initialized
		assert.True(t, client.initialized, "Client should be marked as initialized after notifications/initialized")

		// Verify the response code is 202 Accepted
		assert.Equal(t, http.StatusAccepted, w.Code, "Response status should be 202 Accepted")

		// No actual response body should be sent for notifications
		select {
		case <-client.channel:
			t.Fatal("No response expected for notifications")
		case <-time.After(50 * time.Millisecond):
			// This is the expected outcome - no response
		}
	})

	t.Run("initializedClient", func(t *testing.T) {
		// Create an already initialized client
		client := addTestClient(mock.server, "test-client-initialized", true)
		assert.True(t, client.initialized, "Client should start as initialized")

		// Directly call processNotificationInitialized
		mock.server.processNotificationInitialized(client)

		// Verify client remains initialized
		assert.True(t, client.initialized, "Client should remain initialized after notifications/initialized")

		// No response expected
		select {
		case <-client.channel:
			t.Fatal("No response expected for notifications")
		case <-time.After(50 * time.Millisecond):
			// This is the expected outcome - no response
		}
	})

	t.Run("errorOnIncorrectUsage", func(t *testing.T) {
		// Create a test client
		client := addTestClient(mock.server, "test-client-error", false)

		// Create a request with ID (incorrect usage - should be a notification)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      123, // Adding ID makes this an incorrect usage - should be notification
			Method:  methodNotificationsInitialized,
			Params:  json.RawMessage(`{}`),
		}

		// Process through the request handler
		jsonBody, _ := json.Marshal(req)
		r := httptest.NewRequest(http.MethodPost, "/?session_id="+client.id, bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()
		mock.server.handleRequest(w, r)

		// Should get an error response
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain an error")
			assert.Contains(t, response, "Method should be used as a notification", "Response should explain notification usage")
			assert.Contains(t, response, `"id":123`, "Response should include the original ID")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}

		// Client should not be initialized due to error
		assert.False(t, client.initialized, "Client should not be initialized after error")
	})
}

func TestSendResponse(t *testing.T) {
	t.Run("bad response", func(t *testing.T) {
		mock := newMockMcpServer(t)
		defer mock.shutdown()

		// Create a test client
		client := addTestClient(mock.server, "test-client", true)

		// Create a response
		response := Response{
			JsonRpc: jsonRpcVersion,
			ID:      1,
			Result:  make(chan int),
		}

		// Send the response
		mock.server.sendResponse(context.Background(), client, 1, response)

		// Check the response in the client's channel
		select {
		case res := <-client.channel:
			evt, err := parseEvent(res)
			require.NoError(t, err, "Should parse event without error")
			errMsg, ok := evt.Data["error"].(map[string]any)
			require.True(t, ok, "Should have error in response")
			assert.Equal(t, float64(errCodeInternalError), errMsg["code"], "Error code should match")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for response")
		}
	})

	t.Run("channel full", func(t *testing.T) {
		mock := newMockMcpServer(t)
		defer mock.shutdown()

		// Create a test client
		client := addTestClient(mock.server, "test-client", true)
		for i := 0; i < eventChanSize; i++ {
			client.channel <- "test"
		}

		// Create a response
		response := Response{
			JsonRpc: jsonRpcVersion,
			ID:      1,
			Result:  "foo",
		}

		buf := logtest.NewCollector(t)
		// Send the response
		mock.server.sendResponse(context.Background(), client, 1, response)
		// Check the response in the client's channel
		assert.Contains(t, buf.String(), "channel is full")
	})
}

func TestSendErrorResponse(t *testing.T) {
	t.Run("channel full", func(t *testing.T) {
		mock := newMockMcpServer(t)
		defer mock.shutdown()

		// Create a test client
		client := addTestClient(mock.server, "test-client", true)
		for i := 0; i < eventChanSize; i++ {
			client.channel <- "test"
		}

		buf := logtest.NewCollector(t)
		// Send the response
		mock.server.sendErrorResponse(context.Background(), client, 1, "foo", errCodeInternalError)
		// Check the response in the client's channel
		assert.Contains(t, buf.String(), "channel is full")
	})
}

// TestMethodToolsCall tests the handling of tools/call method through handleRequest
func TestMethodToolsCall(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	t.Run("validToolCall", func(t *testing.T) {
		// Register a test tool
		mock.registerExampleTool()

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-valid", true)

		// Create a tools call request with progress token metadata
		params := struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
			Meta      struct {
				ProgressToken string `json:"progressToken"`
			} `json:"_meta"`
		}{
			Name: "test.tool",
			Arguments: map[string]any{
				"input": "test-input",
			},
			Meta: struct {
				ProgressToken string `json:"progressToken"`
			}{
				ProgressToken: "token123",
			},
		}

		paramBytes, err := json.Marshal(params)
		require.NoError(t, err, "Failed to marshal tool call parameters")

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      42, // Specific ID to verify in response
			Method:  methodToolsCall,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-valid", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest (full path)
		mock.server.handleRequest(w, r)

		// Verify the HTTP response
		assert.Equal(t, http.StatusAccepted, w.Code, "HTTP status should be 202 Accepted")

		// Check the response in client's channel
		select {
		case response := <-client.channel:
			// Verify it's a message event with the expected format
			assert.Contains(t, response, "event: message", "Response should be a message event")

			// Parse the JSON part of the SSE message
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")
			jsonStr := response[jsonStart : jsonEnd+1]

			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Should be able to parse response JSON")

			// Validate the structure
			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			// Verify ID matches our request
			id, ok := parsed["id"].(float64)
			assert.True(t, ok, "Response should have an ID")
			assert.Equal(t, float64(42), id, "Response ID should match request ID")

			// Verify content
			content, ok := result["content"].([]any)
			require.True(t, ok, "Result should have content array")
			assert.NotEmpty(t, content, "Content should not be empty")

			// Check for progress token in metadata
			meta, hasMeta := result["_meta"].(map[string]any)
			assert.True(t, hasMeta, "Response should include _meta with progress token")
			if hasMeta {
				token, hasToken := meta["progressToken"].(string)
				assert.True(t, hasToken, "Meta should include progress token")
				assert.Equal(t, "token123", token, "Progress token should match")
			}

			// Check actual result content
			if len(content) > 0 {
				firstItem, ok := content[0].(map[string]any)
				if ok {
					assert.Contains(t, firstItem[ContentTypeText], "Processed: test-input", "Content should include processed input")
				}
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for tool call response")
		}
	})

	t.Run("invalidToolName", func(t *testing.T) {
		// Create an initialized client
		client := addTestClient(mock.server, "test-client-invalid", true)

		// Create a tools call request with invalid tool name
		params := struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}{
			Name:      "non-existent-tool",
			Arguments: map[string]any{},
		}

		paramBytes, err := json.Marshal(params)
		require.NoError(t, err, "Failed to marshal tool call parameters")

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      43,
			Method:  methodToolsCall,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-invalid", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Verify response contains error about non-existent tool
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "not found", "Error should mention tool not found")
			assert.Contains(t, response, "non-existent-tool", "Error should mention the invalid tool name")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("clientNotInitialized", func(t *testing.T) {
		// Register a tool
		mock.registerExampleTool()

		// Create an uninitialized client
		client := addTestClient(mock.server, "test-client-uninitialized", false)

		// Create a valid tools call request
		params := struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}{
			Name: "test.tool",
			Arguments: map[string]any{
				"input": "test-input",
			},
		}

		paramBytes, err := json.Marshal(params)
		require.NoError(t, err, "Failed to marshal tool call parameters")

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      44,
			Method:  methodToolsCall,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-uninitialized", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Verify response contains error about client not being initialized
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "not fully initialized", "Error should mention client not initialized")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})
}

// TestMethodPromptsGet tests the handling of prompts/get method through handleRequest
func TestMethodPromptsGet(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	t.Run("staticPrompt", func(t *testing.T) {
		// Register a test prompt with static content
		testPrompt := Prompt{
			Name:        "static-prompt",
			Description: "A static test prompt with placeholders",
			Arguments: []PromptArgument{
				{
					Name:        "name",
					Description: "Name to use in greeting",
					Required:    true,
				},
				{
					Name:        "topic",
					Description: "Topic to discuss",
				},
			},
			Content: "Hello {{name}}! Let's talk about {{topic}}.",
		}
		mock.server.RegisterPrompt(testPrompt)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-static", true)

		// Create a prompts/get request
		params := struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}{
			Name: "static-prompt",
			Arguments: map[string]string{
				"name": "Test User",
				// Intentionally not providing "topic" to test default values
			},
		}

		paramBytes, err := json.Marshal(params)
		require.NoError(t, err, "Failed to marshal prompt get parameters")

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      70,
			Method:  methodPromptsGet,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-static", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest (full path)
		mock.server.handleRequest(w, r)

		// Verify the HTTP response
		assert.Equal(t, http.StatusAccepted, w.Code, "HTTP status should be 202 Accepted")

		// Check the response in client's channel
		select {
		case response := <-client.channel:
			// Verify it's a message event with the expected format
			assert.Contains(t, response, "event: message", "Response should be a message event")

			// Parse the JSON part of the SSE message
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")
			jsonStr := response[jsonStart : jsonEnd+1]

			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Should be able to parse response JSON")

			// Validate the structure
			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			// Verify ID matches our request
			id, ok := parsed["id"].(float64)
			assert.True(t, ok, "Response should have an ID")
			assert.Equal(t, float64(70), id, "Response ID should match request ID")

			// Verify description
			description, ok := result["description"].(string)
			assert.True(t, ok, "Response should include prompt description")
			assert.Equal(t, "A static test prompt with placeholders", description, "Description should match")

			// Verify messages
			messages, ok := result["messages"].([]any)
			require.True(t, ok, "Result should have messages array")
			assert.Len(t, messages, 1, "Should have 1 message")

			// Check message content - should have placeholder substitutions
			if len(messages) > 0 {
				message, ok := messages[0].(map[string]any)
				require.True(t, ok, "Message should be an object")
				assert.Equal(t, string(RoleUser), message["role"], "Role should be 'user'")

				content, ok := message["content"].(map[string]any)
				require.True(t, ok, "Should have content object")
				assert.Equal(t, ContentTypeText, content["type"], "Content type should be text")
				assert.Contains(t, content[ContentTypeText], "Hello Test User", "Content should include the name argument")
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for prompt get response")
		}
	})

	t.Run("dynamicPrompt", func(t *testing.T) {
		// Register a test prompt with a handler function
		testPrompt := Prompt{
			Name:        "dynamic-prompt",
			Description: "A dynamic test prompt with a handler",
			Arguments: []PromptArgument{
				{
					Name:        "username",
					Description: "User's name",
					Required:    true,
				},
				{
					Name:        "question",
					Description: "User's question",
				},
			},
			Handler: func(ctx context.Context, args map[string]string) ([]PromptMessage, error) {
				username := args["username"]
				question := args["question"]

				// Create a system message
				systemMessage := PromptMessage{
					Role: RoleAssistant,
					Content: TextContent{
						Text: "You are a helpful assistant.",
					},
				}

				// Create a user message
				userMessage := PromptMessage{
					Role: RoleUser,
					Content: TextContent{
						Text: fmt.Sprintf("Hi, I'm %s and I'm wondering: %s", username, question),
					},
				}

				return []PromptMessage{systemMessage, userMessage}, nil
			},
		}
		mock.server.RegisterPrompt(testPrompt)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-dynamic", true)

		// Create a prompts/get request
		params := struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}{
			Name: "dynamic-prompt",
			Arguments: map[string]string{
				"username": "Dynamic User",
				"question": "How to test this?",
			},
		}

		paramBytes, err := json.Marshal(params)
		require.NoError(t, err, "Failed to marshal prompt get parameters")

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      71,
			Method:  methodPromptsGet,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-dynamic", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check the response
		select {
		case response := <-client.channel:
			// Extract and parse JSON
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			jsonStr := response[jsonStart : jsonEnd+1]

			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Should be able to parse response JSON")

			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			// Verify messages - should have 2 messages from handler
			messages, ok := result["messages"].([]any)
			require.True(t, ok, "Result should have messages array")
			assert.Len(t, messages, 2, "Should have 2 messages")

			// Check message content
			if len(messages) >= 2 {
				// First message should be assistant
				message1, _ := messages[0].(map[string]any)
				assert.Equal(t, string(RoleAssistant), message1["role"], "First role should be 'system'")

				content1, _ := message1["content"].(map[string]any)
				assert.Contains(t, content1[ContentTypeText], "helpful assistant", "System message should be correct")

				// Second message should be user
				message2, _ := messages[1].(map[string]any)
				assert.Equal(t, string(RoleUser), message2["role"], "Second role should be 'user'")

				content2, _ := message2["content"].(map[string]any)
				assert.Contains(t, content2[ContentTypeText], "Dynamic User", "User message should contain username")
				assert.Contains(t, content2[ContentTypeText], "How to test this?", "User message should contain question")
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for prompt get response")
		}
	})

	t.Run("missingRequiredArgument", func(t *testing.T) {
		// Register a test prompt with a required argument
		testPrompt := Prompt{
			Name:        "required-arg-prompt",
			Description: "A prompt with required arguments",
			Arguments: []PromptArgument{
				{
					Name:        "required_arg",
					Description: "This argument is required",
					Required:    true,
				},
			},
		}
		mock.server.RegisterPrompt(testPrompt)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-missing-arg", true)

		// Create a prompts/get request with missing required argument
		params := struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}{
			Name:      "required-arg-prompt",
			Arguments: map[string]string{}, // Empty arguments
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      72,
			Method:  methodPromptsGet,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-missing-arg", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check for error response about missing required argument
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "Missing required arguments", "Error should mention missing arguments")
			assert.Contains(t, response, "required_arg", "Error should name the missing argument")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("promptNotFound", func(t *testing.T) {
		// Create an initialized client
		client := addTestClient(mock.server, "test-client-prompt-not-found", true)

		// Create a prompts/get request with non-existent prompt
		params := struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}{
			Name:      "non-existent-prompt",
			Arguments: map[string]string{},
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      73,
			Method:  methodPromptsGet,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-prompt-not-found", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check for error response about non-existent prompt
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "Prompt 'non-existent-prompt' not found", "Error should mention prompt not found")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("handlerError", func(t *testing.T) {
		// Register a test prompt with a handler that returns an error
		testPrompt := Prompt{
			Name:        "error-handler-prompt",
			Description: "A prompt with a handler that returns an error",
			Arguments:   []PromptArgument{},
			Handler: func(ctx context.Context, args map[string]string) ([]PromptMessage, error) {
				return nil, fmt.Errorf("test handler error")
			},
		}
		mock.server.RegisterPrompt(testPrompt)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-handler-error", true)

		// Create a prompts/get request
		params := struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}{
			Name:      "error-handler-prompt",
			Arguments: map[string]string{},
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      74,
			Method:  methodPromptsGet,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-handler-error", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check for error response about handler error
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "Error generating prompt content", "Error should mention generating content")
			assert.Contains(t, response, "test handler error", "Error should include the handler error message")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("invalidParameters", func(t *testing.T) {
		// Create an invalid JSON request
		invalidJson := []byte(`{"not valid json`)

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      75,
			Method:  methodPromptsGet,
			Params:  invalidJson,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-invalid-params", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("clientNotInitialized", func(t *testing.T) {
		// Register a basic prompt
		testPrompt := Prompt{
			Name:        "basic-prompt",
			Description: "A basic test prompt",
		}
		mock.server.RegisterPrompt(testPrompt)

		// Create an uninitialized client
		client := addTestClient(mock.server, "test-client-uninit", false)

		// Create a valid prompts/get request
		params := struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}{
			Name:      "basic-prompt",
			Arguments: map[string]string{},
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      76,
			Method:  methodPromptsGet,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-uninit", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Verify response contains error about client not being initialized
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "not fully initialized", "Error should mention client not initialized")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})
}

func TestMethodResourcesList(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	t.Run("validResourceWithHandler", func(t *testing.T) {
		// Register a test resource with handler
		testResource := Resource{
			Name:        "test-resource",
			URI:         "file:///test/resource.txt",
			Description: "A test resource with handler",
			MimeType:    "text/plain",
			Handler: func(ctx context.Context) (ResourceContent, error) {
				return ResourceContent{
					URI:      "file:///test/resource.txt",
					MimeType: "text/plain",
					Text:     "This is test resource content",
				}, nil
			},
		}
		mock.server.RegisterResource(testResource)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-resources", true)

		// Create a resources/read request
		params := PaginatedParams{
			Cursor: "next-cursor",
			Meta: struct {
				ProgressToken any `json:"progressToken"`
			}{
				ProgressToken: "token",
			},
		}

		paramBytes, err := json.Marshal(params)
		require.NoError(t, err, "Failed to marshal resource read parameters")

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      50,
			Method:  methodResourcesList,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-resources", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest (full path)
		mock.server.handleRequest(w, r)

		// Verify the HTTP response
		assert.Equal(t, http.StatusAccepted, w.Code, "HTTP status should be 202 Accepted")

		// Check the response in client's channel
		select {
		case response := <-client.channel:
			evt, err := parseEvent(response)
			assert.NoError(t, err)
			result, ok := evt.Data["result"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, "next-cursor", result["nextCursor"])
			assert.Equal(t, "token", result["_meta"].(map[string]any)["progressToken"])
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for resource read response")
		}
	})
}

// TestMethodResourcesRead tests the handling of resources/read method
func TestMethodResourcesRead(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	t.Run("validResourceWithHandler", func(t *testing.T) {
		// Register a test resource with handler
		testResource := Resource{
			Name:        "test-resource",
			URI:         "file:///test/resource.txt",
			Description: "A test resource with handler",
			MimeType:    "text/plain",
			Handler: func(ctx context.Context) (ResourceContent, error) {
				return ResourceContent{
					URI:      "file:///test/resource.txt",
					MimeType: "text/plain",
					Text:     "This is test resource content",
				}, nil
			},
		}
		mock.server.RegisterResource(testResource)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-resources", true)

		// Create a resources/read request
		params := ResourceReadParams{
			URI: "file:///test/resource.txt",
		}

		paramBytes, err := json.Marshal(params)
		require.NoError(t, err, "Failed to marshal resource read parameters")

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      50,
			Method:  methodResourcesRead,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-resources", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest (full path)
		mock.server.handleRequest(w, r)

		// Verify the HTTP response
		assert.Equal(t, http.StatusAccepted, w.Code, "HTTP status should be 202 Accepted")

		// Check the response in client's channel
		select {
		case response := <-client.channel:
			// Verify it's a message event with the expected format
			assert.Contains(t, response, "event: message", "Response should be a message event")

			// Parse the JSON part of the SSE message
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")
			jsonStr := response[jsonStart : jsonEnd+1]

			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Should be able to parse response JSON")

			// Validate the structure
			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			// Verify ID matches our request
			id, ok := parsed["id"].(float64)
			assert.True(t, ok, "Response should have an ID")
			assert.Equal(t, float64(50), id, "Response ID should match request ID")

			// Verify contents
			contents, ok := result["contents"].([]any)
			require.True(t, ok, "Result should have contents array")
			assert.Len(t, contents, 1, "Contents array should have 1 item")

			// Check content details
			if len(contents) > 0 {
				content, ok := contents[0].(map[string]any)
				require.True(t, ok, "Content should be an object")
				assert.Equal(t, "file:///test/resource.txt", content["uri"], "URI should match")
				assert.Equal(t, "text/plain", content["mimeType"], "MimeType should match")
				assert.Equal(t, "This is test resource content", content[ContentTypeText], "Text content should match")
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for resource read response")
		}
	})

	t.Run("resourceWithoutHandler", func(t *testing.T) {
		// Register a test resource without handler
		testResource := Resource{
			Name:        "no-handler-resource",
			URI:         "file:///test/no-handler.txt",
			Description: "A test resource without handler",
			MimeType:    "text/plain",
			// No handler provided
		}
		mock.server.RegisterResource(testResource)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-no-handler", true)

		// Create a resources/read request
		params := ResourceReadParams{
			URI: "file:///test/no-handler.txt",
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      51,
			Method:  methodResourcesRead,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-no-handler", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check for response with empty content
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "event: message", "Response should be a message event")

			// Extract and parse JSON
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			jsonStr := response[jsonStart : jsonEnd+1]

			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Should be able to parse response JSON")

			// Check contents exists but has empty text
			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			contents, ok := result["contents"].([]any)
			require.True(t, ok, "Result should have contents array")
			assert.Len(t, contents, 1, "Contents array should have 1 item")

			// Check content details - should have URI and MimeType but empty text
			if len(contents) > 0 {
				content, ok := contents[0].(map[string]any)
				require.True(t, ok, "Content should be an object")
				assert.Equal(t, "file:///test/no-handler.txt", content["uri"], "URI should match")
				assert.Equal(t, "text/plain", content["mimeType"], "MimeType should match")
				_, ok = content[ContentTypeText]
				assert.False(t, ok, "Text content should be empty string")
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for resource read response")
		}
	})

	t.Run("resourceNotFound", func(t *testing.T) {
		// Create an initialized client
		client := addTestClient(mock.server, "test-client-not-found", true)

		// Create a resources/read request with non-existent URI
		params := ResourceReadParams{
			URI: "file:///test/non-existent.txt",
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      52,
			Method:  methodResourcesRead,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-not-found", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check for error response about resource not found
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "Resource with URI", "Error should mention resource URI")
			assert.Contains(t, response, "not found", "Error should indicate resource not found")
			assert.Contains(t, response, "file:///test/non-existent.txt", "Error should mention the requested URI")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("invalidParameters", func(t *testing.T) {
		// Create an invalid JSON request
		invalidJson := []byte(`{"not valid json`)

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      53,
			Method:  methodResourcesRead,
			Params:  invalidJson,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-invalid-params", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalidParameters direct", func(t *testing.T) {
		// Create an invalid JSON request
		invalidJson := []byte(`{"not valid json`)

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      53,
			Method:  methodResourcesRead,
			Params:  invalidJson,
		}

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-resources", true)

		// Process through handleRequest
		mock.server.processResourcesRead(context.Background(), client, req)

		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "Invalid parameters", "Error should mention invalid parameters")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("handlerError", func(t *testing.T) {
		// Register a test resource with handler that returns an error
		testResource := Resource{
			Name:        "error-resource",
			URI:         "file:///test/error.txt",
			Description: "A test resource with handler that returns error",
			MimeType:    "text/plain",
			Handler: func(ctx context.Context) (ResourceContent, error) {
				return ResourceContent{}, fmt.Errorf("test handler error")
			},
		}
		mock.server.RegisterResource(testResource)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-handler-error", true)

		// Create a resources/read request
		params := ResourceReadParams{
			URI: "file:///test/error.txt",
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      54,
			Method:  methodResourcesRead,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-handler-error", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check for error response about handler error
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "Error reading resource", "Error should mention reading resource")
			assert.Contains(t, response, "test handler error", "Error should include handler error message")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("handlerMissingURIAndMimeType", func(t *testing.T) {
		// Register a test resource with handler that returns content without URI and MimeType
		testResource := Resource{
			Name:        "missing-fields-resource",
			URI:         "file:///test/missing-fields.txt",
			Description: "A test resource with handler that returns content missing fields",
			MimeType:    "text/plain",
			Handler: func(ctx context.Context) (ResourceContent, error) {
				// Return ResourceContent without URI and MimeType
				return ResourceContent{
					Text: "Content with missing fields",
				}, nil
			},
		}
		mock.server.RegisterResource(testResource)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-missing-fields", true)

		// Create a resources/read request
		params := ResourceReadParams{
			URI: "file:///test/missing-fields.txt",
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      55,
			Method:  methodResourcesRead,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-missing-fields", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check response - server should fill in the missing fields
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "event: message", "Response should be a message event")

			// Extract and parse JSON
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			jsonStr := response[jsonStart : jsonEnd+1]

			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Should be able to parse response JSON")

			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")

			contents, ok := result["contents"].([]any)
			require.True(t, ok, "Result should have contents array")
			assert.Len(t, contents, 1, "Contents array should have 1 item")

			// Check content details - server should fill in missing URI and MimeType
			if len(contents) > 0 {
				content, ok := contents[0].(map[string]any)
				require.True(t, ok, "Content should be an object")
				assert.Equal(t, "file:///test/missing-fields.txt", content["uri"], "URI should be filled from request")
				assert.Equal(t, "text/plain", content["mimeType"], "MimeType should be filled from resource")
				assert.Equal(t, "Content with missing fields", content[ContentTypeText], "Text content should match")
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for resource read response")
		}
	})
}

// TestMethodResourcesSubscribe tests the handling of resources/subscribe method
func TestMethodResourcesSubscribe(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	t.Run("validSubscription", func(t *testing.T) {
		// Register a test resource
		testResource := Resource{
			Name:        "subscribe-resource",
			URI:         "file:///test/subscribe.txt",
			Description: "A test resource for subscription",
			MimeType:    "text/plain",
		}
		mock.server.RegisterResource(testResource)

		// Create an initialized client
		client := addTestClient(mock.server, "test-client-subscribe", true)

		// Create a resources/subscribe request
		params := ResourceSubscribeParams{
			URI: "file:///test/subscribe.txt",
		}

		paramBytes, err := json.Marshal(params)
		require.NoError(t, err, "Failed to marshal resource subscribe parameters")

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      60,
			Method:  methodResourcesSubscribe,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-subscribe", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest (full path)
		mock.server.handleRequest(w, r)

		// Verify the HTTP response
		assert.Equal(t, http.StatusAccepted, w.Code, "HTTP status should be 202 Accepted")

		// Check the response in client's channel - should be an empty success response
		select {
		case response := <-client.channel:
			// Verify it's a message event with the expected format
			assert.Contains(t, response, "event: message", "Response should be a message event")

			// Parse the JSON part of the SSE message
			jsonStart := strings.Index(response, "{")
			jsonEnd := strings.LastIndex(response, "}")
			require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")
			jsonStr := response[jsonStart : jsonEnd+1]

			var parsed map[string]any
			err := json.Unmarshal([]byte(jsonStr), &parsed)
			require.NoError(t, err, "Should be able to parse response JSON")

			// Verify ID matches our request
			id, ok := parsed["id"].(float64)
			assert.True(t, ok, "Response should have an ID")
			assert.Equal(t, float64(60), id, "Response ID should match request ID")

			// Verify the result exists and is an empty object
			result, ok := parsed["result"].(map[string]any)
			require.True(t, ok, "Response should have a result object")
			assert.Empty(t, result, "Result should be an empty object for successful subscription")

		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for subscription response")
		}
	})

	t.Run("resourceNotFound", func(t *testing.T) {
		// Create an initialized client
		client := addTestClient(mock.server, "test-client-sub-not-found", true)

		// Create a resources/subscribe request with non-existent URI
		params := ResourceSubscribeParams{
			URI: "file:///test/non-existent-subscription.txt",
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      61,
			Method:  methodResourcesSubscribe,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-sub-not-found", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check for error response about resource not found
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "Resource with URI", "Error should mention resource URI")
			assert.Contains(t, response, "not found", "Error should indicate resource not found")
			assert.Contains(t, response, "file:///test/non-existent-subscription.txt", "Error should mention the requested URI")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("invalidParameters", func(t *testing.T) {
		// Create an invalid JSON request
		invalidJson := []byte(`{"not valid json`)

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      62,
			Method:  methodResourcesSubscribe,
			Params:  invalidJson,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-sub-invalid-params", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Code, "HTTP status should be 400 Bad Request")
	})

	t.Run("invalidParameters direct", func(t *testing.T) {
		// Create an invalid JSON request
		invalidJson := []byte(`{"not valid json`)

		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      62,
			Method:  methodResourcesSubscribe,
			Params:  invalidJson,
		}

		client := addTestClient(mock.server, "test-client-sub-not-found", true)
		mock.server.processResourceSubscribe(context.Background(), client, req)

		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "Invalid parameters", "Error should mention invalid parameters")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("clientNotInitialized", func(t *testing.T) {
		// Register a test resource
		testResource := Resource{
			Name:        "subscribe-resource-uninit",
			URI:         "file:///test/subscribe-uninit.txt",
			Description: "A test resource for subscription with uninitialized client",
		}
		mock.server.RegisterResource(testResource)

		// Create an uninitialized client
		client := addTestClient(mock.server, "test-client-sub-uninitialized", false)

		// Create a valid resources/subscribe request
		params := ResourceSubscribeParams{
			URI: "file:///test/subscribe-uninit.txt",
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      63,
			Method:  methodResourcesSubscribe,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-sub-uninitialized", bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Verify response contains error about client not being initialized
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "not fully initialized", "Error should mention client not initialized")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})

	t.Run("missingURIParameter", func(t *testing.T) {
		// Create an initialized client
		client := addTestClient(mock.server, "test-client-sub-missing-uri", true)

		// Create a subscription request with empty URI
		params := ResourceSubscribeParams{
			URI: "", // Empty URI
		}

		paramBytes, _ := json.Marshal(params)
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      64,
			Method:  methodResourcesSubscribe,
			Params:  paramBytes,
		}
		jsonBody, _ := json.Marshal(req)

		// Create HTTP request
		r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-sub-missing-uri",
			bytes.NewReader(jsonBody))
		w := httptest.NewRecorder()

		// Process through handleRequest
		mock.server.handleRequest(w, r)

		// Check for error response about resource not found (empty URI)
		select {
		case response := <-client.channel:
			assert.Contains(t, response, "error", "Response should contain error")
			assert.Contains(t, response, "Resource with URI", "Error should mention resource URI")
			assert.Contains(t, response, "not found", "Error should indicate resource not found")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timed out waiting for error response")
		}
	})
}

// TestToolCallUnmarshalError tests the error handling when unmarshaling invalid JSON in processToolCall
func TestToolCallUnmarshalError(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Create an initialized client
	client := addTestClient(mock.server, "test-client-unmarshal-error", true)

	// Create a request with invalid JSON in Params
	req := Request{
		JsonRpc: "2.0",
		ID:      100,
		Method:  methodToolsCall,
		Params:  []byte(`{"name": "test.tool", "arguments": {"input": invalid_json}}`), // This is invalid JSON
	}

	// Process the tool call directly
	mock.server.processToolCall(context.Background(), client, req)

	// Check for error response about invalid JSON
	select {
	case response := <-client.channel:
		assert.Contains(t, response, "error", "Response should contain an error")
		assert.Contains(t, response, "Invalid tool call parameters", "Error should mention invalid parameters")

		// Extract error code from response
		jsonStart := strings.Index(response, "{")
		jsonEnd := strings.LastIndex(response, "}")
		require.True(t, jsonStart >= 0 && jsonEnd > jsonStart, "Response should contain valid JSON")
		jsonStr := response[jsonStart : jsonEnd+1]

		var parsed struct {
			Error struct {
				Code int `json:"code"`
			} `json:"error"`
		}
		err := json.Unmarshal([]byte(jsonStr), &parsed)
		require.NoError(t, err, "Should be able to parse response JSON")

		// Verify correct error code was returned
		assert.Equal(t, errCodeInvalidParams, parsed.Error.Code, "Error code should be errCodeInvalidParams")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for error response")
	}
}

// TestToolCallWithInvalidParams tests the handling when calling handleRequest with invalid JSON params
func TestToolCallWithInvalidParams(t *testing.T) {
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Register a tool to make sure it exists
	mock.registerExampleTool()

	// Create a request with invalid JSON
	req := Request{
		JsonRpc: "2.0",
		ID:      101,
		Method:  methodToolsCall,
		Params:  []byte(`{"name": "test.tool", "arguments": {this_is_invalid_json}}`),
	}

	jsonBody, _ := json.Marshal(req)

	// Create HTTP request
	r := httptest.NewRequest(http.MethodPost, "/?session_id=test-client-invalid-json", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	// Process through handleRequest
	mock.server.handleRequest(w, r)

	// Verify HTTP status is Accepted (even for errors, we accept the request)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

type mockResponseWriter struct {
}

func (m *mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockResponseWriter) Write(i []byte) (int, error) {
	return len(i), nil
}

func (m *mockResponseWriter) WriteHeader(_ int) {
}

type notFlusherResponseWriter struct {
	mockResponseWriter
	code int
}

func (m *notFlusherResponseWriter) WriteHeader(code int) {
	m.code = code
}

type cantWriteResponseWriter struct {
	mockResponseWriter
	code int
}

func (m *cantWriteResponseWriter) Flush() {
}

func (m *cantWriteResponseWriter) Write(_ []byte) (int, error) {
	return 0, fmt.Errorf("can't write")
}

type writeOnceResponseWriter struct {
	mockResponseWriter
	times int32
}

func (m *writeOnceResponseWriter) Flush() {
}

func (m *writeOnceResponseWriter) Write(i []byte) (int, error) {
	if atomic.AddInt32(&m.times, 1) > 1 {
		return 0, fmt.Errorf("write once")
	}
	return len(i), nil
}

// testResponseWriter is a custom http.ResponseWriter that captures writes and detects ping messages
type testResponseWriter struct {
	*httptest.ResponseRecorder
	writes       []string
	mu           sync.Mutex
	pingDetected bool
	done         chan struct{}
}

// Write overrides the ResponseRecorder's Write method to detect ping messages
func (w *testResponseWriter) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	written, err := w.ResponseRecorder.Write(b)
	if err != nil {
		return written, err
	}

	content := string(b)
	w.writes = append(w.writes, content)

	// Check if this is a ping message
	if strings.Contains(content, "event: ping") {
		w.pingDetected = true
		// Signal that we've detected a ping
		select {
		case w.done <- struct{}{}:
		default:
			// Channel might be closed or already signaled
		}
	}

	return written, nil
}

// Flush implements the http.Flusher interface
func (w *testResponseWriter) Flush() {
	w.ResponseRecorder.Flush()
}
