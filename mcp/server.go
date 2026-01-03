package mcp

import (
	"net/http"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
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
}

// NewMcpServer creates a new MCP server using the official SDK
func NewMcpServer(c McpConf) McpServer {
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

	mcpServer := sdkmcp.NewServer(impl, nil)

	s := &mcpServerImpl{
		conf:       c,
		httpServer: httpServer,
		mcpServer:  mcpServer,
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

// setupSSETransport configures the server to use SSE transport (2024-11-05 spec)
func (s *mcpServerImpl) setupSSETransport() {
	// Create SSE handler that returns our MCP server for each connection
	handler := sdkmcp.NewSSEHandler(func(r *http.Request) *sdkmcp.Server {
		logx.Infof("New SSE connection from %s", r.RemoteAddr)
		return s.mcpServer
	}, nil)

	// Register the SSE endpoint
	s.httpServer.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    s.conf.Mcp.SseEndpoint,
		Handler: handler.ServeHTTP,
	}, rest.WithSSE(), rest.WithTimeout(s.conf.Mcp.SseTimeout))

	// The SSE handler also handles POST requests to message endpoints
	// We need to route those as well
	s.httpServer.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    s.conf.Mcp.SseEndpoint,
		Handler: handler.ServeHTTP,
	}, rest.WithTimeout(s.conf.Mcp.MessageTimeout))
}

// setupStreamableTransport configures the server to use Streamable HTTP transport (2025-03-26 spec)
func (s *mcpServerImpl) setupStreamableTransport() {
	// Create Streamable HTTP handler
	handler := sdkmcp.NewStreamableHTTPHandler(func(r *http.Request) *sdkmcp.Server {
		logx.Infof("New streamable connection from %s", r.RemoteAddr)
		return s.mcpServer
	}, nil)

	// Register the message endpoint (handles both GET for SSE and POST for messages)
	s.httpServer.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    s.conf.Mcp.MessageEndpoint,
		Handler: handler.ServeHTTP,
	}, rest.WithSSE(), rest.WithTimeout(s.conf.Mcp.SseTimeout))

	s.httpServer.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    s.conf.Mcp.MessageEndpoint,
		Handler: handler.ServeHTTP,
	}, rest.WithTimeout(s.conf.Mcp.MessageTimeout))
}
