package tracespec

type SpanContext interface {
	TraceId() string
	SpanId() string
	Visit(fn func(key, val string) bool)
}
