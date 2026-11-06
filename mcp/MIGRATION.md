# Migration to Official MCP SDK

This document describes the migration from the custom MCP implementation to the official [go-sdk](https://github.com/modelcontextprotocol/go-sdk).

## Changes

### Dependencies

Added the official MCP SDK:
```bash
go get github.com/modelcontextprotocol/go-sdk@v1.2.0
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
    Server() *sdkmcp.Server  // Returns underlying SDK server
}
```

**Important**: The `AddTool`, `AddPrompt`, and `AddResource` methods have been removed. Use the SDK directly:

```go
// Old (no longer supported)
server.AddTool(tool, handler)

// New (use SDK directly)
sdkmcp.AddTool(server.Server(), tool, handler)
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

// Register with explicit type parameters
sdkmcp.AddTool(server.Server(), tool, handler)
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
   - Endpoint: `/sse` (configurable)
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

// Use SDK directly - no error return
sdkmcp.AddTool(server.Server(), tool, handler)
```

## Benefits

1. **Official SDK**: Uses the official Model Context Protocol SDK
2. **Type Safety**: Go generics provide compile-time type checking
3. **Auto Schema**: JSON schemas generated automatically from struct tags
4. **Dual Transport**: Supports both SSE and Streamable HTTP transports
5. **Maintained**: SDK is actively maintained by the MCP team

## Breaking Changes

1. `server.AddTool()` removed → use `sdkmcp.AddTool(server.Server(), ...)`
2. `server.AddPrompt()` removed (SDK v1.2.0 limitation)
3. `server.AddResource()` removed (SDK v1.2.0 limitation)
4. Config fields `ProtocolVersion` and `BaseUrl` removed
5. All types now come from SDK (re-exported for convenience)

## Migration Checklist

- [ ] Update imports: add `sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"`
- [ ] Replace `server.AddTool()` with `sdkmcp.AddTool(server.Server(), ...)`
- [ ] Remove error handling for tool registration (SDK doesn't return errors)
- [ ] Update config: remove `ProtocolVersion` and `BaseUrl`, add `UseStreamable`
- [ ] Test with both SSE and Streamable transports
- [ ] Update documentation/examples
