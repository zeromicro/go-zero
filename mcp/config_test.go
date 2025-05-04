package mcp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestMcpConfDefaults(t *testing.T) {
	// Test default values are set correctly when unmarshalled from JSON
	jsonConfig := `name: test-service
port: 8080
mcp:
  name: test-mcp-server
  version: 1.0.0
`

	var c McpConf
	err := conf.LoadFromYamlBytes([]byte(jsonConfig), &c)
	assert.NoError(t, err)

	// Check default values
	assert.Equal(t, "test-mcp-server", c.Mcp.Name)
	assert.Equal(t, "1.0.0", c.Mcp.Version, "Default version should be 1.0.0")
	assert.Equal(t, "2024-11-05", c.Mcp.ProtocolVersion, "Default protocol version should be 2024-11-05")
	assert.Equal(t, "/sse", c.Mcp.SseEndpoint, "Default SSE endpoint should be /sse")
	assert.Equal(t, "/message", c.Mcp.MessageEndpoint, "Default message endpoint should be /message")
	assert.Equal(t, 30*time.Second, c.Mcp.MessageTimeout, "Default message timeout should be 30s")
}

func TestMcpConfCustomValues(t *testing.T) {
	// Test custom values can be set
	jsonConfig := `{
		"Name": "test-service",
		"Port": 8080,
		"Mcp": {
			"Name": "test-mcp-server",
			"Version": "2.0.0",
			"ProtocolVersion": "2025-01-01",
			"BaseUrl": "http://example.com",
			"SseEndpoint": "/custom-sse",
			"MessageEndpoint": "/custom-message",
			"Cors": ["http://localhost:3000", "http://example.com"],
			"MessageTimeout": "60s"
		}
	}`

	var c McpConf
	err := conf.LoadFromJsonBytes([]byte(jsonConfig), &c)
	assert.NoError(t, err)

	// Check custom values
	assert.Equal(t, "test-mcp-server", c.Mcp.Name, "Name should be inherited from RestConf")
	assert.Equal(t, "2.0.0", c.Mcp.Version, "Version should be customizable")
	assert.Equal(t, "2025-01-01", c.Mcp.ProtocolVersion, "Protocol version should be customizable")
	assert.Equal(t, "http://example.com", c.Mcp.BaseUrl, "BaseUrl should be customizable")
	assert.Equal(t, "/custom-sse", c.Mcp.SseEndpoint, "SSE endpoint should be customizable")
	assert.Equal(t, "/custom-message", c.Mcp.MessageEndpoint, "Message endpoint should be customizable")
	assert.Equal(t, []string{"http://localhost:3000", "http://example.com"}, c.Mcp.Cors, "CORS settings should be customizable")
	assert.Equal(t, 60*time.Second, c.Mcp.MessageTimeout, "Tool timeout should be customizable")
}
