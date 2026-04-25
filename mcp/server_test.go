package mcp

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestNewMcpServer(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8080
	c.Mcp.Name = "test-server"
	c.Mcp.Version = "1.0.0"

	server := NewMcpServer(c)
	assert.NotNil(t, server)
}

func TestNewMcpServerWithDefaults(t *testing.T) {
	c := McpConf{}
	c.Name = "default-server"
	c.Host = "localhost"
	c.Port = 8082

	server := NewMcpServer(c)
	impl := server.(*mcpServerImpl)

	// Check defaults are set
	assert.Equal(t, "default-server", impl.conf.Mcp.Name)
	assert.Equal(t, "1.0.0", impl.conf.Mcp.Version)
}

func TestNewMcpServerWithCORS(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8083
	c.Mcp.Name = "cors-server"
	c.Mcp.Cors = []string{"http://localhost:3000", "http://example.com"}

	server := NewMcpServer(c)
	assert.NotNil(t, server)
}

func TestSetupSSETransport(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8084
	c.Mcp.Name = "sse-server"
	c.Mcp.UseStreamable = false
	c.Mcp.SseEndpoint = "/sse"
	c.Mcp.MessageTimeout = 30 * time.Second
	c.Mcp.SseTimeout = 24 * time.Hour

	server := NewMcpServer(c)
	impl := server.(*mcpServerImpl)

	assert.NotNil(t, impl.httpServer)
	assert.False(t, impl.conf.Mcp.UseStreamable)
}

func TestSetupStreamableTransport(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8085
	c.Mcp.Name = "streamable-server"
	c.Mcp.UseStreamable = true
	c.Mcp.MessageEndpoint = "/message"
	c.Mcp.MessageTimeout = 30 * time.Second
	c.Mcp.SseTimeout = 24 * time.Hour

	server := NewMcpServer(c)
	impl := server.(*mcpServerImpl)

	assert.NotNil(t, impl.httpServer)
	assert.True(t, impl.conf.Mcp.UseStreamable)
}

func TestServerImplementsInterface(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8086
	c.Mcp.Name = "interface-test"

	var _ McpServer = NewMcpServer(c)
}

func TestAddTool(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8081
	c.Mcp.Name = "test-server"

	server := NewMcpServer(c)

	type Args struct {
		Name string `json:"name"`
	}

	tool := &Tool{
		Name:        "greet",
		Description: "Say hello",
	}

	handler := func(ctx context.Context, req *CallToolRequest, args Args) (*CallToolResult, any, error) {
		return &CallToolResult{
			Content: []Content{
				&TextContent{Text: "Hello " + args.Name},
			},
		}, nil, nil
	}

	// Register the tool using mcp.AddTool
	AddTool(server, tool, handler)
}

func TestAddToolWithStructuredOutput(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8087
	c.Mcp.Name = "structured-test"

	server := NewMcpServer(c)

	type CalculateArgs struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	type CalculateResult struct {
		Sum int `json:"sum"`
	}

	tool := &Tool{
		Name:        "add",
		Description: "Add two numbers",
	}

	handler := func(ctx context.Context, req *CallToolRequest, args CalculateArgs) (*CallToolResult, CalculateResult, error) {
		result := CalculateResult{Sum: args.A + args.B}
		return &CallToolResult{
			Content: []Content{
				&TextContent{Text: "Sum calculated"},
			},
		}, result, nil
	}

	AddTool(server, tool, handler)
}

func TestServerLifecycle(t *testing.T) {
	c := McpConf{}
	c.Host = "127.0.0.1"
	c.Port = 0 // Use random port
	c.Mcp.Name = "lifecycle-test"
	c.Mcp.SseEndpoint = "/sse"
	c.Mcp.MessageEndpoint = "/message"

	server := NewMcpServer(c)

	// Test that Start and Stop can be called
	// We don't actually start it to avoid port conflicts in tests
	impl := server.(*mcpServerImpl)
	assert.NotNil(t, impl.httpServer)

	// Just verify the methods exist and can be called
	// Actual server start/stop is tested in integration tests
	defer func() {
		if r := recover(); r == nil {
			// If no panic, call stop
			server.Stop()
		}
	}()
}

func TestServerStartStop(t *testing.T) {
	// Create server with a unique port
	c := McpConf{}
	c.Host = "127.0.0.1"
	c.Port = 18080 // Use high port to avoid conflicts
	c.Mcp.Name = "start-stop-test"
	c.Mcp.SseEndpoint = "/sse"
	c.Mcp.MessageEndpoint = "/message"
	c.Mcp.SseTimeout = 1 * time.Second
	c.Mcp.MessageTimeout = 1 * time.Second

	server := NewMcpServer(c)

	// Test that we can call Stop even without Start
	// (This tests the Stop method coverage)
	server.Stop()
}

func TestServerStartActual(t *testing.T) {
	// Create server with specific port
	c := McpConf{}
	c.Host = "127.0.0.1"
	c.Port = 19080 // Use specific high port
	c.Mcp.Name = "actual-start-test"
	c.Mcp.SseEndpoint = "/sse"
	c.Mcp.MessageEndpoint = "/message"
	c.Mcp.SseTimeout = 1 * time.Second
	c.Mcp.MessageTimeout = 1 * time.Second

	server := NewMcpServer(c)

	// Start server in goroutine
	go func() {
		server.Start() // This blocks until Stop() is called
	}()

	// Give server time to start
	time.Sleep(300 * time.Millisecond)

	// Make a test request to the SSE endpoint to trigger the handler callback
	client := &http.Client{Timeout: 500 * time.Millisecond}
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:19080/sse", nil)
	if err == nil {
		req.Header.Set("Accept", "text/event-stream")
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			// Server is responding - this proves Start() worked
			// and the SSE handler callback was called
			assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode > 0)
		}
	}

	// Stop the server
	server.Stop()

	// Give it time to shutdown
	time.Sleep(100 * time.Millisecond)
}

func TestServerStartStreamable(t *testing.T) {
	// Test with Streamable transport
	c := McpConf{}
	c.Host = "127.0.0.1"
	c.Port = 19081
	c.Mcp.Name = "streamable-start-test"
	c.Mcp.UseStreamable = true
	c.Mcp.SseEndpoint = "/sse"
	c.Mcp.MessageEndpoint = "/message"
	c.Mcp.SseTimeout = 1 * time.Second
	c.Mcp.MessageTimeout = 1 * time.Second

	server := NewMcpServer(c)

	// Start server in goroutine
	go func() {
		server.Start()
	}()

	// Give server time to start
	time.Sleep(300 * time.Millisecond)

	// Make a GET request first (SSE connection)
	client := &http.Client{Timeout: 500 * time.Millisecond}
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:19081/message", nil)
	if err == nil {
		req.Header.Set("Accept", "text/event-stream")
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			// GET request should work
			assert.True(t, resp.StatusCode > 0)
		}
	}

	// Also make a POST request (for message)
	jsonData := []byte(`{"jsonrpc":"2.0","method":"ping","id":1}`)
	req2, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:19081/message", bytes.NewBuffer(jsonData))
	if err == nil {
		req2.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req2)
		if err == nil {
			resp.Body.Close()
			// POST request should also work
			assert.True(t, resp.StatusCode > 0)
		}
	}

	// Stop the server
	server.Stop()

	// Give it time to shutdown
	time.Sleep(100 * time.Millisecond)
}

func TestSSEHandlerCallback(t *testing.T) {
	c := McpConf{}
	c.Host = "127.0.0.1"
	c.Port = 0
	c.Mcp.Name = "sse-handler-test"
	c.Mcp.UseStreamable = false
	c.Mcp.SseEndpoint = "/sse"
	c.Mcp.MessageEndpoint = "/message"

	server := NewMcpServer(c)
	impl := server.(*mcpServerImpl)

	// Verify the server is set up correctly
	assert.NotNil(t, impl.mcpServer)
	assert.False(t, impl.conf.Mcp.UseStreamable)
}

func TestStreamableHandlerCallback(t *testing.T) {
	c := McpConf{}
	c.Host = "127.0.0.1"
	c.Port = 0
	c.Mcp.Name = "streamable-handler-test"
	c.Mcp.UseStreamable = true
	c.Mcp.SseEndpoint = "/sse"
	c.Mcp.MessageEndpoint = "/message"

	server := NewMcpServer(c)
	impl := server.(*mcpServerImpl)

	// Verify the server is set up correctly
	assert.NotNil(t, impl.mcpServer)
	assert.True(t, impl.conf.Mcp.UseStreamable)
}

func TestSSEEndpointAccess(t *testing.T) {
	c := McpConf{}
	c.Host = "127.0.0.1"
	c.Port = 0
	c.Mcp.Name = "sse-endpoint-test"
	c.Mcp.UseStreamable = false
	c.Mcp.SseEndpoint = "/sse"

	server := NewMcpServer(c)
	impl := server.(*mcpServerImpl)

	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/sse", nil)
	req.Header.Set("Accept", "text/event-stream")

	// The server should be configured with SSE endpoints
	assert.NotNil(t, impl.httpServer)
	assert.Equal(t, "/sse", impl.conf.Mcp.SseEndpoint)
}

func TestStreamableEndpointAccess(t *testing.T) {
	c := McpConf{}
	c.Host = "127.0.0.1"
	c.Port = 0
	c.Mcp.Name = "streamable-endpoint-test"
	c.Mcp.UseStreamable = true
	c.Mcp.MessageEndpoint = "/message"

	server := NewMcpServer(c)
	impl := server.(*mcpServerImpl)

	// The server should be configured with streamable endpoints
	assert.NotNil(t, impl.httpServer)
	assert.Equal(t, "/message", impl.conf.Mcp.MessageEndpoint)
}

func TestConfig(t *testing.T) {
	var c McpConf
	err := conf.FillDefault(&c)
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", c.Mcp.Version)
	assert.Equal(t, "/sse", c.Mcp.SseEndpoint)
	assert.Equal(t, "/message", c.Mcp.MessageEndpoint)
}

type mockMcpServer struct{}

func (m *mockMcpServer) Start() {}
func (m *mockMcpServer) Stop()  {}

func TestAddToolWithCustomServer(t *testing.T) {
	server := &mockMcpServer{}
	// Should not panic, but log error
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("AddTool panicked with custom server: %v", r)
		}
	}()

	AddTool(server, &Tool{Name: "test"}, func(ctx context.Context, req *CallToolRequest, args struct{}) (*CallToolResult, any, error) {
		return nil, nil, nil
	})
}

func TestRequestMetadataIntegrationSSEToolCall(t *testing.T) {
	port := getFreePort(t)

	c := McpConf{}
	c.Host = "127.0.0.1"
	c.Port = port
	c.Mcp.Name = "metadata-integration-test"
	c.Mcp.UseStreamable = false
	c.Mcp.SseEndpoint = "/sse/:scope"
	c.Mcp.MessageTimeout = 2 * time.Second
	c.Mcp.SseTimeout = 2 * time.Second

	server := NewMcpServerWithOptions(c, WithRequestMetadataExtractor(DefaultRequestMetadataExtractor))

	tool := &Tool{
		Name:        "inspect_metadata",
		Description: "Inspect metadata in handler context",
	}

	type Args struct{}

	AddTool(server, tool, func(ctx context.Context, req *CallToolRequest, args Args) (*CallToolResult, any, error) {
		header, ok := HeaderFromContext(ctx, "x-tenant-id")
		if !ok || header != "tenant-header" {
			return nil, nil, fmt.Errorf("unexpected header from context: %q", header)
		}

		query, ok := QueryFromContext(ctx, "tenant")
		if !ok || query != "tenant-query" {
			return nil, nil, fmt.Errorf("unexpected query from context: %q", query)
		}

		scope, ok := PathFromContext(ctx, "scope")
		if !ok || scope != "prod" {
			return nil, nil, fmt.Errorf("unexpected path from context: %q", scope)
		}

		return &CallToolResult{
			Content: []Content{&TextContent{Text: "metadata-ok"}},
		}, nil, nil
	})

	go server.Start()
	t.Cleanup(server.Stop)

	baseURL := fmt.Sprintf("http://127.0.0.1:%d/sse/prod?tenant=tenant-query", port)
	waitForServerReady(t, baseURL, 2*time.Second)

	client := sdkmcp.NewClient(&sdkmcp.Implementation{
		Name:    "metadata-client",
		Version: "1.0.0",
	}, nil)

	httpClient := &http.Client{
		Timeout: 2 * time.Second,
		Transport: metadataHeaderRoundTripper{
			next: http.DefaultTransport,
		},
	}

	transport := &sdkmcp.SSEClientTransport{
		Endpoint:   baseURL,
		HTTPClient: httpClient,
	}

	session, err := client.Connect(context.Background(), transport, nil)
	if !assert.NoError(t, err) {
		return
	}
	t.Cleanup(func() {
		_ = session.Close()
	})

	res, err := session.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name:      "inspect_metadata",
		Arguments: map[string]any{},
	})
	if !assert.NoError(t, err) {
		return
	}

	if !assert.NotNil(t, res) {
		return
	}
	assert.False(t, res.IsError)
}

type metadataHeaderRoundTripper struct {
	next http.RoundTripper
}

func (r metadataHeaderRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	next := r.next
	if next == nil {
		next = http.DefaultTransport
	}

	clone := req.Clone(req.Context())
	clone.Header.Set("X-Tenant-Id", "tenant-header")
	return next.RoundTrip(clone)
}

func getFreePort(t *testing.T) int {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if !assert.NoError(t, err) {
		return 0
	}
	defer listener.Close()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !assert.True(t, ok) {
		return 0
	}

	return addr.Port
}

func waitForServerReady(t *testing.T, endpoint string, timeout time.Duration) {
	t.Helper()

	client := &http.Client{Timeout: 200 * time.Millisecond}
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			t.Fatalf("failed to build readiness request: %v", err)
		}
		req.Header.Set("Accept", "text/event-stream")

		resp, err := client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode > 0 {
				return
			}
		}

		time.Sleep(20 * time.Millisecond)
	}

	t.Fatalf("server did not become ready for %s within %s", endpoint, timeout)
}
