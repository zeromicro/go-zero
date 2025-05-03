package mcp

import (
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
)

// Protocol constants
const (
	// JSON-RPC version as defined in the specification
	jsonRpcVersion = "2.0"

	// Session identifier key used in request URLs
	sessionIdKey = "session_id"

	// progressTokenKey is used to track progress of long-running tasks
	progressTokenKey = "progressToken"
)

// Server-Sent Events (SSE) event types
const (
	// Standard message event for JSON-RPC responses
	eventMessage = "message"

	// Endpoint event for sending endpoint URL to clients
	eventEndpoint = "endpoint"
)

// Content type identifiers
const (
	// ContentTypeObject is object content type
	ContentTypeObject = "object"

	// ContentTypeText is text content type
	ContentTypeText = "text"

	// ContentTypeImage is image content type
	ContentTypeImage = "image"

	// ContentTypeAudio is audio content type
	ContentTypeAudio = "audio"

	// ContentTypeResource is resource content type
	ContentTypeResource = "resource"
)

// Collection keys for broadcast events
const (
	// Key for prompts collection
	keyPrompts = "prompts"

	// Key for resources collection
	keyResources = "resources"

	// Key for tools collection
	keyTools = "tools"
)

// JSON-RPC error codes
// Standard error codes from JSON-RPC 2.0 spec
const (
	// Invalid JSON was received by the server
	errCodeInvalidRequest = -32600

	// The method does not exist / is not available
	errCodeMethodNotFound = -32601

	// Invalid method parameter(s)
	errCodeInvalidParams = -32602

	// Internal JSON-RPC error
	errCodeInternalError = -32603

	// Tool execution timed out
	errCodeTimeout = -32001

	// Resource not found error
	errCodeResourceNotFound = -32002

	// Client hasn't completed initialization
	errCodeClientNotInitialized = -32800
)

// User and assistant role definitions
const (
	// RoleUser is the "user" role - the entity asking questions
	RoleUser RoleType = "user"

	// RoleAssistant is the "assistant" role - the entity providing responses
	RoleAssistant RoleType = "assistant"
)

// Method names as defined in the MCP specification
const (
	// Initialize the connection between client and server
	methodInitialize = "initialize"

	// List available tools
	methodToolsList = "tools/list"

	// Call a specific tool
	methodToolsCall = "tools/call"

	// List available prompts
	methodPromptsList = "prompts/list"

	// Get a specific prompt
	methodPromptsGet = "prompts/get"

	// List available resources
	methodResourcesList = "resources/list"

	// Read a specific resource
	methodResourcesRead = "resources/read"

	// Subscribe to resource updates
	methodResourcesSubscribe = "resources/subscribe"

	// Simple ping to check server availability
	methodPing = "ping"

	// Notification that client is fully initialized
	methodNotificationsInitialized = "notifications/initialized"

	// Notification that a request was canceled
	methodNotificationsCancelled = "notifications/cancelled"
)

// Event names for Server-Sent Events (SSE)
const (
	// Notification of tool list changes
	eventToolsListChanged = "tools/list_changed"

	// Notification of prompt list changes
	eventPromptsListChanged = "prompts/list_changed"

	// Notification of resource list changes
	eventResourcesListChanged = "resources/list_changed"
)

var (
	// Default channel size for events
	eventChanSize = 10

	// Default ping interval for checking connection availability
	// use syncx.ForAtomicDuration to ensure atomicity in test race
	pingInterval = syncx.ForAtomicDuration(30 * time.Second)
)
