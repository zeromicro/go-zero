package mcp

import "net/http"

// RequestMetadataExtractor extracts request metadata for downstream handlers.
type RequestMetadataExtractor func(*http.Request) RequestMetadata

// McpOption customizes MCP server construction.
type McpOption interface {
	apply(*serverOptions)
}

type mcpOptionFunc func(*serverOptions)

func (f mcpOptionFunc) apply(opts *serverOptions) {
	f(opts)
}

type serverOptions struct {
	requestMetadataExtractor RequestMetadataExtractor
}

func defaultServerOptions() serverOptions {
	return serverOptions{}
}

// WithRequestMetadataExtractor installs an extractor that runs for each incoming
// MCP HTTP request, and stores the extracted metadata into handler context.
func WithRequestMetadataExtractor(extractor RequestMetadataExtractor) McpOption {
	return mcpOptionFunc(func(opts *serverOptions) {
		opts.requestMetadataExtractor = extractor
	})
}
