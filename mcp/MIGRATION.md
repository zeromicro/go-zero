# Migration to Official MCP SDK

This document describes the migration from the custom MCP implementation to the official [go-sdk](https://github.com/modelcontextprotocol/go-sdk).

## Changes

### Dependencies

Added the official MCP SDK:
```bash
go get github.com/modelcontextprotocol/go-sdk@v1.4.0
```

### Type System

All types are now re-exported from the official SDK:
- `Tool` → `sdkmcp.Tool`
- `CallToolRequest` → `sdkmcp.CallToolRequest`
- `CallToolResult` → `sdkmcp.CallToolResult`
- Content types (`TextContent`, `ImageContent`, etc.)
- `Prompt`, `Resource`, `Server`, `ServerSession`

### Server Interface

The `McpServer` interface has been simplified:

```go
type McpServer interface {
    Start()
    Stop()
}
```

**Important**: The `AddTool`, `AddPrompt`, and `AddResource` methods have been removed from the server interface. Use the helper functions or the underlying SDK server:

```go
// Old (no longer supported)
server.AddTool(tool, handler)

// New (go-zero helper)
mcp.AddTool(server, tool, handler)

// Or use the underlying SDK server directly
sdkServer := mcp.SDKServer(server)
sdkmcp.AddTool(sdkServer, tool, handler)
```

For advanced use cases, use the additive constructor with options:

```go
server := mcp.NewMcpServerWithOptions(c,
    mcp.WithServerOptions(&mcp.ServerOptions{HasTools: true}),
    mcp.WithServerHook(func(s *mcp.Server) {
        // middleware, prompts, resources, dynamic tool changes, etc.
    }),
    mcp.WithRequestMetadataExtractor(func(r *http.Request) mcp.RequestMetadata {
        return mcp.RequestMetadata{
            Headers: map[string]string{"x-tenant": r.Header.Get("X-Tenant")},
        }
    }),
    mcp.WithServerSelector(func(r *http.Request) *mcp.Server {
        return nil // fallback to the default SDK server
    }),
)
```

Handlers can then read selected metadata from `context.Context` without changing existing handler signatures:

```go
tenant := mcp.HeaderFromContext(ctx, "x-tenant")
variant := mcp.QueryFromContext(ctx, "variant")
pathID := mcp.PathFromContext(ctx, "id")
```

### Configuration

Updated configuration structure:
- Removed: `ProtocolVersion`, `BaseUrl` (SDK manages these)
- Added: `UseStreamable` (choose between SSE and Streamable HTTP transport)

```yaml
mcp:
  name: my-server
  version: 1.0.0
  useStreamable: false  # false = SSE (2024-11-05), true = Streamable HTTP (2025-03-26)
  sseEndpoint: /sse
  messageEndpoint: /message
  sseTimeout: 24h
  messageTimeout: 30s
  cors:
    - http://localhost:3000
```

### Tool Registration

The SDK uses Go generics for type-safe tool registration:

```go
import sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

type MyArgs struct {
    Value string `json:"value" jsonschema:"description=Input value"`
}

tool := &mcp.Tool{
    Name:        "my_tool",
    Description: "Description",
}

handler := func(ctx context.Context, req *mcp.CallToolRequest, args MyArgs) (*mcp.CallToolResult, any, error) {
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: "Result"},
        },
    }, nil, nil
}

// Register with the go-zero helper
mcp.AddTool(server, tool, handler)

// Or use the underlying SDK server directly
sdkmcp.AddTool(mcp.SDKServer(server), tool, handler)
```

The SDK automatically generates JSON schemas from struct tags.

### Transport Support

Two transports are supported:

1. **SSE (Server-Sent Events)**: 2024-11-05 MCP spec
   - Default (`UseStreamable: false`)
   - Endpoint: `/sse` (configurable)
   - Bidirectional: client sends messages to `/message`

2. **Streamable HTTP**: 2025-03-26 MCP spec
   - Opt-in (`UseStreamable: true`)
   - Endpoint: `/message` (configurable)
   - Newer protocol with improved streaming

### Example Migration

**Before:**
```go
server := mcp.NewMcpServer(c)

tool := &mcp.Tool{Name: "greet", Description: "Greet"}
handler := func(ctx context.Context, req *mcp.CallToolRequest, args GreetArgs) (*mcp.CallToolResult, any, error) {
    return &mcp.CallToolResult{
        Content: []mcp.Content{&mcp.TextContent{Text: "Hello"}},
    }, nil, nil
}

if err := server.AddTool(tool, handler); err != nil {
    log.Fatal(err)
}
```

**After:**
```go
import sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

server := mcp.NewMcpServer(c)

tool := &mcp.Tool{Name: "greet", Description: "Greet"}
handler := func(ctx context.Context, req *mcp.CallToolRequest, args GreetArgs) (*mcp.CallToolResult, any, error) {
    return &mcp.CallToolResult{
        Content: []mcp.Content{&mcp.TextContent{Text: "Hello"}},
    }, nil, nil
}

// Use the go-zero helper
mcp.AddTool(server, tool, handler)

// Or use the underlying SDK server directly
sdkmcp.AddTool(mcp.SDKServer(server), tool, handler)
```

## Benefits

1. **Official SDK**: Uses the official Model Context Protocol SDK
2. **Type Safety**: Go generics provide compile-time type checking
3. **Auto Schema**: JSON schemas generated automatically from struct tags
4. **Dual Transport**: Supports both SSE and Streamable HTTP transports
5. **Maintained**: SDK is actively maintained by the MCP team

## Breaking Changes

1. `server.AddTool()` removed → use `mcp.AddTool(server, ...)` or `sdkmcp.AddTool(mcp.SDKServer(server), ...)`
2. `server.AddPrompt()` removed → use `mcp.SDKServer(server).AddPrompt(...)`
3. `server.AddResource()` removed → use `mcp.SDKServer(server).AddResource(...)`
4. Config fields `ProtocolVersion` and `BaseUrl` removed
5. All types now come from SDK (re-exported for convenience)
6. Advanced runtime customization is now opt-in via `NewMcpServerWithOptions(...)`

## Migration Checklist

- [ ] Update imports: add `sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"`
- [ ] Replace `server.AddTool()` with `mcp.AddTool(server, ...)` or `sdkmcp.AddTool(mcp.SDKServer(server), ...)`
- [ ] Remove error handling for tool registration (SDK doesn't return errors)
- [ ] Update config: remove `ProtocolVersion` and `BaseUrl`, add `UseStreamable`
- [ ] Test with both SSE and Streamable transports
- [ ] Update documentation/examples
