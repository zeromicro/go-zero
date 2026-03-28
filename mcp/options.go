package mcp

import (
	"net/http"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ServerHook customizes the underlying SDK server before the HTTP server starts.
type ServerHook func(*sdkmcp.Server)

// ServerSelector returns the SDK server that should handle the incoming request.
// Returning nil falls back to the default server created by NewMcpServer.
type ServerSelector func(*http.Request) *sdkmcp.Server

// Option customizes MCP server construction in a backward-compatible way.
type Option func(*serverOptions)

type serverOptions struct {
	sdkOptions     *sdkmcp.ServerOptions
	serverHooks    []ServerHook
	serverSelector ServerSelector
}

// WithServerOptions configures the underlying SDK server options.
func WithServerOptions(opts *sdkmcp.ServerOptions) Option {
	return func(o *serverOptions) {
		o.sdkOptions = opts
	}
}

// WithServerHook applies a hook to the underlying SDK server after creation.
func WithServerHook(hook ServerHook) Option {
	return func(o *serverOptions) {
		if hook != nil {
			o.serverHooks = append(o.serverHooks, hook)
		}
	}
}

// WithServerSelector allows callers to select different SDK server instances
// for different requests while keeping the current default server behavior.
func WithServerSelector(selector ServerSelector) Option {
	return func(o *serverOptions) {
		o.serverSelector = selector
	}
}
