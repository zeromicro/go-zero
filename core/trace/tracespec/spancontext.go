package tracespec

// SpanContext interface that represents a span context.
type SpanContext interface {
	TraceId() string
	SpanId() string
	Visit(fn func(key, val string) bool)
}
