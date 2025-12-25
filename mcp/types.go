package mcp

import (
	"context"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Re-export commonly used SDK types for convenience
type (
	// Tool types
	Tool            = sdkmcp.Tool
	CallToolParams  = sdkmcp.CallToolParams
	CallToolResult  = sdkmcp.CallToolResult
	CallToolRequest = sdkmcp.CallToolRequest

	// Content types
	Content      = sdkmcp.Content
	TextContent  = sdkmcp.TextContent
	ImageContent = sdkmcp.ImageContent
	AudioContent = sdkmcp.AudioContent

	// Prompt types
	Prompt          = sdkmcp.Prompt
	PromptMessage   = sdkmcp.PromptMessage
	GetPromptParams = sdkmcp.GetPromptParams
	GetPromptResult = sdkmcp.GetPromptResult

	// Resource types
	Resource           = sdkmcp.Resource
	ResourceContents   = sdkmcp.ResourceContents
	ReadResourceParams = sdkmcp.ReadResourceParams
	ReadResourceResult = sdkmcp.ReadResourceResult

	// Session and server types
	Server         = sdkmcp.Server
	ServerSession  = sdkmcp.ServerSession
	ServerOptions  = sdkmcp.ServerOptions
	Implementation = sdkmcp.Implementation

	// Transport types
	SSEHandler            = sdkmcp.SSEHandler
	StreamableHTTPHandler = sdkmcp.StreamableHTTPHandler
)

// ToolHandler is a generic function signature for tool handlers.
// Handlers should accept context, request, and typed arguments, and return
// a result, metadata, and error.
//
// Deprecated: Use ToolHandlerFor directly from the SDK types.
type ToolHandler[Args any, Meta any] func(
	ctx context.Context,
	req *CallToolRequest,
	args Args,
) (*CallToolResult, Meta, error)

// PromptHandler is a function signature for prompt handlers.
type PromptHandler func(
	ctx context.Context,
	req *sdkmcp.GetPromptRequest,
	args map[string]string,
) (*GetPromptResult, error)

// ResourceHandler is a function signature for resource handlers.
type ResourceHandler func(
	ctx context.Context,
	req *sdkmcp.ReadResourceRequest,
	uri string,
) (*ReadResourceResult, error)

// AddTool registers a tool with the MCP server using type-safe generics.
// The SDK automatically generates JSON schema from the Args struct tags.
//
// Example:
//
//	type GreetArgs struct {
//	    Name string `json:"name" jsonschema:"description=Name to greet"`
//	}
//
//	tool := &mcp.Tool{
//	    Name: "greet",
//	    Description: "Greet someone",
//	}
//
//	handler := func(ctx context.Context, req *mcp.CallToolRequest, args GreetArgs) (*mcp.CallToolResult, any, error) {
//	    return &mcp.CallToolResult{
//	        Content: []mcp.Content{&mcp.TextContent{Text: "Hello " + args.Name}},
//	    }, nil, nil
//	}
//
//	mcp.AddTool(server, tool, handler)
func AddTool[In, Out any](server McpServer, tool *Tool, handler func(context.Context, *CallToolRequest, In) (*CallToolResult, Out, error)) {
	sdkmcp.AddTool(server.Server(), tool, handler)
}
