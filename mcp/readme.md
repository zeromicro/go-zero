# Model Context Protocol (MCP) SDK Implementation

## Overview
This package implements a Model Context Protocol (MCP) server in Go that facilitates real-time communication between AI models and clients using Server-Sent Events (SSE). The implementation provides a framework for building AI-assisted applications with bidirectional communication capabilities.

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

To create and use an MCP server, see the examples directory for practical implementation examples including:
- Tool registration and execution
- Static and dynamic prompt creation
- Resource handling with proper URI identification
- Embedded resources in prompt responses
- Client connection management