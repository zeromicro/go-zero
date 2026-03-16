# GitHub Copilot Instructions for go-zero

This document provides guidelines for GitHub Copilot when assisting with development in the go-zero project.

## Project Overview

go-zero is a web and RPC framework with lots of built-in engineering practices designed to ensure the stability of busy services with resilience design. It has been serving sites with tens of millions of users for years.

### Key Architecture Components

- **REST API framework** (`rest/`) - HTTP service framework with middleware chain support
- **RPC framework** (`zrpc/`) - gRPC-based RPC framework with etcd service discovery and p2c_ewma load balancing
- **Gateway** (`gateway/`) - API gateway supporting both HTTP and gRPC upstreams with proto-based routing
- **MCP Server** (`mcp/`) - Model Context Protocol server for AI agent integration via SSE
- **Core utilities** (`core/`) - Production-grade components:
  - Resilience: circuit breakers (`breaker/`), rate limiters (`limit/`), adaptive load shedding (`load/`)
  - Storage: SQL with cache (`stores/sqlc/`), Redis (`stores/redis/`), MongoDB (`stores/mongo/`)
  - Concurrency: MapReduce (`mr/`), worker pools (`executors/`), sync primitives (`syncx/`)
  - Observability: metrics (`metric/`), tracing (`trace/`), structured logging (`logx/`)
- **Code generation tool** (`tools/goctl/`) - CLI tool for generating Go code from `.api` and `.proto` files

## Coding Standards and Conventions

### Code Style

1. **Follow Go conventions**: Use `gofmt` for formatting, follow effective Go practices
2. **Package naming**: Use lowercase, single-word package names when possible
3. **Error handling**: Always handle errors explicitly, use `errorx.BatchError` for multiple errors
4. **Context propagation**: Always pass `context.Context` as the first parameter for functions that may block
5. **Configuration structures**: Use struct tags with JSON annotations, defaults, and validation

**Pattern**: All service configs embed `service.ServiceConf` for common fields (Name, Log, Mode, Telemetry)
```go
type Config struct {
    service.ServiceConf              // Always embed for services
    Host     string `json:",default=0.0.0.0"`
    Port     int    // Required field (no default)
    Timeout  int64  `json:",default=3000"`  // Timeouts in milliseconds
    Optional string `json:",optional"`      // Optional field
    Mode     string `json:",default=pro,options=dev|test|rt|pre|pro"`  // Validated options
}
```

**Service modes**: `dev`/`test`/`rt` disable load shedding and stats; `pre`/`pro` enable all resilience features

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

1. **API Definition**: Use `.api` files to define REST APIs with goctl codegen
2. **Handler pattern**: Separate business logic into logic packages (handlers call logic layer)
3. **Middleware chain**: Middlewares wrap via `chain.Chain` interface - use `Append()` or `Prepend()` to control order
   - Built-in middlewares (all in `rest/handler/`): tracing, logging, metrics, recovery, breaker, shedding, timeout, maxconns, maxbytes, gunzip
   - Custom middleware: `func(http.Handler) http.Handler` - call `next.ServeHTTP(w, r)` to continue chain
4. **Response handling**: Use `httpx.WriteJson(w, code, v)` for JSON responses
5. **Error handling**: Use `httpx.Error(w, err)` or `httpx.ErrorCtx(ctx, w, err)` for HTTP error responses
6. **Route registration**: Routes defined with `Method`, `Path`, and `Handler` - wildcards use `:param` syntax

### RPC Development

1. **Protocol Buffers**: Use protobuf for service definitions, generate code with goctl
2. **Service discovery**: Use etcd for dynamic service registration/discovery, or direct endpoints for static routing
3. **Load balancing**: Default is `p2c_ewma` (power of 2 choices with EWMA), configurable via `BalancerName`
4. **Client configuration**: Support `Etcd`, `Endpoints`, or `Target` - use `BuildTarget()` to construct connection string
5. **Interceptors**: Implement gRPC interceptors for cross-cutting concerns (auth, logging, metrics)
6. **Health checks**: gRPC health checks enabled by default (`Health: true`)

### Database Operations

1. **SQL operations**: Use `sqlx.SqlConn` interface - methods always end with `Ctx` for context support
2. **Caching pattern**: `stores/sqlc` provides `CachedConn` for automatic cache-aside pattern
   - `QueryRowCtx`: Query with cache key, auto-populate on cache miss
   - `ExecCtx`: Execute and delete cache keys
3. **Transactions**: Use `sqlx.SqlConn.TransactCtx()` to get transaction session
4. **Connection pooling**: Managed automatically (64 max idle/open, 1min lifetime)
5. **Test helpers**: Use `redistest.CreateRedis(t)` for Redis, SQL mocks for DB testing

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

## GitHub Issue Management

### Understanding and Categorizing Issues

When analyzing GitHub issues, consider these common categories:

1. **Bug Reports**: Stack traces, version info, reproduction steps
2. **Feature Requests**: Use case, proposed solution, alternatives
3. **Questions**: Usage, configuration, or architecture
4. **Documentation Issues**: Missing, unclear, or incorrect docs
5. **Performance Issues**: Benchmarks, profiling data, resource usage

### Issue Analysis Checklist

- Identify affected component (REST, RPC, Gateway, MCP, Core utilities, goctl)
- Check versions (go-zero, Go)
- Look for reproduction steps or code examples
- Review code snippets, logs, or stack traces
- Check if related to resilience features (breaker, load shedding, rate limiting)
- Determine production impact

### Responding to Issues

Be helpful and professional. Ask clarifying questions when needed. Reference relevant documentation and code files. Provide code examples following project conventions. Suggest workarounds when applicable.

### Chinese to English Translation

go-zero has an international user base. When encountering issues or comments written in Chinese, translate them to English to ensure all contributors can participate in discussions.

#### Translation Guidelines

1. **Update issue titles**: Edit the issue title to include English translation only
2. **Translate comments in place**: Add a comment with the English translation, followed by the original Chinese text
3. **Keep original Chinese**: After translating, include the original Chinese text in a blockquote for verification
4. **Encourage English communication**: Politely suggest users write in English for better collaboration
5. **Maintain technical accuracy**: Preserve technical terms, component names, and code exactly
6. **Translate naturally**: Avoid literal word-by-word translation; use idiomatic English
7. **Preserve formatting**: Keep markdown formatting, code blocks, and links intact
8. **Keep URLs unchanged**: Don't translate URLs or file paths

#### Common Technical Terms (Chinese → English)

- 框架 → **Framework** | 中间件 → **Middleware** | 负载均衡 → **Load Balancing**
- 熔断器 → **Circuit Breaker** | 限流 → **Rate Limiting** | 降载/过载保护 → **Load Shedding**
- 服务发现 → **Service Discovery** | 配置 → **Configuration** | 弹性/容错 → **Resilience** | 微服务 → **Microservices**

#### Translation Example

**Original Chinese Title:** `goctl 执行环境问题`
**Updated Title:** `goctl Execution Environment Issue`

**Original Chinese Comment:** `我在项目中遇到熔断器配置问题`
**Translation in Comment:**
```markdown
I encountered a circuit breaker configuration issue in my project.

> Original (原文): 我在项目中遇到熔断器配置问题
```

### Common Issue Patterns and Solutions

#### Configuration Issues
- Check `service.ServiceConf` embedding and struct tags
- Verify YAML syntax, defaults, and validation rules
- Reference: [rest/config.go](rest/config.go), [zrpc/config.go](zrpc/config.go)

#### Code Generation (goctl) Issues
- Verify `.api` or `.proto` file syntax and goctl version
- Reference: `tools/goctl/` directory

#### RPC Connection Issues
- Check etcd configuration, service discovery, and endpoints
- Verify load balancing settings (p2c_ewma)

#### Database/Cache Issues
- Verify `sqlx.SqlConn` usage with context
- Check cache key generation, invalidation, and connection pools
- Use test helpers (`redistest`, `mongtest`)

#### Performance Issues
- Check if load shedding is enabled (mode: `pre`/`pro`)
- Review circuit breaker thresholds, rate limiting, and context timeouts

### Referencing Codebase

When explaining issues, reference specific files and patterns:
- REST API: `rest/`, `rest/handler/`, `rest/httpx/`
- RPC: `zrpc/`, `zrpc/internal/`
- Core utilities: `core/breaker/`, `core/limit/`, `core/load/`, etc.
- Gateway: `gateway/`
- MCP: `mcp/`
- Code generation: `tools/goctl/`
- Examples: `adhoc/` directory contains various examples

### Encouraging Best Practices

When responding to issues, gently guide users toward:
- Proper error handling with context
- Using resilience features (breakers, rate limiters)
- Following testing patterns with table-driven tests
- Implementing proper resource cleanup
- Reading existing documentation in `docs/` and `readme.md`

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
- Code generation:
  - REST API: `goctl api go -api *.api -dir .`
  - RPC: `goctl rpc protoc *.proto --go_out=. --go-grpc_out=. --zrpc_out=.`
  - Model from SQL: `goctl model mysql datasource -url="user:pass@tcp(host:port)/db" -table="*" -dir="./model"`

## Critical Architecture Patterns

### Resilience Design Philosophy
go-zero implements defense-in-depth with multiple protection layers:
1. **Circuit Breaker** (`core/breaker`): Google SRE breaker - tracks success/failure, opens on error threshold
2. **Adaptive Load Shedding** (`core/load`): CPU-based auto-rejection when system overloaded (disabled in dev/test/rt modes)
3. **Rate Limiting** (`core/limit`): Token bucket (Redis-based) and period limiters
4. **Timeout Control**: Cascading timeouts via context - set at multiple levels (client, server, handler)

### Middleware Chain Architecture
`rest/chain` provides middleware composition:
```go
// Middleware signature
type Middleware func(http.Handler) http.Handler

// Chain operations
chain := chain.New(m1, m2)
chain.Append(m3)    // Adds to end: m1 -> m2 -> m3
chain.Prepend(m0)   // Adds to start: m0 -> m1 -> m2 -> m3
handler := chain.Then(finalHandler)
```

### Concurrency Patterns
- **MapReduce** (`core/mr`): Parallel processing with worker pools - use for batch operations
- **Executors** (`core/executors`): Bulk/period executors for batching operations
- **SingleFlight** (`core/syncx`): Deduplicates concurrent identical requests

Remember to run tests and ensure all checks pass before submitting changes. The project emphasizes high quality, performance, and reliability, so these should be primary considerations in all development work.