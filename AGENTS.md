# AGENTS.md - Go-Zero Project

This file provides guidance to AI agents working on the Go-Zero project.

## Project Overview

Go-Zero is a web and RPC framework for building microservices. It provides:

- **Language**: Go (Golang) 1.21+
- **Type**: Microservices framework
- **Purpose**: Building scalable web and RPC services
- **Features**: HTTP server, RPC server, service discovery, monitoring, etc.

## Key Configuration Files

- `go.mod` - Go module definition and dependencies
- `go.sum` - Dependency checksums
- `readme.md` - Main project documentation
- `readme-cn.md` - Chinese documentation
- `CONTRIBUTING.md` - Contribution guidelines
- `code-of-conduct.md` - Code of conduct

## Build and Test Commands

### Installation
```bash
# Install the framework
go get github.com/zeromicro/go-zero

# Install dependencies
go mod tidy
```

### Development
```bash
# Build the project
go build ./...

# Build specific components
go build ./core/...      # Core components
go build ./rest/...      # REST components
go build ./zrpc/...      # RPC components
go build ./gateway/...   # Gateway components
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific component
go test ./core/...
go test ./rest/...

# Run manual tests (if any)
go test -tags manual ./...
```

### Code Quality
```bash
# Format code
gofmt -w .

# Check for formatting issues
gofmt -d .

# Run linter (if available)
golangci-lint run
```

## Project Structure

```
core/          # Core framework components
  cache/       # Caching functionality
  conf/        # Configuration management
  logging/     # Logging functionality
  monitoring/  # Monitoring and metrics
  service/     # Service management
  sync/        # Synchronization primitives
gateway/      # API gateway components
internal/      # Internal packages
mcp/           # Management control protocol
rest/          # REST API components
  handler/     # HTTP handlers
  router/      # Routing functionality
tools/         # Development tools
zrpc/          # RPC components
  client/      # RPC client
  server/      # RPC server
demo/          # Example applications
```

## Code Style Guidelines

- **Go Standards**: Follow official Go code review comments
- **Formatting**: Use `gofmt` for consistent formatting
- **Naming**: Use camelCase for variables, PascalCase for exported types
- **Error Handling**: Explicit error handling (no panic for expected errors)
- **Documentation**: Add godoc comments for exported functions/types
- **Consistency**: Maintain consistent patterns across the framework

## Testing Instructions

- **Unit Tests**: Each component has its own tests
- **Integration Tests**: Test component interactions
- **Manual Tests**: Some tests may require manual execution
- **Coverage**: Aim for high test coverage
- **Mocking**: Use interfaces for mocking dependencies

## Framework Components

### Core Components
- **Cache**: Distributed caching support
- **Configuration**: Flexible configuration management
- **Logging**: Structured logging
- **Monitoring**: Metrics and tracing
- **Service Discovery**: Service registration and discovery

### REST Components
- **Routing**: HTTP request routing
- **Handlers**: Request handling
- **Middleware**: Request processing pipeline
- **Validation**: Request validation

### RPC Components
- **Client**: RPC client implementation
- **Server**: RPC server implementation
- **Code Generation**: Protobuf code generation
- **Service Discovery**: RPC service discovery

### Gateway Components
- **API Gateway**: Request routing and aggregation
- **Load Balancing**: Traffic distribution
- **Circuit Breaking**: Fault tolerance
- **Rate Limiting**: Traffic control

## Security Considerations

- **Input Validation**: Validate all external inputs
- **Authentication**: Secure authentication mechanisms
- **Authorization**: Proper access control
- **Error Handling**: Don't expose sensitive information
- **Dependency Management**: Regularly update dependencies
- **Network Security**: Secure network communications

## Performance Considerations

- **Concurrency**: Efficient goroutine usage
- **Memory Management**: Minimize allocations
- **I/O Operations**: Optimize network and disk I/O
- **Caching**: Effective caching strategies
- **Benchmarking**: Performance testing and optimization

## Microservices Best Practices

- **Service Decomposition**: Proper service boundaries
- **API Design**: RESTful and RPC API design
- **Resilience**: Fault tolerance patterns
- **Observability**: Monitoring and logging
- **Deployment**: Containerization and orchestration

## Git Conventions

- **Commit Messages**: Clear, descriptive commit messages
- **Branching**: Use feature branches for new development
- **Pull Requests**: Required for merging to main branch
- **Tags**: Use semantic versioning for releases
- **Changelog**: Maintain changelog for significant changes

## CI/CD

- **GitHub Actions**: Configured in `.github/workflows/`
- **Automated Testing**: Runs on every push/PR
- **Build Verification**: Ensures all components build
- **Test Coverage**: Reports test coverage metrics
- **Release Process**: Automated release workflows

## Documentation

- **readme.md**: Main project documentation
- **readme-cn.md**: Chinese documentation
- **Godoc**: Use godoc comments for inline documentation
- **Examples**: `demo/` directory contains example applications
- **Tutorials**: Step-by-step guides for using the framework

## Dependency Management

- **Go Modules**: Uses Go modules for dependency management
- **Version Pinning**: Specific versions for stability
- **Updates**: Regularly update dependencies
- **Compatibility**: Ensure backward compatibility

## Microservices Architecture

- **Service-Oriented**: Design for service-oriented architecture
- **Decentralized**: Independent service development
- **Scalable**: Horizontal scaling capabilities
- **Resilient**: Fault tolerance and recovery
- **Observable**: Comprehensive monitoring and logging

## API Design Principles

- **RESTful**: Follow REST principles for HTTP APIs
- **RPC**: Efficient RPC communication
- **Versioning**: API versioning strategies
- **Documentation**: Clear API documentation
- **Consistency**: Consistent API design across services

## Deployment Considerations

- **Containerization**: Docker support
- **Orchestration**: Kubernetes integration
- **Configuration**: Environment-specific configuration
- **Scaling**: Horizontal scaling strategies
- **Monitoring**: Production monitoring setup

## Future Enhancements

- **Additional Features**: Support for more protocols and patterns
- **Performance**: Optimize framework components
- **Monitoring**: Enhanced observability features
- **Documentation**: Expand examples and tutorials
- **Tooling**: Additional development tools

## Community and Ecosystem

- **Contributions**: Welcome community contributions
- **Plugins**: Support for framework extensions
- **Integrations**: Integration with other tools
- **Support**: Community support channels
- **Roadmap**: Clear project roadmap
