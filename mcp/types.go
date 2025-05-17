package mcp

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/zeromicro/go-zero/rest"
)

// Cursor is an opaque token used for pagination
type Cursor string

// Request represents a generic MCP request following JSON-RPC 2.0 specification
type Request struct {
	SessionId string          `form:"session_id"` // Session identifier for client tracking
	JsonRpc   string          `json:"jsonrpc"`    // Must be "2.0" per JSON-RPC spec
	ID        int64           `json:"id"`         // Request identifier for matching responses
	Method    string          `json:"method"`     // Method name to invoke
	Params    json.RawMessage `json:"params"`     // Parameters for the method
}

type PaginatedParams struct {
	Cursor string `json:"cursor"`
	Meta   struct {
		ProgressToken any `json:"progressToken"`
	} `json:"_meta"`
}

// Result is the base interface for all results
type Result struct {
	Meta map[string]any `json:"_meta,omitempty"` // Optional metadata
}

// PaginatedResult is a base for results that support pagination
type PaginatedResult struct {
	Result
	NextCursor Cursor `json:"nextCursor,omitempty"` // Opaque token for fetching next page
}

// ListToolsResult represents the response to a tools/list request
type ListToolsResult struct {
	PaginatedResult
	Tools []Tool `json:"tools"` // List of available tools
}

// Message Content Types

// RoleType represents the sender or recipient of messages in a conversation
type RoleType string

// PromptArgument defines a single argument that can be passed to a prompt
type PromptArgument struct {
	Name        string `json:"name"`                  // Argument name
	Description string `json:"description,omitempty"` // Human-readable description
	Required    bool   `json:"required,omitempty"`    // Whether this argument is required
}

// PromptHandler is a function that dynamically generates prompt content
type PromptHandler func(ctx context.Context, args map[string]string) ([]PromptMessage, error)

// Prompt represents an MCP Prompt definition
type Prompt struct {
	Name        string           `json:"name"`                  // Unique identifier for the prompt
	Description string           `json:"description,omitempty"` // Human-readable description
	Arguments   []PromptArgument `json:"arguments,omitempty"`   // Arguments for customization
	Content     string           `json:"-"`                     // Static content (internal use only)
	Handler     PromptHandler    `json:"-"`                     // Handler for dynamic content generation
}

// PromptMessage represents a message in a conversation
type PromptMessage struct {
	Role    RoleType `json:"role"`    // Message sender role
	Content any      `json:"content"` // Message content (TextContent, ImageContent, etc.)
}

// TextContent represents text content in a message
type TextContent struct {
	Text        string       `json:"text"`                  // The text content
	Annotations *Annotations `json:"annotations,omitempty"` // Optional annotations
}

type typedTextContent struct {
	Type string `json:"type"`
	TextContent
}

// ImageContent represents image data in a message
type ImageContent struct {
	Data     string `json:"data"`     // Base64-encoded image data
	MimeType string `json:"mimeType"` // MIME type (e.g., "image/png")
}

type typedImageContent struct {
	Type string `json:"type"`
	ImageContent
}

// AudioContent represents audio data in a message
type AudioContent struct {
	Data     string `json:"data"`     // Base64-encoded audio data
	MimeType string `json:"mimeType"` // MIME type (e.g., "audio/mp3")
}

type typedAudioContent struct {
	Type string `json:"type"`
	AudioContent
}

// FileContent represents file content
type FileContent struct {
	URI      string `json:"uri"`      // URI identifying the file
	MimeType string `json:"mimeType"` // MIME type of the file
	Text     string `json:"text"`     // File content as text
}

// EmbeddedResource represents a resource embedded in a message
type EmbeddedResource struct {
	Type     string          `json:"type"`     // Always "resource"
	Resource ResourceContent `json:"resource"` // The resource data
}

// Annotations provides additional metadata for content
type Annotations struct {
	Audience []RoleType `json:"audience,omitempty"` // Who should see this content
	Priority *float64   `json:"priority,omitempty"` // Optional priority (0-1)
}

// Tool-related Types

// ToolHandler is a function that handles tool calls
type ToolHandler func(ctx context.Context, params map[string]any) (any, error)

// Tool represents a Model Context Protocol Tool definition
type Tool struct {
	Name        string      `json:"name"`        // Unique identifier for the tool
	Description string      `json:"description"` // Human-readable description
	InputSchema InputSchema `json:"inputSchema"` // JSON Schema for parameters
	Handler     ToolHandler `json:"-"`           // Not sent to clients
}

// InputSchema represents tool's input schema in JSON Schema format
type InputSchema struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`         // Property definitions
	Required   []string       `json:"required,omitempty"` // List of required properties
}

// CallToolResult represents a tool call result that conforms to the MCP schema
type CallToolResult struct {
	Result
	Content []any `json:"content"`           // Content items (text, images, etc.)
	IsError bool  `json:"isError,omitempty"` // True if tool execution failed
}

// Resource represents a Model Context Protocol Resource definition
type Resource struct {
	URI         string          `json:"uri"`                   // Unique resource identifier (RFC3986)
	Name        string          `json:"name"`                  // Human-readable name
	Description string          `json:"description,omitempty"` // Optional description
	MimeType    string          `json:"mimeType,omitempty"`    // Optional MIME type
	Handler     ResourceHandler `json:"-"`                     // Internal handler not sent to clients
}

// ResourceHandler is a function that handles resource read requests
type ResourceHandler func(ctx context.Context) (ResourceContent, error)

// ResourceContent represents the content of a resource
type ResourceContent struct {
	URI      string `json:"uri"`                // Resource URI (required)
	MimeType string `json:"mimeType,omitempty"` // MIME type of the resource
	Text     string `json:"text,omitempty"`     // Text content (if available)
	Blob     string `json:"blob,omitempty"`     // Base64 encoded blob data (if available)
}

// ResourcesListResult represents the response to a resources/list request
type ResourcesListResult struct {
	PaginatedResult
	Resources []Resource `json:"resources"` // List of available resources
}

// ResourceReadParams contains parameters for a resources/read request
type ResourceReadParams struct {
	URI string `json:"uri"` // URI of the resource to read
}

// ResourceReadResult contains the result of a resources/read request
type ResourceReadResult struct {
	Result
	Contents []ResourceContent `json:"contents"` // Array of resource content
}

// ResourceSubscribeParams contains parameters for a resources/subscribe request
type ResourceSubscribeParams struct {
	URI string `json:"uri"` // URI of the resource to subscribe to
}

// ResourceUpdateNotification represents a notification about a resource update
type ResourceUpdateNotification struct {
	URI     string          `json:"uri"`     // URI of the updated resource
	Content ResourceContent `json:"content"` // New resource content
}

// Client and Server Types

// mcpClient represents an SSE client connection
type mcpClient struct {
	id          string      // Unique client identifier
	channel     chan string // Channel for sending SSE messages
	initialized bool        // Tracks if client has sent notifications/initialized
}

// McpServer defines the interface for Model Context Protocol servers
type McpServer interface {
	Start()
	Stop()
	RegisterTool(tool Tool) error
	RegisterPrompt(prompt Prompt)
	RegisterResource(resource Resource)
}

// sseMcpServer implements the McpServer interface using SSE
type sseMcpServer struct {
	conf          McpConf
	server        *rest.Server
	clients       map[string]*mcpClient
	clientsLock   sync.Mutex
	tools         map[string]Tool
	toolsLock     sync.Mutex
	prompts       map[string]Prompt
	promptsLock   sync.Mutex
	resources     map[string]Resource
	resourcesLock sync.Mutex
}

// Response Types

// errorObj represents a JSON-RPC error object
type errorObj struct {
	Code    int    `json:"code"`    // Error code
	Message string `json:"message"` // Error message
}

// Response represents a JSON-RPC response
type Response struct {
	JsonRpc string    `json:"jsonrpc"`         // Always "2.0"
	ID      int64     `json:"id"`              // Same as request ID
	Result  any       `json:"result"`          // Result object (null if error)
	Error   *errorObj `json:"error,omitempty"` // Error object (null if success)
}

// Server Information Types

// serverInfo provides information about the server
type serverInfo struct {
	Name    string `json:"name"`    // Server name
	Version string `json:"version"` // Server version
}

// capabilities describes the server's capabilities
type capabilities struct {
	Logging struct{} `json:"logging"`
	Prompts struct {
		ListChanged bool `json:"listChanged"` // Server will notify on prompt changes
	} `json:"prompts"`
	Resources struct {
		Subscribe   bool `json:"subscribe"`   // Server supports resource subscriptions
		ListChanged bool `json:"listChanged"` // Server will notify on resource changes
	} `json:"resources"`
	Tools struct {
		ListChanged bool `json:"listChanged"` // Server will notify on tool changes
	} `json:"tools"`
}

// initializationResponse is sent in response to an initialize request
type initializationResponse struct {
	ProtocolVersion string       `json:"protocolVersion"` // Protocol version
	Capabilities    capabilities `json:"capabilities"`    // Server capabilities
	ServerInfo      serverInfo   `json:"serverInfo"`      // Server information
}

// ToolCallParams contains the parameters for a tool call
type ToolCallParams struct {
	Name       string         `json:"name"`       // Tool name
	Parameters map[string]any `json:"parameters"` // Tool parameters
}

// ToolResult contains the result of a tool execution
type ToolResult struct {
	Type    string `json:"type"`    // Content type (text, image, etc.)
	Content any    `json:"content"` // Result content
}

// errorMessage represents a detailed error message
type errorMessage struct {
	Code    int    `json:"code"`       // Error code
	Message string `json:"message"`    // Error message
	Data    any    `json:",omitempty"` // Additional error data
}
