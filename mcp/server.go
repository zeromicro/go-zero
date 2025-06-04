package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

func NewMcpServer(c McpConf) McpServer {
	var server *rest.Server
	if len(c.Mcp.Cors) == 0 {
		server = rest.MustNewServer(c.RestConf)
	} else {
		server = rest.MustNewServer(c.RestConf, rest.WithCors(c.Mcp.Cors...))
	}

	if len(c.Mcp.Name) == 0 {
		c.Mcp.Name = c.Name
	}
	if len(c.Mcp.BaseUrl) == 0 {
		c.Mcp.BaseUrl = fmt.Sprintf("http://localhost:%d", c.Port)
	}

	s := &sseMcpServer{
		conf:      c,
		server:    server,
		clients:   make(map[string]*mcpClient),
		tools:     make(map[string]Tool),
		prompts:   make(map[string]Prompt),
		resources: make(map[string]Resource),
	}

	// SSE endpoint for real-time updates
	s.server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    s.conf.Mcp.SseEndpoint,
		Handler: s.handleSSE,
	}, rest.WithSSE(), rest.WithTimeout(c.Mcp.SseTimeout))

	// JSON-RPC message endpoint for regular requests
	s.server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    s.conf.Mcp.MessageEndpoint,
		Handler: s.handleRequest,
	}, rest.WithTimeout(c.Mcp.MessageTimeout))

	return s
}

// RegisterPrompt registers a new prompt with the server
func (s *sseMcpServer) RegisterPrompt(prompt Prompt) {
	s.promptsLock.Lock()
	s.prompts[prompt.Name] = prompt
	s.promptsLock.Unlock()
	// Notify clients about the new prompt
	s.broadcast(eventPromptsListChanged, map[string][]Prompt{keyPrompts: {prompt}})
}

// RegisterResource registers a new resource with the server
func (s *sseMcpServer) RegisterResource(resource Resource) {
	s.resourcesLock.Lock()
	s.resources[resource.URI] = resource
	s.resourcesLock.Unlock()
	// Notify clients about the new resource
	s.broadcast(eventResourcesListChanged, map[string][]Resource{keyResources: {resource}})
}

// RegisterTool registers a new tool with the server
func (s *sseMcpServer) RegisterTool(tool Tool) error {
	if tool.Handler == nil {
		return fmt.Errorf("tool '%s' has no handler function", tool.Name)
	}

	s.toolsLock.Lock()
	s.tools[tool.Name] = tool
	s.toolsLock.Unlock()
	// Notify clients about the new tool
	s.broadcast(eventToolsListChanged, map[string][]Tool{keyTools: {tool}})
	return nil
}

// Start implements McpServer.
func (s *sseMcpServer) Start() {
	s.server.Start()
}

func (s *sseMcpServer) Stop() {
	s.server.Stop()
}

// broadcast sends a message to all connected clients
// It uses Server-Sent Events (SSE) format for real-time communication
func (s *sseMcpServer) broadcast(event string, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logx.Errorf("Failed to marshal broadcast data: %v", err)
		return
	}

	// Lock only while reading the clients map
	s.clientsLock.Lock()
	clients := make([]*mcpClient, 0, len(s.clients))
	for _, client := range s.clients {
		clients = append(clients, client)
	}
	s.clientsLock.Unlock()

	clientCount := len(clients)
	if clientCount == 0 {
		return
	}

	logx.Infof("Broadcasting event '%s' to %d clients", event, clientCount)

	// Use CRLF line endings as per SSE specification
	message := fmt.Sprintf("event: %s\r\ndata: %s\r\n\r\n", event, string(jsonData))

	// Send messages without holding the lock
	for _, client := range clients {
		select {
		case client.channel <- message:
			// Message sent successfully
		default:
			// Channel buffer is full, log warning and continue
			logx.Errorf("Client channel buffer full, dropping message for client %s", client.id)
		}
	}
}

// cleanupClient removes a client from the active clients map
func (s *sseMcpServer) cleanupClient(sessionId string) {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	if client, exists := s.clients[sessionId]; exists {
		// Close the channel to signal any goroutines waiting on it
		close(client.channel)
		// Remove from active clients
		delete(s.clients, sessionId)
		logx.Infof("Cleaned up client %s (remaining clients: %d)", sessionId, len(s.clients))
	}
}

// handleRequest handles MCP JSON-RPC requests
func (s *sseMcpServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	// Extract sessionId from query parameters
	sessionId := r.URL.Query().Get(sessionIdKey)
	if len(sessionId) == 0 {
		http.Error(w, fmt.Sprintf("Missing %s", sessionIdKey), http.StatusBadRequest)
		return
	}

	// Check if the client with this sessionId exists
	s.clientsLock.Lock()
	client, exists := s.clients[sessionId]
	s.clientsLock.Unlock()

	if !exists {
		http.Error(w, fmt.Sprintf("Invalid or expired %s", sessionIdKey), http.StatusBadRequest)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)

	// For notification methods (no ID), we don't send a response
	isNotification := req.ID == 0

	// Special handling for initialization sequence
	// Always allow initialize and notifications/initialized regardless of client state
	if req.Method == methodInitialize {
		logx.Infof("Processing initialize request with ID: %d", req.ID)
		s.processInitialize(r.Context(), client, req)
		logx.Infof("Sent initialize response for ID: %d, waiting for notifications/initialized", req.ID)
		return
	} else if req.Method == methodNotificationsInitialized {
		// Handle initialized notification
		logx.Info("Received notifications/initialized notification")
		if !isNotification {
			s.sendErrorResponse(r.Context(), client, req.ID,
				"Method should be used as a notification", errCodeInvalidRequest)
			return
		}
		s.processNotificationInitialized(client)
		return
	} else if !client.initialized && req.Method != methodNotificationsCancelled {
		// Block most requests until client is initialized (except for cancellations)
		s.sendErrorResponse(r.Context(), client, req.ID,
			"Client not fully initialized, waiting for notifications/initialized",
			errCodeClientNotInitialized)
		return
	}

	// Process normal requests only after initialization
	switch req.Method {
	case methodToolsCall:
		logx.Infof("Received tools call request with ID: %d", req.ID)
		s.processToolCall(r.Context(), client, req)
		logx.Infof("Sent tools call response for ID: %d", req.ID)
	case methodToolsList:
		logx.Infof("Processing tools/list request with ID: %d", req.ID)
		s.processListTools(r.Context(), client, req)
		logx.Infof("Sent tools/list response for ID: %d", req.ID)
	case methodPromptsList:
		logx.Infof("Processing prompts/list request with ID: %d", req.ID)
		s.processListPrompts(r.Context(), client, req)
		logx.Infof("Sent prompts/list response for ID: %d", req.ID)
	case methodPromptsGet:
		logx.Infof("Processing prompts/get request with ID: %d", req.ID)
		s.processGetPrompt(r.Context(), client, req)
		logx.Infof("Sent prompts/get response for ID: %d", req.ID)
	case methodResourcesList:
		logx.Infof("Processing resources/list request with ID: %d", req.ID)
		s.processListResources(r.Context(), client, req)
		logx.Infof("Sent resources/list response for ID: %d", req.ID)
	case methodResourcesRead:
		logx.Infof("Processing resources/read request with ID: %d", req.ID)
		s.processResourcesRead(r.Context(), client, req)
		logx.Infof("Sent resources/read response for ID: %d", req.ID)
	case methodResourcesSubscribe:
		logx.Infof("Processing resources/subscribe request with ID: %d", req.ID)
		s.processResourceSubscribe(r.Context(), client, req)
		logx.Infof("Sent resources/subscribe response for ID: %d", req.ID)
	case methodPing:
		logx.Infof("Processing ping request with ID: %d", req.ID)
		s.processPing(r.Context(), client, req)
	case methodNotificationsCancelled:
		logx.Infof("Received notifications/cancelled notification: %d", req.ID)
		s.processNotificationCancelled(r.Context(), client, req)
	default:
		logx.Infof("Unknown method: %s from client: %d", req.Method, req.ID)
		s.sendErrorResponse(r.Context(), client, req.ID, "Method not found", errCodeMethodNotFound)
	}
}

// handleSSE handles Server-Sent Events connections
func (s *sseMcpServer) handleSSE(w http.ResponseWriter, r *http.Request) {
	// Generate a unique session ID for this client
	sessionId := uuid.New().String()

	// Create new client with buffered channel to prevent blocking
	client := &mcpClient{
		id:      sessionId,
		channel: make(chan string, eventChanSize),
	}

	// Add client to active clients map
	s.clientsLock.Lock()
	s.clients[sessionId] = client
	activeClients := len(s.clients)
	s.clientsLock.Unlock()

	logx.Infof("New SSE connection established for client %s (active clients: %d)",
		sessionId, activeClients)

	// Set proper SSE headers
	w.Header().Set("Transfer-Encoding", "chunked")

	// Enable streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		logx.Error("Streaming not supported by the underlying http.ResponseWriter")
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send the message endpoint URL to the client
	endpoint := fmt.Sprintf("%s%s?%s=%s",
		s.conf.Mcp.BaseUrl, s.conf.Mcp.MessageEndpoint, sessionIdKey, sessionId)

	// Format and send the endpoint message
	endpointMsg := formatSSEMessage(eventEndpoint, []byte(endpoint))
	if _, err := fmt.Fprint(w, endpointMsg); err != nil {
		logx.Errorf("Failed to send endpoint message to client %s: %v", sessionId, err)
		s.cleanupClient(sessionId)
		return
	}
	flusher.Flush()

	// Set up keep-alive ping and client cleanup
	ticker := time.NewTicker(pingInterval.Load())
	defer func() {
		ticker.Stop()
		s.cleanupClient(sessionId)
		logx.Infof("SSE connection closed for client %s", sessionId)
	}()

	// Message processing loop
	for {
		select {
		case message, ok := <-client.channel:
			if !ok {
				// Channel was closed, end connection
				logx.Infof("Client channel was closed for %s", sessionId)
				return
			}

			// Write message to the response
			if _, err := fmt.Fprint(w, message); err != nil {
				logx.Infof("Failed to write message to client %s: %v", sessionId, err)
				return
			}
			flusher.Flush()
		case <-ticker.C:
			// Send keep-alive ping to maintain connection
			ping := fmt.Sprintf(`{"type":"ping","timestamp":"%s"}`, time.Now().String())
			pingMsg := formatSSEMessage("ping", []byte(ping))
			if _, err := fmt.Fprint(w, pingMsg); err != nil {
				logx.Errorf("Failed to send ping to client %s, closing connection: %v", sessionId, err)
				return
			}
			flusher.Flush()
		case <-r.Context().Done():
			// Client disconnected or request was canceled or timed out
			logx.Infof("Client %s disconnected: context done", sessionId)
			return
		}
	}
}

// processInitialize processes the initialize request
func (s *sseMcpServer) processInitialize(ctx context.Context, client *mcpClient, req Request) {
	// Create a proper JSON-RPC response that preserves the client's request ID
	result := initializationResponse{
		ProtocolVersion: s.conf.Mcp.ProtocolVersion,
		Capabilities: capabilities{
			Prompts: struct {
				ListChanged bool `json:"listChanged"`
			}{
				ListChanged: true,
			},
			Resources: struct {
				Subscribe   bool `json:"subscribe"`
				ListChanged bool `json:"listChanged"`
			}{
				Subscribe:   true,
				ListChanged: true,
			},
			Tools: struct {
				ListChanged bool `json:"listChanged"`
			}{
				ListChanged: true,
			},
		},
		ServerInfo: serverInfo{
			Name:    s.conf.Mcp.Name,
			Version: s.conf.Mcp.Version,
		},
	}

	// Mark client as initialized
	client.initialized = true

	// Send response with client's original request ID
	s.sendResponse(ctx, client, req.ID, result)
}

// processListTools processes the tools/list request
func (s *sseMcpServer) processListTools(ctx context.Context, client *mcpClient, req Request) {
	// Extract pagination params if any
	var nextCursor string
	var progressToken any

	// Extract meta data including progress token
	if req.Params != nil {
		var metaParams struct {
			Cursor string `json:"cursor"`
			Meta   struct {
				ProgressToken any `json:"progressToken"`
			} `json:"_meta"`
		}
		if err := json.Unmarshal(req.Params, &metaParams); err == nil {
			if len(metaParams.Cursor) > 0 {
				nextCursor = metaParams.Cursor
			}
			progressToken = metaParams.Meta.ProgressToken
		}
	}

	var toolsList []Tool
	s.toolsLock.Lock()
	for _, tool := range s.tools {
		if len(tool.InputSchema.Type) == 0 {
			tool.InputSchema.Type = ContentTypeObject
		}
		toolsList = append(toolsList, tool)
	}
	s.toolsLock.Unlock()

	result := ListToolsResult{
		PaginatedResult: PaginatedResult{
			Result:     Result{},
			NextCursor: Cursor(nextCursor),
		},
		Tools: toolsList,
	}

	// Add meta information if progress token was provided
	if progressToken != nil {
		result.Result.Meta = map[string]any{
			progressTokenKey: progressToken,
		}
	}

	s.sendResponse(ctx, client, req.ID, result)
}

// processListPrompts processes the prompts/list request
func (s *sseMcpServer) processListPrompts(ctx context.Context, client *mcpClient, req Request) {
	// Extract pagination params if any
	var nextCursor string
	if req.Params != nil {
		var cursorParams struct {
			Cursor string `json:"cursor"`
		}
		if err := json.Unmarshal(req.Params, &cursorParams); err == nil && cursorParams.Cursor != "" {
			// If we have a valid cursor, we could use it for pagination
			// For now, we're not actually implementing pagination, so this is just
			// to show how it would be extracted from the request
			_ = cursorParams.Cursor
		}
	}

	// Prepare prompt list
	var promptsList []Prompt
	s.promptsLock.Lock()
	for _, prompt := range s.prompts {
		promptsList = append(promptsList, prompt)
	}
	s.promptsLock.Unlock()

	// In a real implementation, you'd handle pagination here
	// For now, we'll return all prompts at once
	result := struct {
		Prompts    []Prompt  `json:"prompts"`
		NextCursor string    `json:"nextCursor,omitempty"`
		Meta       *struct{} `json:"_meta,omitempty"`
	}{
		Prompts:    promptsList,
		NextCursor: nextCursor,
	}

	s.sendResponse(ctx, client, req.ID, result)
}

// processListResources processes the resources/list request
func (s *sseMcpServer) processListResources(ctx context.Context, client *mcpClient, req Request) {
	// Extract pagination params if any
	var nextCursor string
	var progressToken any

	// Extract meta information including progress token if available
	if req.Params != nil {
		var metaParams PaginatedParams
		if err := json.Unmarshal(req.Params, &metaParams); err == nil {
			if len(metaParams.Cursor) > 0 {
				nextCursor = metaParams.Cursor
			}
			progressToken = metaParams.Meta.ProgressToken
		}
	}

	var resourcesList []Resource
	s.resourcesLock.Lock()
	for _, resource := range s.resources {
		// Create a copy without the handler function which shouldn't be sent to clients
		resourceCopy := Resource{
			URI:         resource.URI,
			Name:        resource.Name,
			Description: resource.Description,
			MimeType:    resource.MimeType,
		}
		resourcesList = append(resourcesList, resourceCopy)
	}
	s.resourcesLock.Unlock()

	// Create proper ResourcesListResult according to MCP specification
	result := ResourcesListResult{
		PaginatedResult: PaginatedResult{
			Result:     Result{},
			NextCursor: Cursor(nextCursor),
		},
		Resources: resourcesList,
	}

	// Add meta information if progress token was provided
	if progressToken != nil {
		result.Result.Meta = map[string]any{
			progressTokenKey: progressToken,
		}
	}

	s.sendResponse(ctx, client, req.ID, result)
}

// processGetPrompt processes the prompts/get request
func (s *sseMcpServer) processGetPrompt(ctx context.Context, client *mcpClient, req Request) {
	type GetPromptParams struct {
		Name      string            `json:"name"`
		Arguments map[string]string `json:"arguments,omitempty"`
	}

	var params GetPromptParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendErrorResponse(ctx, client, req.ID, "Invalid parameters", errCodeInvalidParams)
		return
	}

	// Check if prompt exists
	s.promptsLock.Lock()
	prompt, exists := s.prompts[params.Name]
	s.promptsLock.Unlock()
	if !exists {
		message := fmt.Sprintf("Prompt '%s' not found", params.Name)
		s.sendErrorResponse(ctx, client, req.ID, message, errCodeInvalidParams)
		return
	}

	logx.Infof("Processing prompt request: %s with %d arguments", prompt.Name, len(params.Arguments))

	// Validate required arguments
	missingArgs := validatePromptArguments(prompt, params.Arguments)
	if len(missingArgs) > 0 {
		message := fmt.Sprintf("Missing required arguments: %s", strings.Join(missingArgs, ", "))
		s.sendErrorResponse(ctx, client, req.ID, message, errCodeInvalidParams)
		return
	}

	// Ensure arguments are initialized to an empty map if nil
	if params.Arguments == nil {
		params.Arguments = make(map[string]string)
	}
	args := params.Arguments

	// Generate messages using handler or static content
	var messages []PromptMessage
	var err error

	if prompt.Handler != nil {
		// Use dynamic handler to generate messages
		messages, err = prompt.Handler(ctx, args)
		if err != nil {
			logx.Errorf("Error from prompt handler: %v", err)
			s.sendErrorResponse(ctx, client, req.ID,
				fmt.Sprintf("Error generating prompt content: %v", err), errCodeInternalError)
			return
		}
	} else {
		// No handler, generate messages from static content
		var messageText string
		if len(prompt.Content) > 0 {
			messageText = prompt.Content

			// Apply argument substitutions to static content
			for key, value := range args {
				placeholder := fmt.Sprintf("{{%s}}", key)
				messageText = strings.Replace(messageText, placeholder, value, -1)
			}
		}

		// Create a single user message with the content
		messages = []PromptMessage{
			{
				Role: RoleUser,
				Content: TextContent{
					Text: messageText,
				},
			},
		}
	}

	// Construct the response according to MCP spec
	result := struct {
		Description string          `json:"description,omitempty"`
		Messages    []PromptMessage `json:"messages"`
	}{
		Description: prompt.Description,
		Messages:    toTypedPromptMessages(messages),
	}

	s.sendResponse(ctx, client, req.ID, result)
}

// processToolCall processes the tools/call request
func (s *sseMcpServer) processToolCall(ctx context.Context, client *mcpClient, req Request) {
	var toolCallParams struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments,omitempty"`
		Meta      struct {
			ProgressToken any `json:"progressToken"`
		} `json:"_meta,omitempty"`
	}

	// Handle different types of req.Params
	// If it's a RawMessage (JSON), unmarshal it
	if err := json.Unmarshal(req.Params, &toolCallParams); err != nil {
		logx.Errorf("Failed to unmarshal tool call params: %v", err)
		s.sendErrorResponse(ctx, client, req.ID, "Invalid tool call parameters", errCodeInvalidParams)
		return
	}

	// Extract progress token if available
	progressToken := toolCallParams.Meta.ProgressToken

	// Find the requested tool
	s.toolsLock.Lock()
	tool, exists := s.tools[toolCallParams.Name]
	s.toolsLock.Unlock()
	if !exists {
		s.sendErrorResponse(ctx, client, req.ID, fmt.Sprintf("Tool '%s' not found",
			toolCallParams.Name), errCodeInvalidParams)
		return
	}

	// Log parameters before execution
	logx.Infof("Executing tool '%s' with arguments: %#v", toolCallParams.Name, toolCallParams.Arguments)

	// Execute the tool handler with timeout handling
	var result any
	var err error

	// Create a channel to receive the result
	// make sure to have 1 size buffer to avoid channel leak if timeout
	resultCh := make(chan struct {
		result any
		err    error
	}, 1)

	// Execute the tool handler in a goroutine
	go func() {
		toolResult, toolErr := tool.Handler(ctx, toolCallParams.Arguments)
		resultCh <- struct {
			result any
			err    error
		}{
			result: toolResult,
			err:    toolErr,
		}
	}()

	// Wait for either the result or a timeout
	select {
	case res := <-resultCh:
		result = res.result
		err = res.err
	case <-ctx.Done():
		// Handle request timeout
		logx.Errorf("Tool execution timed out after %v: %s", s.conf.Mcp.MessageTimeout, toolCallParams.Name)
		s.sendErrorResponse(ctx, client, req.ID, "Tool execution timed out", errCodeTimeout)
		return
	}

	// Create the base result structure with metadata
	callToolResult := CallToolResult{
		Result:  Result{},
		Content: []any{},
		IsError: false,
	}

	// Add meta information if progress token was provided
	if progressToken != nil {
		callToolResult.Result.Meta = map[string]any{
			progressTokenKey: progressToken,
		}
	}

	// Check if there was an error during tool execution
	if err != nil {
		// According to the spec, for tool-level errors (as opposed to protocol-level errors),
		// we should report them inside the result with isError=true
		logx.Errorf("Tool execution reported error: %v", err)

		callToolResult.Content = []any{
			TextContent{
				Text: fmt.Sprintf("Error: %v", err),
			},
		}
		callToolResult.IsError = true
		s.sendResponse(ctx, client, req.ID, callToolResult)
		return
	}

	// Format the response according to the CallToolResult schema
	switch v := result.(type) {
	case string:
		// Simple string becomes text content
		callToolResult.Content = append(callToolResult.Content, TextContent{
			Text: v,
			Annotations: &Annotations{
				Audience: []RoleType{RoleUser, RoleAssistant},
			},
		})
	case map[string]any:
		// JSON-like object becomes formatted JSON text
		jsonStr, err := json.Marshal(v)
		if err != nil {
			jsonStr = []byte(err.Error())
		}
		callToolResult.Content = append(callToolResult.Content, TextContent{
			Text: string(jsonStr),
			Annotations: &Annotations{
				Audience: []RoleType{RoleUser, RoleAssistant},
			},
		})
	case TextContent:
		callToolResult.Content = append(callToolResult.Content, v)
	case ImageContent:
		callToolResult.Content = append(callToolResult.Content, v)
	case []any:
		callToolResult.Content = v
	case ToolResult:
		// Handle legacy ToolResult type
		switch v.Type {
		case ContentTypeText:
			callToolResult.Content = append(callToolResult.Content, TextContent{
				Text: fmt.Sprintf("%v", v.Content),
				Annotations: &Annotations{
					Audience: []RoleType{RoleUser, RoleAssistant},
				},
			})
		case ContentTypeImage:
			if imgData, ok := v.Content.(map[string]any); ok {
				callToolResult.Content = append(callToolResult.Content, ImageContent{
					Data:     fmt.Sprintf("%v", imgData["data"]),
					MimeType: fmt.Sprintf("%v", imgData["mimeType"]),
				})
			}
		default:
			callToolResult.Content = append(callToolResult.Content, TextContent{
				Text: fmt.Sprintf("%v", v.Content),
				Annotations: &Annotations{
					Audience: []RoleType{RoleUser, RoleAssistant},
				},
			})
		}
	default:
		// For any other type, convert to string
		callToolResult.Content = append(callToolResult.Content, TextContent{
			Text: fmt.Sprintf("%v", v),
			Annotations: &Annotations{
				Audience: []RoleType{RoleUser, RoleAssistant},
			},
		})
	}

	callToolResult.Content = toTypedContents(callToolResult.Content)
	logx.Infof("Tool call result: %#v", callToolResult)

	s.sendResponse(ctx, client, req.ID, callToolResult)
}

// processResourcesRead processes the resources/read request
func (s *sseMcpServer) processResourcesRead(ctx context.Context, client *mcpClient, req Request) {
	var params ResourceReadParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendErrorResponse(ctx, client, req.ID, "Invalid parameters", errCodeInvalidParams)
		return
	}

	// Find resource that matches the URI
	s.resourcesLock.Lock()
	resource, exists := s.resources[params.URI]
	s.resourcesLock.Unlock()

	if !exists {
		s.sendErrorResponse(ctx, client, req.ID, fmt.Sprintf("Resource with URI '%s' not found",
			params.URI), errCodeResourceNotFound)
		return
	}

	// If no handler is provided, return an empty content array
	if resource.Handler == nil {
		result := ResourceReadResult{
			Contents: []ResourceContent{
				{
					URI:      params.URI,
					MimeType: resource.MimeType,
					Text:     "",
				},
			},
		}
		s.sendResponse(ctx, client, req.ID, result)
		return
	}

	// Execute the resource handler
	content, err := resource.Handler(ctx)
	if err != nil {
		s.sendErrorResponse(ctx, client, req.ID, fmt.Sprintf("Error reading resource: %v", err),
			errCodeInternalError)
		return
	}

	// Ensure the URI is set if not already provided by the handler
	if len(content.URI) == 0 {
		content.URI = params.URI
	}

	// Ensure MimeType is set if available from the resource definition
	if len(content.MimeType) == 0 && len(resource.MimeType) > 0 {
		content.MimeType = resource.MimeType
	}

	// Create response with contents from the handler
	// The MCP specification requires a contents array
	result := ResourceReadResult{
		Contents: []ResourceContent{content},
	}

	s.sendResponse(ctx, client, req.ID, result)
}

// processResourceSubscribe processes the resources/subscribe request
func (s *sseMcpServer) processResourceSubscribe(ctx context.Context, client *mcpClient, req Request) {
	var params ResourceSubscribeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendErrorResponse(ctx, client, req.ID, "Invalid parameters", errCodeInvalidParams)
		return
	}

	// Check if the resource exists
	s.resourcesLock.Lock()
	_, exists := s.resources[params.URI]
	s.resourcesLock.Unlock()

	if !exists {
		s.sendErrorResponse(ctx, client, req.ID, fmt.Sprintf("Resource with URI '%s' not found",
			params.URI), errCodeResourceNotFound)
		return
	}

	// Send success response for the subscription
	s.sendResponse(ctx, client, req.ID, struct{}{})
}

// processNotificationCancelled processes the notifications/cancelled notification
func (s *sseMcpServer) processNotificationCancelled(ctx context.Context, client *mcpClient, req Request) {
	// Extract the requestId that was canceled
	type CancelParams struct {
		RequestId int64  `json:"requestId"`
		Reason    string `json:"reason"`
	}

	var params CancelParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		logx.Errorf("Failed to parse cancellation params: %v", err)
		return
	}

	logx.Infof("Request %d was cancelled by client. Reason: %s", params.RequestId, params.Reason)
}

// processNotificationInitialized processes the notifications/initialized notification
func (s *sseMcpServer) processNotificationInitialized(client *mcpClient) {
	// Mark the client as properly initialized
	client.initialized = true
	logx.Infof("Client %s is now fully initialized and ready for normal operations", client.id)
}

// processPing processes the ping request and responds immediately
func (s *sseMcpServer) processPing(ctx context.Context, client *mcpClient, req Request) {
	// A ping request should simply respond with an empty result to confirm the server is alive
	logx.Infof("Received ping request with ID: %d", req.ID)

	// Send an empty response with client's original request ID
	s.sendResponse(ctx, client, req.ID, struct{}{})
}

// sendErrorResponse sends an error response via the SSE channel
func (s *sseMcpServer) sendErrorResponse(ctx context.Context, client *mcpClient,
	id int64, message string, code int) {
	errorResponse := struct {
		JsonRpc string       `json:"jsonrpc"`
		ID      int64        `json:"id"`
		Error   errorMessage `json:"error"`
	}{
		JsonRpc: jsonRpcVersion,
		ID:      id,
		Error: errorMessage{
			Code:    code,
			Message: message,
		},
	}

	// all fields are primitive types, impossible to fail
	jsonData, _ := json.Marshal(errorResponse)
	// Use CRLF line endings as requested
	sseMessage := fmt.Sprintf("event: %s\r\ndata: %s\r\n\r\n", eventMessage, string(jsonData))
	logx.Infof("Sending error for ID %d: %s", id, sseMessage)

	// cannot receive from ctx.Done() because we're sending to the channel for SSE messages
	select {
	case client.channel <- sseMessage:
	default:
		// Channel buffer is full, log warning and continue
		logx.Infof("Client %s channel is full while sending error response with ID %d", client.id, id)
	}
}

// sendResponse sends a success response via the SSE channel
func (s *sseMcpServer) sendResponse(ctx context.Context, client *mcpClient, id int64, result any) {
	response := Response{
		JsonRpc: jsonRpcVersion,
		ID:      id,
		Result:  result,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		s.sendErrorResponse(ctx, client, id, "Failed to marshal response", errCodeInternalError)
		return
	}

	// Use CRLF line endings as requested
	sseMessage := fmt.Sprintf("event: %s\r\ndata: %s\r\n\r\n", eventMessage, string(jsonData))
	logx.Infof("Sending response for ID %d: %s", id, sseMessage)

	// cannot receive from ctx.Done() because we're sending to the channel for SSE messages
	select {
	case client.channel <- sseMessage:
	default:
		// Channel buffer is full, log warning and continue
		logx.Infof("Client %s channel is full while sending response with ID %d", client.id, id)
	}
}
