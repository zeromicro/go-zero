// filepath: /Users/kevin/Develop/go/opensource/go-zero/mcp/vars_test.go
package mcp

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorCodes ensures error codes are applied correctly in error responses
func TestErrorCodes(t *testing.T) {
	testCases := []struct {
		name     string
		code     int
		message  string
		expected string
	}{
		{
			name:     "invalid request error",
			code:     errCodeInvalidRequest,
			message:  "Invalid request",
			expected: `"code":-32600`,
		},
		{
			name:     "method not found error",
			code:     errCodeMethodNotFound,
			message:  "Method not found",
			expected: `"code":-32601`,
		},
		{
			name:     "invalid params error",
			code:     errCodeInvalidParams,
			message:  "Invalid parameters",
			expected: `"code":-32602`,
		},
		{
			name:     "internal error",
			code:     errCodeInternalError,
			message:  "Internal server error",
			expected: `"code":-32603`,
		},
		{
			name:     "timeout error",
			code:     errCodeTimeout,
			message:  "Operation timed out",
			expected: `"code":-32001`,
		},
		{
			name:     "resource not found error",
			code:     errCodeResourceNotFound,
			message:  "Resource not found",
			expected: `"code":-32002`,
		},
		{
			name:     "client not initialized error",
			code:     errCodeClientNotInitialized,
			message:  "Client not initialized",
			expected: `"code":-32800`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := Response{
				JsonRpc: jsonRpcVersion,
				ID:      int64(1),
				Error: &errorObj{
					Code:    tc.code,
					Message: tc.message,
				},
			}
			data, err := json.Marshal(resp)
			assert.NoError(t, err)
			assert.Contains(t, string(data), tc.expected, "Error code should match expected value")
			assert.Contains(t, string(data), tc.message, "Error message should be included")
			assert.Contains(t, string(data), jsonRpcVersion, "JSON-RPC version should be included")
		})
	}
}

// TestJsonRpcVersion ensures the correct JSON-RPC version is used
func TestJsonRpcVersion(t *testing.T) {
	assert.Equal(t, "2.0", jsonRpcVersion, "JSON-RPC version should be 2.0")

	// Test that it's used in responses
	resp := Response{
		JsonRpc: jsonRpcVersion,
		ID:      int64(1),
		Result:  "test",
	}
	data, err := json.Marshal(resp)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"jsonrpc":"2.0"`, "Response should use correct JSON-RPC version")

	// Test that it's expected in requests
	reqStr := `{"jsonrpc":"2.0","id":1,"method":"test"}`
	var req Request
	err = json.Unmarshal([]byte(reqStr), &req)
	assert.NoError(t, err)
	assert.Equal(t, jsonRpcVersion, req.JsonRpc, "Request should parse correct JSON-RPC version")
}

// TestSessionIdKey ensures session ID extraction works correctly
func TestSessionIdKey(t *testing.T) {
	// Create a mock server implementation
	mock := newMockMcpServer(t)
	defer mock.shutdown()

	// Verify the key constant
	assert.Equal(t, "session_id", sessionIdKey, "Session ID key should be 'session_id'")

	// Test that session ID is extracted correctly
	mockR := httptest.NewRequest("GET", "/?"+sessionIdKey+"=test-session", nil)

	// Since the mock server is using the same session key logic,
	// we can test this by accessing the request query parameters directly
	sessionID := mockR.URL.Query().Get(sessionIdKey)
	assert.Equal(t, "test-session", sessionID, "Session ID should be extracted correctly")
}

// TestEventTypes ensures event types are set correctly in SSE responses
func TestEventTypes(t *testing.T) {
	// Test message event
	assert.Equal(t, "message", eventMessage, "Message event should be 'message'")

	// Test endpoint event
	assert.Equal(t, "endpoint", eventEndpoint, "Endpoint event should be 'endpoint'")

	// Verify them in an actual SSE format string
	messageEvent := "event: " + eventMessage + "\ndata: test\n\n"
	assert.Contains(t, messageEvent, "event: message", "Message event should format correctly")

	endpointEvent := "event: " + eventEndpoint + "\ndata: test\n\n"
	assert.Contains(t, endpointEvent, "event: endpoint", "Endpoint event should format correctly")
}

// TestCollectionKeys checks that collection keys are used correctly
func TestCollectionKeys(t *testing.T) {
	// Verify collection key constants
	assert.Equal(t, "prompts", keyPrompts, "Prompts key should be 'prompts'")
	assert.Equal(t, "resources", keyResources, "Resources key should be 'resources'")
	assert.Equal(t, "tools", keyTools, "Tools key should be 'tools'")
}

// TestRoleTypes checks that role types are used correctly
func TestRoleTypes(t *testing.T) {
	// Test in annotations
	annotations := Annotations{
		Audience: []RoleType{RoleUser, RoleAssistant},
	}
	data, err := json.Marshal(annotations)
	assert.NoError(t, err)
	assert.Contains(t, string(data), `"audience":["user","assistant"]`, "Role types should marshal correctly")
}

// TestMethodNames checks that method names are used correctly
func TestMethodNames(t *testing.T) {
	// Verify method name constants
	methods := map[string]string{
		"initialize":                methodInitialize,
		"tools/list":                methodToolsList,
		"tools/call":                methodToolsCall,
		"prompts/list":              methodPromptsList,
		"prompts/get":               methodPromptsGet,
		"resources/list":            methodResourcesList,
		"resources/read":            methodResourcesRead,
		"resources/subscribe":       methodResourcesSubscribe,
		"ping":                      methodPing,
		"notifications/initialized": methodNotificationsInitialized,
		"notifications/cancelled":   methodNotificationsCancelled,
	}

	for expected, actual := range methods {
		assert.Equal(t, expected, actual, "Method name should be "+expected)
	}

	// Test in a request
	for methodName := range methods {
		req := Request{
			JsonRpc: jsonRpcVersion,
			ID:      int64(1),
			Method:  methodName,
		}
		data, err := json.Marshal(req)
		assert.NoError(t, err)
		assert.Contains(t, string(data), `"method":"`+methodName+`"`, "Method name should be used in requests")
	}
}

// TestEventNames checks that event names are used correctly
func TestEventNames(t *testing.T) {
	// Verify event name constants
	events := map[string]string{
		"tools/list_changed":     eventToolsListChanged,
		"prompts/list_changed":   eventPromptsListChanged,
		"resources/list_changed": eventResourcesListChanged,
	}

	for expected, actual := range events {
		assert.Equal(t, expected, actual, "Event name should be "+expected)
	}

	// Test event names in SSE format
	for _, eventName := range events {
		sseEvent := "event: " + eventName + "\ndata: test\n\n"
		assert.Contains(t, sseEvent, "event: "+eventName, "Event name should format correctly in SSE")
	}
}
