package mcp

import (
	"context"
	"net/http"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

// McpServer defines the interface for Model Context Protocol servers using the official SDK
type McpServer interface {
	// Start starts the HTTP server
	Start()
	// Stop stops the HTTP server
	Stop()
}

type mcpServerImpl struct {
	conf       McpConf
	httpServer *rest.Server
	mcpServer  *sdkmcp.Server
	selector   ServerSelector
	metadata   RequestMetadataExtractor
}

// NewMcpServer creates a new MCP server using the official SDK
func NewMcpServer(c McpConf) McpServer {
	return NewMcpServerWithOptions(c)
}

// NewMcpServerWithOptions creates a new MCP server with optional SDK hooks,
// SDK server options, and request-based server selection.
func NewMcpServerWithOptions(c McpConf, opts ...Option) McpServer {
	var options serverOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&options)
		}
	}

	// Create the underlying rest HTTP server
	var httpServer *rest.Server
	if len(c.Mcp.Cors) == 0 {
		httpServer = rest.MustNewServer(c.RestConf)
	} else {
		httpServer = rest.MustNewServer(c.RestConf, rest.WithCors(c.Mcp.Cors...))
	}

	// Set defaults
	if len(c.Mcp.Name) == 0 {
		c.Mcp.Name = c.Name
	}
	if len(c.Mcp.Version) == 0 {
		c.Mcp.Version = "1.0.0"
	}

	// Create the MCP SDK server
	impl := &sdkmcp.Implementation{
		Name:    c.Mcp.Name,
		Version: c.Mcp.Version,
	}

	mcpServer := sdkmcp.NewServer(impl, options.sdkOptions)
	for _, hook := range options.serverHooks {
		hook(mcpServer)
	}

	s := &mcpServerImpl{
		conf:       c,
		httpServer: httpServer,
		mcpServer:  mcpServer,
		selector:   options.serverSelector,
		metadata:   options.metadata,
	}

	// Choose transport based on configuration
	if c.Mcp.UseStreamable {
		s.setupStreamableTransport()
	} else {
		s.setupSSETransport()
	}

	return s
}

// Start implements McpServer.
func (s *mcpServerImpl) Start() {
	logx.Infof("Starting MCP server %s v%s on %s:%d",
		s.conf.Mcp.Name, s.conf.Mcp.Version, s.conf.Host, s.conf.Port)
	s.httpServer.Start()
}

// Stop implements McpServer.
func (s *mcpServerImpl) Stop() {
	logx.Info("Stopping MCP server")
	s.httpServer.Stop()
}

func (s *mcpServerImpl) selectServer(r *http.Request) *sdkmcp.Server {
	if s.selector != nil {
		if server := s.selector(r); server != nil {
			return server
		}
	}

	return s.mcpServer
}

func (s *mcpServerImpl) sseServer(r *http.Request) *sdkmcp.Server {
	logx.Infof("New SSE connection from %s", r.RemoteAddr)
	return s.selectServer(r)
}

func (s *mcpServerImpl) streamableServer(r *http.Request) *sdkmcp.Server {
	logx.Infof("New streamable connection from %s", r.RemoteAddr)
	return s.selectServer(r)
}

// setupSSETransport configures the server to use SSE transport (2024-11-05 spec)
func (s *mcpServerImpl) setupSSETransport() {
	// Create SSE handler that returns our MCP server for each connection
	handler := sdkmcp.NewSSEHandler(s.sseServer, nil)

	s.registerRoutes(handler, s.conf.Mcp.SseEndpoint)
}

// setupStreamableTransport configures the server to use Streamable HTTP transport (2025-03-26 spec)
func (s *mcpServerImpl) setupStreamableTransport() {
	// Create Streamable HTTP handler
	handler := sdkmcp.NewStreamableHTTPHandler(s.streamableServer, nil)

	s.registerRoutes(handler, s.conf.Mcp.MessageEndpoint)
}

func (s *mcpServerImpl) registerRoutes(handler http.Handler, endpoint string) {
	handler = s.wrapHandler(handler)

	// Register the endpoint (handles both GET for SSE and POST for messages)
	s.httpServer.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    endpoint,
		Handler: handler.ServeHTTP,
	}, rest.WithSSE(), rest.WithTimeout(s.conf.Mcp.SseTimeout))

	s.httpServer.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    endpoint,
		Handler: handler.ServeHTTP,
	}, rest.WithTimeout(s.conf.Mcp.MessageTimeout))
}

func (s *mcpServerImpl) wrapHandler(next http.Handler) http.Handler {
	if s.metadata == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metadata := s.metadata(r)
		if metadata.Path == nil {
			metadata.Path = cloneStringMap(pathvar.Vars(r))
		}

		if isEmptyRequestMetadata(metadata) {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), requestMetadataContextKey{}, cloneRequestMetadata(metadata))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func cloneRequestMetadata(metadata RequestMetadata) RequestMetadata {
	return RequestMetadata{
		Headers: cloneStringMap(metadata.Headers),
		Query:   cloneStringMap(metadata.Query),
		Path:    cloneStringMap(metadata.Path),
	}
}

func cloneStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return nil
	}

	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}

	return cloned
}

func isEmptyRequestMetadata(metadata RequestMetadata) bool {
	return len(metadata.Headers) == 0 && len(metadata.Query) == 0 && len(metadata.Path) == 0
}
