package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// syncResponseRecorder is a thread-safe wrapper around httptest.ResponseRecorder
type syncResponseRecorder struct {
	*httptest.ResponseRecorder
	mu sync.Mutex
}

// Create a new synchronized response recorder
func newSyncResponseRecorder() *syncResponseRecorder {
	return &syncResponseRecorder{
		ResponseRecorder: httptest.NewRecorder(),
	}
}

// Override Write method to synchronize access
func (srr *syncResponseRecorder) Write(p []byte) (int, error) {
	srr.mu.Lock()
	defer srr.mu.Unlock()
	return srr.ResponseRecorder.Write(p)
}

// Override WriteHeader method to synchronize access
func (srr *syncResponseRecorder) WriteHeader(statusCode int) {
	srr.mu.Lock()
	defer srr.mu.Unlock()
	srr.ResponseRecorder.WriteHeader(statusCode)
}

// Override Result method to synchronize access
func (srr *syncResponseRecorder) Result() *http.Response {
	srr.mu.Lock()
	defer srr.mu.Unlock()
	return srr.ResponseRecorder.Result()
}

// TestHTTPHandlerIntegration tests the HTTP handlers with a real server instance
func TestHTTPHandlerIntegration(t *testing.T) {
	// Skip in short test mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a test configuration
	conf := McpConf{}
	conf.Mcp.Name = "test-integration"
	conf.Mcp.Version = "1.0.0-test"
	conf.Mcp.MessageTimeout = 1 * time.Second

	// Create a mock server directly
	server := &sseMcpServer{
		conf:      conf,
		clients:   make(map[string]*mcpClient),
		tools:     make(map[string]Tool),
		prompts:   make(map[string]Prompt),
		resources: make(map[string]Resource),
	}

	// Register a test tool
	err := server.RegisterTool(Tool{
		Name:        "echo",
		Description: "Echo tool for testing",
		InputSchema: InputSchema{
			Properties: map[string]any{
				"message": map[string]any{
					"type":        "string",
					"description": "Message to echo",
				},
			},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			if msg, ok := params["message"].(string); ok {
				return fmt.Sprintf("Echo: %s", msg), nil
			}
			return "Echo: no message provided", nil
		},
	})
	require.NoError(t, err)

	// Create a test HTTP request to the SSE endpoint
	req := httptest.NewRequest("GET", "/sse", nil)
	w := newSyncResponseRecorder()

	// Create a done channel to signal completion of test
	done := make(chan bool)

	// Start the SSE handler in a goroutine
	go func() {
		// lock.Lock()
		server.handleSSE(w, req)
		// lock.Unlock()
		done <- true
	}()

	// Allow time for the handler to process
	select {
	case <-time.After(100 * time.Millisecond):
		// Expected - handler would normally block indefinitely
	case <-done:
		// This shouldn't happen immediately - the handler should block
		t.Error("SSE handler returned unexpectedly")
	}

	// Check the initial headers
	resp := w.Result()
	assert.Equal(t, "chunked", resp.Header.Get("Transfer-Encoding"))
	resp.Body.Close()

	// The handler creates a client and sends the endpoint message
	var sessionId string

	// Give the handler time to set up the client
	time.Sleep(50 * time.Millisecond)

	// Check that a client was created
	server.clientsLock.Lock()
	assert.Equal(t, 1, len(server.clients))
	for id := range server.clients {
		sessionId = id
	}
	server.clientsLock.Unlock()

	require.NotEmpty(t, sessionId, "Expected a session ID to be created")

	// Now that we have a session ID, we can test the message endpoint
	messageBody, _ := json.Marshal(Request{
		JsonRpc: "2.0",
		ID:      1,
		Method:  methodInitialize,
		Params:  json.RawMessage(`{}`),
	})

	// Create a message request
	reqURL := fmt.Sprintf("/message?%s=%s", sessionIdKey, sessionId)
	msgReq := httptest.NewRequest("POST", reqURL, bytes.NewReader(messageBody))
	msgW := newSyncResponseRecorder()

	// Process the message
	server.handleRequest(msgW, msgReq)

	// Check the response
	msgResp := msgW.Result()
	assert.Equal(t, http.StatusAccepted, msgResp.StatusCode)
	msgResp.Body.Close() // Ensure response body is closed
}

// TestHandlerResponseFlow tests the flow of a full request/response cycle
func TestHandlerResponseFlow(t *testing.T) {
	// Create a mock server for testing
	server := &sseMcpServer{
		conf: McpConf{},
		clients: map[string]*mcpClient{
			"test-session": {
				id:          "test-session",
				channel:     make(chan string, 10),
				initialized: true,
			},
		},
		tools:     make(map[string]Tool),
		prompts:   make(map[string]Prompt),
		resources: make(map[string]Resource),
	}

	// Register test resources
	server.RegisterTool(Tool{
		Name:        "test.tool",
		Description: "Test tool",
		InputSchema: InputSchema{Type: "object"},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			return "tool result", nil
		},
	})

	server.RegisterPrompt(Prompt{
		Name:        "test.prompt",
		Description: "Test prompt",
	})

	server.RegisterResource(Resource{
		Name:        "test.resource",
		URI:         "http://example.com",
		Description: "Test resource",
	})

	// Create a request with session ID parameter
	reqURL := fmt.Sprintf("/message?%s=%s", sessionIdKey, "test-session")

	// Test tools/list request
	toolsListBody, _ := json.Marshal(Request{
		JsonRpc: "2.0",
		ID:      1,
		Method:  methodToolsList,
		Params:  json.RawMessage(`{}`),
	})

	toolsReq := httptest.NewRequest("POST", reqURL, bytes.NewReader(toolsListBody))
	toolsW := newSyncResponseRecorder()

	// Process the request
	server.handleRequest(toolsW, toolsReq)

	// Check the response code
	toolsResp := toolsW.Result()
	assert.Equal(t, http.StatusAccepted, toolsResp.StatusCode)
	toolsResp.Body.Close()

	// Check the channel message
	client := server.clients["test-session"]
	select {
	case message := <-client.channel:
		assert.Contains(t, message, `"tools":[{"name":"test.tool"`)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tools/list response")
	}

	// Test prompts/list request
	promptsListBody, _ := json.Marshal(Request{
		JsonRpc: "2.0",
		ID:      2,
		Method:  methodPromptsList,
		Params:  json.RawMessage(`{}`),
	})

	promptsReq := httptest.NewRequest("POST", reqURL, bytes.NewReader(promptsListBody))
	promptsW := newSyncResponseRecorder()

	// Process the request
	server.handleRequest(promptsW, promptsReq)

	// Check the response code
	promptsResp := promptsW.Result()
	assert.Equal(t, http.StatusAccepted, promptsResp.StatusCode)
	promptsResp.Body.Close()

	// Check the channel message
	select {
	case message := <-client.channel:
		assert.Contains(t, message, `"prompts":[{"name":"test.prompt"`)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for prompts/list response")
	}

	// Test resources/list request
	resourcesListBody, _ := json.Marshal(Request{
		JsonRpc: "2.0",
		ID:      3,
		Method:  methodResourcesList,
		Params:  json.RawMessage(`{}`),
	})

	resourcesReq := httptest.NewRequest("POST", reqURL, bytes.NewReader(resourcesListBody))
	resourcesW := newSyncResponseRecorder()

	// Process the request
	server.handleRequest(resourcesW, resourcesReq)

	// Check the response code
	resourcesResp := resourcesW.Result()
	assert.Equal(t, http.StatusAccepted, resourcesResp.StatusCode)
	resourcesResp.Body.Close()

	// Check the channel message
	select {
	case message := <-client.channel:
		assert.Contains(t, message, `"name":"test.resource"`)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for resources/list response")
	}
}

// TestProcessListMethods tests the list processing methods with pagination
func TestProcessListMethods(t *testing.T) {
	server := &sseMcpServer{
		tools:     make(map[string]Tool),
		prompts:   make(map[string]Prompt),
		resources: make(map[string]Resource),
	}

	// Add some test data
	for i := 1; i <= 5; i++ {
		tool := Tool{
			Name:        fmt.Sprintf("tool%d", i),
			Description: fmt.Sprintf("Tool %d", i),
			InputSchema: InputSchema{Type: "object"},
		}
		server.tools[tool.Name] = tool

		prompt := Prompt{
			Name:        fmt.Sprintf("prompt%d", i),
			Description: fmt.Sprintf("Prompt %d", i),
		}
		server.prompts[prompt.Name] = prompt

		resource := Resource{
			Name:        fmt.Sprintf("resource%d", i),
			URI:         fmt.Sprintf("http://example.com/%d", i),
			Description: fmt.Sprintf("Resource %d", i),
		}
		server.resources[resource.Name] = resource
	}

	// Create a test client
	client := &mcpClient{
		id:          "test-client",
		channel:     make(chan string, 10),
		initialized: true,
	}

	// Test processListTools
	req := Request{
		JsonRpc: "2.0",
		ID:      1,
		Method:  methodToolsList,
		Params:  json.RawMessage(`{"cursor": "", "_meta": {"progressToken": "token1"}}`),
	}

	server.processListTools(context.Background(), client, req)

	// Read response
	select {
	case response := <-client.channel:
		assert.Contains(t, response, `"tools":`)
		assert.Contains(t, response, `"progressToken":"token1"`)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for tools/list response")
	}

	// Test processListPrompts
	req.ID = 2
	req.Method = methodPromptsList
	req.Params = json.RawMessage(`{"cursor": "next"}`)
	server.processListPrompts(context.Background(), client, req)

	// Read response
	select {
	case response := <-client.channel:
		assert.Contains(t, response, `"prompts":`)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for prompts/list response")
	}

	// Test processListResources
	req.ID = 3
	req.Method = methodResourcesList
	req.Params = json.RawMessage(`{"cursor": "next"}`)
	server.processListResources(context.Background(), client, req)

	// Read response
	select {
	case response := <-client.channel:
		assert.Contains(t, response, `"resources":`)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for resources/list response")
	}
}

// TestErrorResponseHandling tests error handling in the server
func TestErrorResponseHandling(t *testing.T) {
	server := &sseMcpServer{
		tools:     make(map[string]Tool),
		prompts:   make(map[string]Prompt),
		resources: make(map[string]Resource),
	}

	// Create a test client
	client := &mcpClient{
		id:          "test-client",
		channel:     make(chan string, 10),
		initialized: true,
	}

	// Test invalid method
	req := Request{
		JsonRpc: "2.0",
		ID:      1,
		Method:  "invalid_method",
		Params:  json.RawMessage(`{}`),
	}

	// Mock handleRequest by directly calling error handler
	server.sendErrorResponse(context.Background(), client, req.ID, "Method not found", errCodeMethodNotFound)

	// Check response
	select {
	case response := <-client.channel:
		assert.Contains(t, response, `"error":{"code":-32601,"message":"Method not found"}`)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for error response")
	}

	// Test invalid tool
	toolReq := Request{
		JsonRpc: "2.0",
		ID:      2,
		Method:  methodToolsCall,
		Params:  json.RawMessage(`{"name":"non_existent_tool"}`),
	}

	// Call process method directly
	server.processToolCall(context.Background(), client, toolReq)

	// Check response
	select {
	case response := <-client.channel:
		assert.Contains(t, response, `"error":{"code":-32602,"message":"Tool 'non_existent_tool' not found"}`)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for error response")
	}

	// Test invalid prompt
	promptReq := Request{
		JsonRpc: "2.0",
		ID:      3,
		Method:  methodPromptsGet,
		Params:  json.RawMessage(`{"name":"non_existent_prompt"}`),
	}

	// Call process method directly
	server.processGetPrompt(context.Background(), client, promptReq)

	// Check response
	select {
	case response := <-client.channel:
		assert.Contains(t, response, `"error":{"code":-32602,"message":"Prompt 'non_existent_prompt' not found"}`)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timed out waiting for error response")
	}
}
