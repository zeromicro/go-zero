package tracespec

import "context"

// Trace interface represents a tracing.
type Trace interface {
	SpanContext
	Finish()
	Fork(ctx context.Context, serviceName, operationName string) (context.Context, Trace)
	Follow(ctx context.Context, serviceName, operationName string) (context.Context, Trace)
}
