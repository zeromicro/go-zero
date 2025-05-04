# Model Context Protocol (MCP) Implementation

## Overview
This package implements the Model Context Protocol (MCP) server specification in Go, providing a framework for real-time communication between AI models and clients using Server-Sent Events (SSE). The implementation follows the standardized protocol for building AI-assisted applications with bidirectional communication capabilities.

## Core Components

### Server-Sent Events (SSE) Communication
- **Real-time Communication**: Robust SSE-based communication system that maintains persistent connections with clients
- **Connection Management**: Client registration, message broadcasting, and client cleanup mechanisms
- **Event Handling**: Event types for tools, prompts, and resources changes

### JSON-RPC Implementation
- **Request Processing**: Complete JSON-RPC request processor for handling MCP protocol methods
- **Response Formatting**: Proper response formatting according to JSON-RPC specifications
- **Error Handling**: Comprehensive error handling with appropriate error codes

### Tool Management
- **Tool Registration**: System to register custom tools with handlers
- **Tool Execution**: Mechanism to execute tool functions with proper timeout handling
- **Result Handling**: Flexible result handling supporting various return types (string, JSON, images)

### Prompt System
- **Prompt Registration**: System for registering both static and dynamic prompts
- **Argument Validation**: Validation for required arguments and default values for optional ones
- **Message Generation**: Handlers that generate properly formatted conversation messages

### Resource Management
- **Resource Registration**: System for managing and accessing external resources
- **Content Delivery**: Handlers for delivering resource content to clients on demand
- **Resource Subscription**: Mechanisms for clients to subscribe to resource updates

### Protocol Features
- **Initialization Sequence**: Proper handshaking with capability negotiation
- **Notification Handling**: Support for both standard and client-specific notifications
- **Message Routing**: Intelligent routing of requests to appropriate handlers

## Technical Highlights

### Configuration System
- **Flexible Configuration**: Configuration system with sensible defaults and customization options
- **CORS Support**: Configurable CORS settings for cross-origin requests
- **Server Information**: Proper server identification and versioning

### Client Session Management
- **Session Tracking**: Client session tracking with unique identifiers
- **Connection Health**: Ping/pong mechanism to maintain connection health
- **Initialization State**: Client initialization state tracking

### Content Handling
- **Multi-format Content**: Support for text, code, and binary content
- **MIME Type Support**: Proper MIME type identification for various content types
- **Audience Annotations**: Content audience annotations for user/assistant targeting

## Usage

### Setting Up an MCP Server

To create and start an MCP server:

```go
package main

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/mcp"
)

func main() {
	// Load configuration from YAML file
	var c mcp.McpConf
	conf.MustLoad("config.yaml", &c)

	// Optional: Disable stats logging
	logx.DisableStat()

	// Create MCP server
	server := mcp.NewMcpServer(c)

	// Register tools, prompts, and resources (examples below)

	// Start the server and ensure it's stopped on exit
	defer server.Stop()
	server.Start()
}
```

Sample configuration file (config.yaml):

```yaml
name: mcp-server
host: localhost
port: 8080
mcp:
  name: my-mcp-server
  messageTimeout: 30s # Timeout for tool calls
  cors:
    - http://localhost:3000 # Optional CORS configuration
```

### Registering Tools

Tools allow AI models to execute custom code through the MCP protocol.

#### Basic Tool Example:

```go
// Register a simple echo tool
echoTool := mcp.Tool{
	Name:        "echo",
	Description: "Echoes back the message provided by the user",
	InputSchema: mcp.InputSchema{
		Properties: map[string]any{
			"message": map[string]any{
				"type":        "string",
				"description": "The message to echo back",
			},
			"prefix": map[string]any{
				"type":        "string",
				"description": "Optional prefix to add to the echoed message",
				"default":     "Echo: ",
			},
		},
		Required: []string{"message"},
	},
	Handler: func(ctx context.Context, params map[string]any) (any, error) {
		var req struct {
			Message string `json:"message"`
			Prefix  string `json:"prefix,optional"`
		}

		if err := mcp.ParseArguments(params, &req); err != nil {
			return nil, fmt.Errorf("failed to parse params: %w", err)
		}

		prefix := "Echo: "
		if len(req.Prefix) > 0 {
			prefix = req.Prefix
		}

		return prefix + req.Message, nil
	},
}

server.RegisterTool(echoTool)
```

#### Tool with Different Response Types:

```go
// Tool returning JSON data
dataTool := mcp.Tool{
	Name:        "data.generate",
	Description: "Generates sample data in various formats",
	InputSchema: mcp.InputSchema{
		Properties: map[string]any{
			"format": map[string]any{
				"type":        "string",
				"description": "Format of data (json, text)",
				"enum":        []string{"json", "text"},
			},
		},
	},
	Handler: func(ctx context.Context, params map[string]any) (any, error) {
		var req struct {
			Format string `json:"format"`
		}

		if err := mcp.ParseArguments(params, &req); err != nil {
			return nil, fmt.Errorf("failed to parse params: %w", err)
		}

		if req.Format == "json" {
			// Return structured data
			return map[string]any{
				"items": []map[string]any{
					{"id": 1, "name": "Item 1"},
					{"id": 2, "name": "Item 2"},
				},
				"count": 2,
			}, nil
		}

		// Default to text
		return "Sample text data", nil
	},
}

server.RegisterTool(dataTool)
```

#### Image Generation Tool Example:

```go
// Tool returning image content
imageTool := mcp.Tool{
	Name:        "image.generate",
	Description: "Generates a simple image",
	InputSchema: mcp.InputSchema{
		Properties: map[string]any{
			"type": map[string]any{
				"type":        "string",
				"description": "Type of image to generate",
				"default":     "placeholder",
			},
		},
	},
	Handler: func(ctx context.Context, params map[string]any) (any, error) {
		// Return image content directly
		return mcp.ImageContent{
			Data:     "base64EncodedImageData...", // Base64 encoded image data
			MimeType: "image/png",
		}, nil
	},
}

server.RegisterTool(imageTool)
```

#### Using ToolResult for Custom Outputs:

```go
// Tool that returns a custom ToolResult type
customResultTool := mcp.Tool{
	Name:        "custom.result",
	Description: "Returns a custom formatted result",
	InputSchema: mcp.InputSchema{
		Properties: map[string]any{
			"resultType": map[string]any{
				"type": "string",
				"enum": []string{"text", "image"},
			},
		},
	},
	Handler: func(ctx context.Context, params map[string]any) (any, error) {
		var req struct {
			ResultType string `json:"resultType"`
		}

		if err := mcp.ParseArguments(params, &req); err != nil {
			return nil, fmt.Errorf("failed to parse params: %w", err)
		}

		if req.ResultType == "image" {
			return mcp.ToolResult{
				Type: mcp.ContentTypeImage,
				Content: map[string]any{
					"data":     "base64EncodedImageData...",
					"mimeType": "image/jpeg",
				},
			}, nil
		}

		// Default to text
		return mcp.ToolResult{
			Type:    mcp.ContentTypeText,
			Content: "This is a text result from ToolResult",
		}, nil
	},
}

server.RegisterTool(customResultTool)
```

### Registering Prompts

Prompts are reusable conversation templates for AI models.

#### Static Prompt Example:

```go
// Register a simple static prompt with placeholders
server.RegisterPrompt(mcp.Prompt{
	Name:        "hello",
	Description: "A simple hello prompt",
	Arguments: []mcp.PromptArgument{
		{
			Name:        "name",
			Description: "The name to greet",
			Required:    false,
		},
	},
	Content: "Say hello to {{name}} and introduce yourself as an AI assistant.",
})
```

#### Dynamic Prompt with Handler Function:

```go
// Register a prompt with a dynamic handler function
server.RegisterPrompt(mcp.Prompt{
	Name:        "dynamic-prompt",
	Description: "A prompt that uses a handler to generate dynamic content",
	Arguments: []mcp.PromptArgument{
		{
			Name:        "username",
			Description: "User's name for personalized greeting",
			Required:    true,
		},
		{
			Name:        "topic",
			Description: "Topic of expertise",
			Required:    true,
		},
	},
	Handler: func(ctx context.Context, args map[string]string) ([]mcp.PromptMessage, error) {
		var req struct {
			Username string `json:"username"`
			Topic    string `json:"topic"`
		}

		if err := mcp.ParseArguments(args, &req); err != nil {
			return nil, fmt.Errorf("failed to parse args: %w", err)
		}

		// Create a user message
		userMessage := mcp.PromptMessage{
			Role: mcp.RoleUser,
			Content: mcp.TextContent{
				Text: fmt.Sprintf("Hello, I'm %s and I'd like to learn about %s.", req.Username, req.Topic),
			},
		}

		// Create an assistant response with current time
		currentTime := time.Now().Format(time.RFC1123)
		assistantMessage := mcp.PromptMessage{
			Role: mcp.RoleAssistant,
			Content: mcp.TextContent{
				Text: fmt.Sprintf("Hello %s! I'm an AI assistant and I'll help you learn about %s. The current time is %s.",
					req.Username, req.Topic, currentTime),
			},
		}

		// Return both messages as a conversation
		return []mcp.PromptMessage{userMessage, assistantMessage}, nil
	},
})
```

#### Multi-Message Prompt with Code Examples:

```go
// Register a prompt that provides code examples in different programming languages
server.RegisterPrompt(mcp.Prompt{
	Name:        "code-example",
	Description: "Provides code examples in different programming languages",
	Arguments: []mcp.PromptArgument{
		{
			Name:        "language",
			Description: "Programming language for the example",
			Required:    true,
		},
		{
			Name:        "complexity",
			Description: "Complexity level (simple, medium, advanced)",
		},
	},
	Handler: func(ctx context.Context, args map[string]string) ([]mcp.PromptMessage, error) {
		var req struct {
			Language   string `json:"language"`
			Complexity string `json:"complexity,optional"`
		}

		if err := mcp.ParseArguments(args, &req); err != nil {
			return nil, fmt.Errorf("failed to parse args: %w", err)
		}

		// Validate language
		supportedLanguages := map[string]bool{"go": true, "python": true, "javascript": true, "rust": true}
		if !supportedLanguages[req.Language] {
			return nil, fmt.Errorf("unsupported language: %s", req.Language)
		}

		// Generate code example based on language and complexity
		var codeExample string

		switch req.Language {
		case "go":
			if req.Complexity == "simple" {
				codeExample = `
package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`
			} else {
				codeExample = `
package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	fmt.Printf("Hello, World! Current time is %s\n", now.Format(time.RFC3339))
}`
			}
		case "python":
			// Python example code
			if req.Complexity == "simple" {
				codeExample = `
def greet(name):
    return f"Hello, {name}!"

print(greet("World"))`
			} else {
				codeExample = `
import datetime

def greet(name, include_time=False):
    message = f"Hello, {name}!"
    if include_time:
        message += f" Current time is {datetime.datetime.now().isoformat()}"
    return message

print(greet("World", include_time=True))`
			}
		}

		// Create messages array according to MCP spec
		messages := []mcp.PromptMessage{
			{
				Role: mcp.RoleAssistant,
				Content: mcp.TextContent{
					Text: fmt.Sprintf("You are a helpful coding assistant specialized in %s programming.", req.Language),
				},
			},
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Text: fmt.Sprintf("Show me a %s example of a Hello World program in %s.", req.Complexity, req.Language),
				},
			},
			{
				Role: mcp.RoleAssistant,
				Content: mcp.TextContent{
					Text: fmt.Sprintf("Here's a %s example in %s:\n\n```%s%s\n```\n\nHow can I help you implement this?",
						req.Complexity, req.Language, req.Language, codeExample),
				},
			},
		}

		return messages, nil
	},
})
```

### Registering Resources

Resources provide access to external content such as files or generated data.

#### Basic Resource Example:

```go
// Register a static resource
server.RegisterResource(mcp.Resource{
	Name:        "example-document",
	URI:         "file:///example/document.txt",
	Description: "An example document",
	MimeType:    "text/plain",
	Handler: func(ctx context.Context) (mcp.ResourceContent, error) {
		return mcp.ResourceContent{
			URI:      "file:///example/document.txt",
			MimeType: "text/plain",
			Text:     "This is an example document content.",
		}, nil
	},
})
```

#### Dynamic Resource with Code Example:

```go
// Register a Go code resource with dynamic handler
server.RegisterResource(mcp.Resource{
	Name:        "go-example",
	URI:         "file:///project/src/main.go",
	Description: "A simple Go example with multiple files",
	MimeType:    "text/x-go",
	Handler: func(ctx context.Context) (mcp.ResourceContent, error) {
		// Return ResourceContent with all required fields
		return mcp.ResourceContent{
			URI:      "file:///project/src/main.go",
			MimeType: "text/x-go",
			Text:     "package main\n\nimport (\n\t\"fmt\"\n\t\"./greeting\"\n)\n\nfunc main() {\n\tfmt.Println(greeting.Hello(\"world\"))\n}",
		}, nil
	},
})

// Register a companion file for the above example
server.RegisterResource(mcp.Resource{
	Name:        "go-greeting",
	URI:         "file:///project/src/greeting/greeting.go",
	Description: "A greeting package for the Go example",
	MimeType:    "text/x-go",
	Handler: func(ctx context.Context) (mcp.ResourceContent, error) {
		return mcp.ResourceContent{
			URI:      "file:///project/src/greeting/greeting.go",
			MimeType: "text/x-go",
			Text:     "package greeting\n\nfunc Hello(name string) string {\n\treturn \"Hello, \" + name + \"!\"\n}",
		}, nil
	},
})
```

#### Binary Resource Example:

```go
// Register a binary resource (like an image)
server.RegisterResource(mcp.Resource{
	Name:        "example-image",
	URI:         "file:///example/image.png",
	Description: "An example image",
	MimeType:    "image/png",
	Handler: func(ctx context.Context) (mcp.ResourceContent, error) {
		// Read image from file or generate it
		imageData := "base64EncodedImageData..." // Base64 encoded image data

		return mcp.ResourceContent{
			URI:      "file:///example/image.png",
			MimeType: "image/png",
			Blob:     imageData, // For binary data
		}, nil
	},
})
```

### Using Resources in Prompts

You can embed resources in prompt responses to create rich interactions with proper MCP-compliant structure:

```go
// Register a prompt that embeds a resource
server.RegisterPrompt(mcp.Prompt{
	Name:        "resource-example",
	Description: "A prompt that embeds a resource",
	Arguments: []mcp.PromptArgument{
		{
			Name:        "file_type",
			Description: "Type of file to show (rust or go)",
			Required:    true,
		},
	},
	Handler: func(ctx context.Context, args map[string]string) ([]mcp.PromptMessage, error) {
		var req struct {
			FileType string `json:"file_type"`
		}

		if err := mcp.ParseArguments(args, &req); err != nil {
			return nil, fmt.Errorf("failed to parse args: %w", err)
		}

		var resourceURI, mimeType, fileContent string
		if req.FileType == "rust" {
			resourceURI = "file:///project/src/main.rs"
			mimeType = "text/x-rust"
			fileContent = "fn main() {\n    println!(\"Hello world!\");\n}"
		} else {
			resourceURI = "file:///project/src/main.go"
			mimeType = "text/x-go"
			fileContent = "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello, world!\")\n}"
		}

		// Create message with embedded resource using proper MCP format
		return []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Text: fmt.Sprintf("Can you explain this %s code?", req.FileType),
				},
			},
			{
				Role: mcp.RoleAssistant,
				Content: mcp.EmbeddedResource{
					Type: mcp.ContentTypeResource,
					Resource: struct {
						URI      string `json:"uri"`
						MimeType string `json:"mimeType"`
						Text     string `json:"text,omitempty"`
						Blob     string `json:"blob,omitempty"`
					}{
						URI:      resourceURI,
						MimeType: mimeType,
						Text:     fileContent,
					},
				},
			},
			{
				Role: mcp.RoleAssistant,
				Content: mcp.TextContent{
					Text: fmt.Sprintf("Above is a simple Hello World example in %s. Let me explain how it works.", req.FileType),
				},
			},
		}, nil
	},
})
```

### Multiple File Resources Example

```go
// Register a prompt that demonstrates embedding multiple resource files
server.RegisterPrompt(mcp.Prompt{
	Name:        "go-code-example",
	Description: "A prompt that correctly embeds multiple resource files",
	Arguments: []mcp.PromptArgument{
		{
			Name:        "format",
			Description: "How to format the code display",
		},
	},
	Handler: func(ctx context.Context, args map[string]string) ([]mcp.PromptMessage, error) {
		var req struct {
			Format string `json:"format,optional"`
		}

		if err := mcp.ParseArguments(args, &req); err != nil {
			return nil, fmt.Errorf("failed to parse args: %w", err)
		}

		// Get the Go code for multiple files
		var mainGoText string = "package main\n\nimport (\n\t\"fmt\"\n\t\"./greeting\"\n)\n\nfunc main() {\n\tfmt.Println(greeting.Hello(\"world\"))\n}"
		var greetingGoText string = "package greeting\n\nfunc Hello(name string) string {\n\treturn \"Hello, \" + name + \"!\"\n}"

		// Create message with properly formatted embedded resource per MCP spec
		messages := []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Text: "Show me a simple Go example with proper imports.",
				},
			},
			{
				Role: mcp.RoleAssistant,
				Content: mcp.TextContent{
					Text: "Here's a simple Go example project:",
				},
			},
			{
				Role: mcp.RoleAssistant,
				Content: mcp.EmbeddedResource{
					Type: mcp.ContentTypeResource,
					Resource: struct {
						URI      string `json:"uri"`
						MimeType string `json:"mimeType"`
						Text     string `json:"text,omitempty"`
						Blob     string `json:"blob,omitempty"`
					}{
						URI:      "file:///project/src/main.go",
						MimeType: "text/x-go",
						Text:     mainGoText,
					},
				},
			},
		}

		// Add explanation and additional file if requested
		if req.Format == "with_explanation" {
			messages = append(messages, mcp.PromptMessage{
				Role: mcp.RoleAssistant,
				Content: mcp.TextContent{
					Text: "This example demonstrates a simple Go application with modular structure. The main.go file imports from a local 'greeting' package that provides the Hello function.",
				},
			})

			// Also show the greeting.go file with correct resource format
			messages = append(messages, mcp.PromptMessage{
				Role: mcp.RoleAssistant,
				Content: mcp.EmbeddedResource{
					Type: mcp.ContentTypeResource,
					Resource: struct {
						URI      string `json:"uri"`
						MimeType string `json:"mimeType"`
						Text     string `json:"text,omitempty"`
						Blob     string `json:"blob,omitempty"`
					}{
						URI:      "file:///project/src/greeting/greeting.go",
						MimeType: "text/x-go",
						Text:     greetingGoText,
					},
				},
			})
		}

		return messages, nil
	},
})
```

### Complete Application Example

Here's a complete example demonstrating all the components:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/mcp"
)

func main() {
	// Load configuration
	var c mcp.McpConf
	if err := conf.Load("config.yaml", &c); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set up logging
	logx.DisableStat()

	// Create MCP server
	server := mcp.NewMcpServer(c)
	defer server.Stop()

	// Register a simple echo tool
	echoTool := mcp.Tool{
		Name:        "echo",
		Description: "Echoes back the message provided by the user",
		InputSchema: mcp.InputSchema{
			Properties: map[string]any{
				"message": map[string]any{
					"type":        "string",
					"description": "The message to echo back",
				},
				"prefix": map[string]any{
					"type":        "string",
					"description": "Optional prefix to add to the echoed message",
					"default":     "Echo: ",
				},
			},
			Required: []string{"message"},
		},
		Handler: func(ctx context.Context, params map[string]any) (any, error) {
			var req struct {
				Message string `json:"message"`
				Prefix  string `json:"prefix,optional"`
			}

			if err := mcp.ParseArguments(params, &req); err != nil {
				return nil, fmt.Errorf("failed to parse args: %w", err)
			}

			prefix := "Echo: "
			if len(req.Prefix) > 0 {
				prefix = req.Prefix
			}

			return prefix + req.Message, nil
		},
	}
	server.RegisterTool(echoTool)

	// Register a static prompt
	server.RegisterPrompt(mcp.Prompt{
		Name:        "greeting",
		Description: "A simple greeting prompt",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "name",
				Description: "The name to greet",
				Required:    true,
			},
		},
		Content: "Hello {{name}}! How can I assist you today?",
	})

	// Register a dynamic prompt
	server.RegisterPrompt(mcp.Prompt{
		Name:        "dynamic-prompt",
		Description: "A prompt that uses a handler to generate dynamic content",
		Arguments: []mcp.PromptArgument{
			{
				Name:        "username",
				Description: "User's name for personalized greeting",
				Required:    true,
			},
			{
				Name:        "topic",
				Description: "Topic of expertise",
				Required:    true,
			},
		},
		Handler: func(ctx context.Context, args map[string]string) ([]mcp.PromptMessage, error) {
			var req struct {
				Username string `json:"username"`
				Topic    string `json:"topic"`
			}

			if err := mcp.ParseArguments(args, &req); err != nil {
				return nil, fmt.Errorf("failed to parse args: %w", err)
			}

			// Create messages with current time
			currentTime := time.Now().Format(time.RFC1123)
			return []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Text: fmt.Sprintf("Hello, I'm %s and I'd like to learn about %s.", req.Username, req.Topic),
					},
				},
				{
					Role: mcp.RoleAssistant,
					Content: mcp.TextContent{
						Text: fmt.Sprintf("Hello %s! I'm an AI assistant and I'll help you learn about %s. The current time is %s.",
							req.Username, req.Topic, currentTime),
					},
				},
			}, nil
		},
	})

	// Register a resource
	server.RegisterResource(mcp.Resource{
		Name:        "example-doc",
		URI:         "file:///example/doc.txt",
		Description: "An example document",
		MimeType:    "text/plain",
		Handler: func(ctx context.Context) (mcp.ResourceContent, error) {
			return mcp.ResourceContent{
				URI:      "file:///example/doc.txt",
				MimeType: "text/plain",
				Text:     "This is the content of the example document.",
			}, nil
		},
	})

	// Start the server
	fmt.Printf("Starting MCP server on %s:%d\n", c.Host, c.Port)
	server.Start()
}
```

## Error Handling

The MCP implementation provides comprehensive error handling:

- Tool execution errors are properly reported back to clients
- Missing or invalid parameters are detected and reported with appropriate error codes
- Resource and prompt lookup failures are handled gracefully
- Timeout handling for long-running tool executions using context
- Panic recovery to prevent server crashes

## Advanced Features

- **Annotations**: Add audience and priority metadata to content
- **Content Types**: Support for text, images, audio, and other content formats
- **Embedded Resources**: Include file resources directly in prompt responses
- **Context Awareness**: All handlers receive context.Context for timeout and cancellation support
- **Progress Tokens**: Support for tracking progress of long-running operations
- **Customizable Timeouts**: Configure execution timeouts for tools and operations

## Performance Considerations

- Tool execution runs with configurable timeouts to prevent blocking
- Efficient client tracking and cleanup to prevent resource leaks
- Proper concurrency handling with mutex protection for shared resources
- Buffered message channels to prevent blocking on client message delivery
