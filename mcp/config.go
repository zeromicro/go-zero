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

		// ProtocolVersion is the MCP protocol version implemented
		ProtocolVersion string `json:",default=2024-11-05"`

		// BaseUrl is the base URL for the server, used in SSE endpoint messages
		// If not set, defaults to http://localhost:{Port}
		BaseUrl string `json:",optional"`

		// SseEndpoint is the path for Server-Sent Events connections
		SseEndpoint string `json:",default=/sse"`

		// MessageEndpoint is the path for JSON-RPC requests
		MessageEndpoint string `json:",default=/message"`

		// Cors contains allowed CORS origins
		Cors []string `json:",optional"`

		// SseTimeout is the maximum time allowed for SSE connections
		SseTimeout time.Duration `json:",default=24h"`

		// MessageTimeout is the maximum time allowed for request execution
		MessageTimeout time.Duration `json:",default=30s"`
	}
}
