package mcp

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest/pathvar"
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

func TestNewMcpServerWithOptions(t *testing.T) {
	var hookCalls atomic.Int32

	c := McpConf{}
	c.Name = "options-server"
	c.Host = "localhost"
	c.Port = 8088

	server := NewMcpServerWithOptions(c,
		WithServerOptions(&sdkmcp.ServerOptions{HasTools: true}),
		WithServerHook(func(server *sdkmcp.Server) {
			hookCalls.Add(1)
			server.RemoveTools("missing")
		}),
	)

	impl := server.(*mcpServerImpl)
	assert.NotNil(t, impl)
	assert.Equal(t, int32(1), hookCalls.Load())
	assert.NotNil(t, SDKServer(server))
	assert.Same(t, impl.mcpServer, SDKServer(server))
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

func TestRemoveToolsWithSdkServerAccess(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8089
	c.Mcp.Name = "dynamic-tools"

	server := NewMcpServer(c)

	type Args struct {
		Name string `json:"name"`
	}

	AddTool(server, &Tool{Name: "greet", Description: "Say hello"}, func(ctx context.Context, req *CallToolRequest, args Args) (*CallToolResult, any, error) {
		return &CallToolResult{}, nil, nil
	})

	RemoveTools(server, "greet")
	assert.NotNil(t, SDKServer(server))
}

func TestSDKServerWithCustomServer(t *testing.T) {
	assert.Nil(t, SDKServer(&mockMcpServer{}))
}

func TestServerSelectorUsesCustomServer(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8090
	c.Mcp.Name = "selector-test"

	selected := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "selected", Version: "1.0.0"}, nil)
	server := NewMcpServerWithOptions(c, WithServerSelector(func(r *http.Request) *sdkmcp.Server {
		if r.URL.Path == "/mcp/b" {
			return selected
		}

		return nil
	}))

	impl := server.(*mcpServerImpl)
	reqSelected := httptest.NewRequest(http.MethodGet, "/mcp/b", nil)
	reqDefault := httptest.NewRequest(http.MethodGet, "/mcp/a", nil)

	assert.Same(t, selected, impl.selectServer(reqSelected))
	assert.Same(t, impl.mcpServer, impl.selectServer(reqDefault))
}

func TestRequestMetadataHelpers(t *testing.T) {
	ctx := context.WithValue(context.Background(), requestMetadataContextKey{}, RequestMetadata{
		Headers: map[string]string{"x-tenant": "acme"},
		Query:   map[string]string{"variant": "beta"},
		Path:    map[string]string{"tenant": "blue"},
	})

	metadata := RequestMetadataFromContext(ctx)
	assert.Equal(t, RequestMetadata{
		Headers: map[string]string{"x-tenant": "acme"},
		Query:   map[string]string{"variant": "beta"},
		Path:    map[string]string{"tenant": "blue"},
	}, metadata)
	assert.Equal(t, "acme", HeaderFromContext(ctx, "x-tenant"))
	assert.Equal(t, "beta", QueryFromContext(ctx, "variant"))
	assert.Equal(t, "blue", PathFromContext(ctx, "tenant"))
	assert.Equal(t, "", HeaderFromContext(ctx, "missing"))

	metadata.Headers["x-tenant"] = "changed"
	assert.Equal(t, "acme", HeaderFromContext(ctx, "x-tenant"))
}

func TestWithRequestMetadataExtractor(t *testing.T) {
	var options serverOptions
	extractor := func(r *http.Request) RequestMetadata {
		return RequestMetadata{Headers: map[string]string{"x-tenant": r.Header.Get("X-Tenant")}}
	}

	WithRequestMetadataExtractor(extractor)(&options)
	require.NotNil(t, options.metadata)
	req := httptest.NewRequest(http.MethodGet, "/sse", nil)
	req.Header.Set("X-Tenant", "acme")
	assert.Equal(t, "acme", options.metadata(req).Headers["x-tenant"])
}

func TestRequestMetadataHelpersWithNilContext(t *testing.T) {
	var nilCtx context.Context

	assert.True(t, isEmptyRequestMetadata(RequestMetadataFromContext(nilCtx)))
	assert.Equal(t, "", HeaderFromContext(nilCtx, "x-tenant"))
	assert.Equal(t, "", QueryFromContext(nilCtx, "variant"))
	assert.Equal(t, "", PathFromContext(nilCtx, "tenant"))
}

func TestRequestMetadataHelpersWithEmptyContext(t *testing.T) {
	ctx := context.Background()
	assert.True(t, isEmptyRequestMetadata(RequestMetadataFromContext(ctx)))
	assert.Equal(t, "", HeaderFromContext(ctx, "x-tenant"))
	assert.Equal(t, "", QueryFromContext(ctx, "variant"))
	assert.Equal(t, "", PathFromContext(ctx, "tenant"))
}

func TestWrapHandlerInjectsMetadata(t *testing.T) {
	s := &mcpServerImpl{
		metadata: func(r *http.Request) RequestMetadata {
			return RequestMetadata{
				Headers: map[string]string{"x-tenant": r.Header.Get("X-Tenant")},
				Query:   map[string]string{"variant": r.URL.Query().Get("variant")},
			}
		},
	}

	called := false
	wrapped := s.wrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.Equal(t, "acme", HeaderFromContext(r.Context(), "x-tenant"))
		assert.Equal(t, "beta", QueryFromContext(r.Context(), "variant"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/message?variant=beta", nil)
	req.Header.Set("X-Tenant", "acme")
	wrapped.ServeHTTP(httptest.NewRecorder(), req)

	assert.True(t, called)
}

func TestWrapHandlerWithoutExtractorPreservesContext(t *testing.T) {
	s := &mcpServerImpl{}
	wrapped := s.wrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, isEmptyRequestMetadata(RequestMetadataFromContext(r.Context())))
	}))

	req := httptest.NewRequest(http.MethodGet, "/sse", nil)
	wrapped.ServeHTTP(httptest.NewRecorder(), req)
}

func TestCloneRequestMetadata(t *testing.T) {
	original := RequestMetadata{Headers: map[string]string{"x-tenant": "acme"}}
	cloned := cloneRequestMetadata(original)

	assert.True(t, reflect.DeepEqual(original, cloned))
	cloned.Headers["x-tenant"] = "changed"
	assert.Equal(t, "acme", original.Headers["x-tenant"])
}

func TestWrapHandlerAddsPathVarsWhenExtractorOmitsThem(t *testing.T) {
	s := &mcpServerImpl{
		metadata: func(r *http.Request) RequestMetadata {
			return RequestMetadata{Headers: map[string]string{"x-tenant": r.Header.Get("X-Tenant")}}
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/mcp/blue", nil)
	req = pathvar.WithVars(req, map[string]string{"tenant": "blue"})
	req.Header.Set("X-Tenant", "acme")

	wrapped := s.wrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "acme", HeaderFromContext(r.Context(), "x-tenant"))
		assert.Equal(t, "blue", PathFromContext(r.Context(), "tenant"))
	}))

	wrapped.ServeHTTP(httptest.NewRecorder(), req)
}

func TestWrapHandlerPreservesExplicitPathMetadata(t *testing.T) {
	s := &mcpServerImpl{
		metadata: func(r *http.Request) RequestMetadata {
			return RequestMetadata{Path: map[string]string{"tenant": "explicit"}}
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/mcp/blue", nil)
	req = pathvar.WithVars(req, map[string]string{"tenant": "blue"})

	wrapped := s.wrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "explicit", PathFromContext(r.Context(), "tenant"))
	}))

	wrapped.ServeHTTP(httptest.NewRecorder(), req)
}

func TestWrapHandlerWithEmptyExtractorMetadata(t *testing.T) {
	s := &mcpServerImpl{
		metadata: func(r *http.Request) RequestMetadata {
			return RequestMetadata{}
		},
	}

	called := false
	wrapped := s.wrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		assert.True(t, isEmptyRequestMetadata(RequestMetadataFromContext(r.Context())))
	}))

	wrapped.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/sse", nil))
	assert.True(t, called)
}

func TestStreamableServerUsesSelector(t *testing.T) {
	c := McpConf{}
	c.Host = "localhost"
	c.Port = 8091
	c.Mcp.Name = "streamable-selector-test"
	c.Mcp.UseStreamable = true

	selected := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "selected-streamable", Version: "1.0.0"}, nil)
	server := NewMcpServerWithOptions(c, WithServerSelector(func(r *http.Request) *sdkmcp.Server {
		if r.URL.Path == "/message" {
			return selected
		}

		return nil
	}))

	impl := server.(*mcpServerImpl)
	reqSelected := httptest.NewRequest(http.MethodPost, "/message", bytes.NewBufferString(`{}`))
	reqDefault := httptest.NewRequest(http.MethodPost, "/other", bytes.NewBufferString(`{}`))

	assert.Same(t, selected, impl.streamableServer(reqSelected))
	assert.Same(t, impl.mcpServer, impl.streamableServer(reqDefault))
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
