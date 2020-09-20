package tracespec

// TracingKey is tracing key for context
var TracingKey = contextKey("X-Trace")

// contextKey a type for context key
type contextKey string

// Printing a context will reveal a fair amount of information about it.
func (c contextKey) String() string {
	return "trace/tracespec context key " + string(c)
}
