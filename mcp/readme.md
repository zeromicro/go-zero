# Model Context Protocol (MCP) Implementation

## Overview

This package provides a go-zero integration for the [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) using the official [go-sdk](https://github.com/modelcontextprotocol/go-sdk). It wraps the official MCP SDK to provide a seamless integration with go-zero's REST server framework.

## Features

- **Official SDK Integration**: Built on top of the official Model Context Protocol Go SDK
- **No SDK Import Required**: Use `mcp.AddTool()` directly without importing the official SDK
- **go-zero Integration**: Seamlessly integrates with go-zero's REST server and configuration system
- **Dual Transport Support**:
  - SSE (Server-Sent Events) transport for 2024-11-05 MCP spec
  - Streamable HTTP transport for 2025-03-26 MCP spec
- **CORS Support**: Configurable CORS settings for cross-origin requests
- **Type-Safe Tool Handlers**: Generic tool handlers with automatic JSON schema generation
- **Prompts and Resources**: Full support for MCP prompts and resources

## Quick Start

### 1. Installation

```bash
go get github.com/zeromicro/go-zero
```

**Note**: The official MCP SDK is a transitive dependency and will be installed automatically. You don't need to import it directly in your code.

### 2. Configuration

Create a configuration file `config.yaml`:

```yaml
name: my-mcp-server
host: localhost
port: 8080
mcp:
  name: my-mcp-server
  version: 1.0.0
  useStreamable: false  # Use SSE transport (default), set to true for Streamable HTTP
  sseEndpoint: /sse
  messageEndpoint: /message
  sseTimeout: 24h
  messageTimeout: 30s
  cors:
    - http://localhost:3000
```

### 3. Create Your Server

```go
package main

import (
	"context"
	"log"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/mcp"
)

type GreetArgs struct {
	Name string `json:"name" jsonschema:"description=Name of the person to greet"`
}

func main() {
	// Load configuration
	var c mcp.McpConf
	conf.MustLoad("config.yaml", &c)

	// Create MCP server
	server := mcp.NewMcpServer(c)

	// Register a tool with automatic schema generation using the SDK directly
	tool := &mcp.Tool{
		Name:        "greet",
		Description: "Greet someone by name",
	}

	handler := func(ctx context.Context, req *mcp.CallToolRequest, args GreetArgs) (*mcp.CallToolResult, any, error) {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Hello, " + args.Name + "!"},
			},
		}, nil, nil
	}

	// Register tool with type-safe generics - no need to import official SDK
	mcp.AddTool(server, tool, handler)

	// Start server
	defer server.Stop()
	server.Start()
}
```

## Adding Tools

Tools are functions that the MCP client can call. The SDK automatically generates JSON schemas from your struct tags. Use `sdkmcp.AddTool` with the server's underlying SDK server:

```go
type CalculateArgs struct {
	Operation string  `json:"operation" jsonschema:"enum=add,enum=subtract,enum=multiply,enum=divide"`
	A         float64 `json:"a" jsonschema:"description=First number"`
	B         float64 `json:"b" jsonschema:"description=Second number"`
}

tool := &mcp.Tool{
	Name:        "calculate",
	Description: "Perform arithmetic operations",
}

handler := func(ctx context.Context, req *mcp.CallToolRequest, args CalculateArgs) (*mcp.CallToolResult, any, error) {
	var result float64
	switch args.Operation {
	case "add":
		result = args.A + args.B
	case "subtract":
		result = args.A - args.B
	case "multiply":
		result = args.A * args.B
	case "divide":
		if args.B == 0 {
			return &mcp.CallToolResult{IsError: true}, nil, fmt.Errorf("division by zero")
		}
		result = args.A / args.B
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Result: %v", result)},
		},
	}, result, nil
}

// Register tool
mcp.AddTool(server, tool, handler)
```

## Adding Prompts

Prompts provide reusable message templates:

```go
prompt := &mcp.Prompt{
	Name:        "code-review",
	Description: "Review code for best practices",
}

handler := func(ctx context.Context, req *sdkmcp.GetPromptRequest, args map[string]string) (*mcp.GetPromptResult, error) {
	code := args["code"]
	language := args["language"]

	return &mcp.GetPromptResult{
		Messages: []mcp.PromptMessage{
			{
				Role: "user",
				Content: &mcp.TextContent{
					Text: fmt.Sprintf("Please review this %s code:\n\n%s", language, code),
				},
			},
		},
	}, nil
}

server.AddPrompt(prompt, handler)
```

## Adding Resources

Resources provide access to data that the model can read:

```go
resource := &mcp.Resource{
	URI:         "file:///docs/readme.md",
	Name:        "README",
	Description: "Project documentation",
	MimeType:    "text/markdown",
}

handler := func(ctx context.Context, req *sdkmcp.ReadResourceRequest, uri string) (*mcp.ReadResourceResult, error) {
	content, err := os.ReadFile("README.md")
	if err != nil {
		return nil, err
	}

	return &mcp.ReadResourceResult{
		Contents: []mcp.ResourceContents{
			{
				URI:      uri,
				MimeType: "text/markdown",
				Text:     string(content),
			},
		},
	}, nil
}

server.AddResource(resource, handler)
```

## Transport Options

### SSE Transport (Default)

The SSE (Server-Sent Events) transport is the original MCP transport from the 2024-11-05 specification:

```yaml
mcp:
  useStreamable: false
  sseEndpoint: /sse
```

### Streamable HTTP Transport

The newer Streamable HTTP transport from the 2025-03-26 specification provides better connection management:

```yaml
mcp:
  useStreamable: true
  messageEndpoint: /message
```

## Configuration Options

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `name` | string | | Server name (required, from RestConf) |
| `host` | string | | Server host (required, from RestConf) |
| `port` | int | | Server port (required, from RestConf) |
| `mcp.name` | string | | MCP server name (defaults to `name`) |
| `mcp.version` | string | `1.0.0` | Server version |
| `mcp.useStreamable` | bool | `false` | Use Streamable HTTP transport instead of SSE |
| `mcp.sseEndpoint` | string | `/sse` | SSE endpoint path |
| `mcp.messageEndpoint` | string | `/message` | Message endpoint path |
| `mcp.sseTimeout` | duration | `24h` | SSE connection timeout |
| `mcp.messageTimeout` | duration | `30s` | Message processing timeout |
| `mcp.cors` | []string | | Allowed CORS origins |

## Examples

See the `adhoc/mcp` directory for a complete working example.

## Official SDK Documentation

For more details on the underlying MCP SDK, see:
- [Official Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [MCP Specification](https://modelcontextprotocol.io/)
- [SDK Documentation](https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk/mcp)

## License

This implementation follows the go-zero project license (MIT).
