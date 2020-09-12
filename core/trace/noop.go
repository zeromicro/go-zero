package trace

import (
	"context"

	"github.com/tal-tech/go-zero/core/trace/tracespec"
)

var emptyNoopSpan = noopSpan{}

type noopSpan struct{}

func (s noopSpan) Finish() {
}

func (s noopSpan) Follow(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	return ctx, emptyNoopSpan
}

func (s noopSpan) Fork(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	return ctx, emptyNoopSpan
}

func (s noopSpan) SpanId() string {
	return ""
}

func (s noopSpan) TraceId() string {
	return ""
}

func (s noopSpan) Visit(fn func(key, val string) bool) {
}
