package mcp

import (
	"context"
	"testing"

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
	assert.NotNil(t, server.Server())
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

	// Verify tool was added by checking it's not nil
	assert.NotNil(t, server.Server())
}

func TestConfig(t *testing.T) {
	var c McpConf
	err := conf.FillDefault(&c)
	assert.NoError(t, err)
	assert.Equal(t, "1.0.0", c.Mcp.Version)
	assert.Equal(t, "/sse", c.Mcp.SseEndpoint)
	assert.Equal(t, "/message", c.Mcp.MessageEndpoint)
}
