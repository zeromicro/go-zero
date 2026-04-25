package mcp

import sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"

// SDKServer exposes the underlying official SDK server for advanced use cases.
// It returns nil for non-go-zero MCP server implementations.
func SDKServer(server McpServer) *sdkmcp.Server {
	if impl, ok := server.(*mcpServerImpl); ok {
		return impl.mcpServer
	}

	return nil
}

// RemoveTools removes tools from the underlying SDK server.
func RemoveTools(server McpServer, names ...string) {
	if sdkServer := SDKServer(server); sdkServer != nil {
		sdkServer.RemoveTools(names...)
	}
}
