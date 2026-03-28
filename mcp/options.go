package mcp

import (
	"context"
	"net/http"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ServerHook customizes the underlying SDK server before the HTTP server starts.
type ServerHook func(*sdkmcp.Server)

// ServerSelector returns the SDK server that should handle the incoming request.
// Returning nil falls back to the default server created by NewMcpServer.
type ServerSelector func(*http.Request) *sdkmcp.Server

// RequestMetadata contains selected HTTP request metadata made available to MCP handlers.
type RequestMetadata struct {
	Headers map[string]string
	Query   map[string]string
	Path    map[string]string
}

// RequestMetadataExtractor extracts selected metadata from the original HTTP request.
type RequestMetadataExtractor func(*http.Request) RequestMetadata

type requestMetadataContextKey struct{}

// Option customizes MCP server construction in a backward-compatible way.
type Option func(*serverOptions)

type serverOptions struct {
	sdkOptions     *sdkmcp.ServerOptions
	serverHooks    []ServerHook
	serverSelector ServerSelector
	metadata       RequestMetadataExtractor
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

// WithRequestMetadataExtractor extracts selected HTTP request metadata and makes it
// available to tool, prompt, and resource handlers through context helpers.
func WithRequestMetadataExtractor(extractor RequestMetadataExtractor) Option {
	return func(o *serverOptions) {
		o.metadata = extractor
	}
}

// RequestMetadataFromContext returns metadata extracted from the original HTTP request.
func RequestMetadataFromContext(ctx context.Context) RequestMetadata {
	if ctx == nil {
		return RequestMetadata{}
	}

	metadata, _ := ctx.Value(requestMetadataContextKey{}).(RequestMetadata)
	if isEmptyRequestMetadata(metadata) {
		return RequestMetadata{}
	}

	return cloneRequestMetadata(metadata)
}

// HeaderFromContext returns a single extracted header value.
func HeaderFromContext(ctx context.Context, key string) string {
	metadata := RequestMetadataFromContext(ctx)
	if len(metadata.Headers) == 0 {
		return ""
	}

	return metadata.Headers[key]
}

// QueryFromContext returns a single extracted query value.
func QueryFromContext(ctx context.Context, key string) string {
	metadata := RequestMetadataFromContext(ctx)
	if len(metadata.Query) == 0 {
		return ""
	}

	return metadata.Query[key]
}

// PathFromContext returns a single extracted path variable value.
func PathFromContext(ctx context.Context, key string) string {
	metadata := RequestMetadataFromContext(ctx)
	if len(metadata.Path) == 0 {
		return ""
	}

	return metadata.Path[key]
}
