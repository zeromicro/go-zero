package mcp

import (
	"time"

	"github.com/zeromicro/go-zero/rest"
)

// McpConf defines the configuration for an MCP server.
// It embeds rest.RestConf for HTTP server settings
// and adds MCP-specific configuration options.
type McpConf struct {
	rest.RestConf
	Mcp struct {
		// Name is the server name reported in initialize responses
		Name string `json:",optional"`

		// Version is the server version reported in initialize responses
		Version string `json:",default=1.0.0"`

		// UseStreamable when true uses Streamable HTTP transport (2025-03-26 spec),
		// otherwise uses SSE transport (2024-11-05 spec)
		UseStreamable bool `json:",default=false"`

		// SseEndpoint is the path for Server-Sent Events connections
		// Used for SSE transport mode
		SseEndpoint string `json:",default=/sse"`

		// MessageEndpoint is the path for JSON-RPC requests
		// Used for Streamable HTTP transport mode
		MessageEndpoint string `json:",default=/message"`

		// Cors contains allowed CORS origins
		Cors []string `json:",optional"`

		// SseTimeout is the maximum time allowed for SSE connections
		SseTimeout time.Duration `json:",default=24h"`

		// MessageTimeout is the maximum time allowed for request execution
		MessageTimeout time.Duration `json:",default=30s"`
	}
}
