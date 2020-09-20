package tracespec

// contextKey a type for context key
type contextKey string

// TracingKey is tracing key for context
var TracingKey = contextKey("X-Trace")
