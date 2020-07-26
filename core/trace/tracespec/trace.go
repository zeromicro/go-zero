package tracespec

import "context"

type Trace interface {
	SpanContext
	Finish()
	Fork(ctx context.Context, serviceName, operationName string) (context.Context, Trace)
	Follow(ctx context.Context, serviceName, operationName string) (context.Context, Trace)
}
