# GitHub Copilot Instructions for go-zero

This document provides guidelines for GitHub Copilot when assisting with development in the go-zero project.

## Project Overview

go-zero is a web and RPC framework with lots of built-in engineering practices designed to ensure the stability of busy services with resilience design. It has been serving sites with tens of millions of users for years.

### Key Architecture Components

- **REST API framework** (`rest/`) - HTTP service framework with middleware support
- **RPC framework** (`zrpc/`) - gRPC-based RPC framework with service discovery
- **Core utilities** (`core/`) - Foundational components including:
  - Circuit breakers, rate limiters, load shedding
  - Caching, stores (Redis, MongoDB, SQL)
  - Concurrency control, metrics, tracing
  - Configuration management
- **Code generation tool** (`tools/goctl/`) - CLI tool for generating code from API files

## Coding Standards and Conventions

### Code Style

1. **Follow Go conventions**: Use `gofmt` for formatting, follow effective Go practices
2. **Package naming**: Use lowercase, single-word package names when possible
3. **Error handling**: Always handle errors explicitly, use `errorx.BatchError` for multiple errors
4. **Context propagation**: Always pass `context.Context` as the first parameter for functions that may block
5. **Configuration structures**: Use struct tags with JSON annotations and default values

Example configuration pattern:
```go
type Config struct {
    Host     string `json:",default=0.0.0.0"`
    Port     int    `json:",default=8080"`
    Timeout  int    `json:",default=3000"`
    Optional string `json:",optional"`
}
```

### Interface Design

1. **Small interfaces**: Follow Go's preference for small, focused interfaces
2. **Context methods**: Provide both context and non-context versions of methods
3. **Options pattern**: Use functional options for complex configuration

Example:
```go
func (c *Client) Get(key string, val any) error {
    return c.GetCtx(context.Background(), key, val)
}

func (c *Client) GetCtx(ctx context.Context, key string, val any) error {
    // implementation
}
```

### Testing Patterns

1. **Test file naming**: Use `*_test.go` suffix
2. **Test function naming**: Use `TestFunctionName` pattern
3. **Use testify/assert**: Prefer `assert` package for assertions
4. **Table-driven tests**: Use table-driven tests for multiple scenarios
5. **Mock interfaces**: Use `go.uber.org/mock` for mocking
6. **Test helpers**: Use `redistest`, `mongtest` helpers for database testing

Example test pattern:
```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {"valid case", "input", "output", false},
        {"error case", "bad", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := SomeFunction(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Framework-Specific Guidelines

### REST API Development

1. **API Definition**: Use `.api` files to define REST APIs
2. **Handler pattern**: Separate business logic into logic packages
3. **Middleware**: Use built-in middlewares (tracing, logging, metrics, recovery)
4. **Response handling**: Use `httpx.WriteJson` for JSON responses
5. **Error handling**: Use `httpx.Error` for HTTP error responses

### RPC Development

1. **Protocol Buffers**: Use protobuf for service definitions
2. **Service discovery**: Integrate with etcd for service registration
3. **Load balancing**: Use built-in load balancing strategies
4. **Interceptors**: Implement interceptors for cross-cutting concerns

### Database Operations

1. **SQL operations**: Use `sqlx` package for database operations
2. **Caching**: Implement caching patterns with `cache` package
3. **Transactions**: Use proper transaction handling
4. **Connection pooling**: Configure appropriate connection pools

Example cache pattern:
```go
err := c.QueryRowCtx(ctx, &dest, key, func(ctx context.Context, conn sqlx.SqlConn) error {
    return conn.QueryRowCtx(ctx, &dest, query, args...)
})
```

### Configuration Management

1. **YAML configuration**: Use YAML for configuration files
2. **Environment variables**: Support environment variable overrides
3. **Validation**: Include proper validation for configuration parameters
4. **Sensible defaults**: Provide reasonable default values

## Error Handling Best Practices

1. **Wrap errors**: Use `fmt.Errorf` with `%w` verb to wrap errors
2. **Custom errors**: Define custom error types when needed
3. **Error logging**: Log errors appropriately with context
4. **Graceful degradation**: Implement fallback mechanisms

## Performance Considerations

1. **Resource pools**: Use connection pools and worker pools
2. **Circuit breakers**: Implement circuit breaker patterns for external calls
3. **Rate limiting**: Apply rate limiting to protect services
4. **Load shedding**: Implement adaptive load shedding
5. **Metrics**: Add appropriate metrics and monitoring

## Security Guidelines

1. **Input validation**: Validate all input parameters
2. **SQL injection prevention**: Use parameterized queries
3. **Authentication**: Implement proper JWT token handling
4. **HTTPS**: Support TLS/HTTPS configurations
5. **CORS**: Configure CORS appropriately for web APIs

## Documentation Standards

1. **Package documentation**: Include package-level documentation
2. **Function documentation**: Document exported functions with examples
3. **API documentation**: Maintain API documentation in sync
4. **README updates**: Update README for significant changes

## Common Patterns to Follow

### Service Configuration
```go
type ServiceConf struct {
    Name string
    Log  logx.LogConf
    Mode string `json:",default=pro,options=[dev,test,pre,pro]"`
    // ... other common fields
}
```

### Middleware Implementation
```go
func SomeMiddleware() rest.Middleware {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            // Pre-processing
            next.ServeHTTP(w, r)
            // Post-processing
        }
    }
}
```

### Resource Management
Always implement proper resource cleanup using defer and context cancellation.

## Build and Test Commands

- Build: `go build ./...`
- Test: `go test ./...`  
- Test with race detection: `go test -race ./...`
- Format: `gofmt -w .`
- Generate code: `goctl api go -api *.api -dir .`

Remember to run tests and ensure all checks pass before submitting changes. The project emphasizes high quality, performance, and reliability, so these should be primary considerations in all development work.