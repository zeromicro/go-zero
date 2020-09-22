package trace

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/core/trace/tracespec"
	"google.golang.org/grpc/metadata"
)

func TestClientSpan(t *testing.T) {
	span := newServerSpan(nil, "service", "operation")
	ctx := context.WithValue(context.Background(), tracespec.TracingKey, span)
	ctx, span = StartClientSpan(ctx, "entrance", "operation")
	defer span.Finish()
	assert.Equal(t, span, ctx.Value(tracespec.TracingKey))

	const serviceName = "authorization"
	const operationName = "verification"
	ctx, childSpan := span.Fork(ctx, serviceName, operationName)
	defer childSpan.Finish()

	assert.Equal(t, childSpan, ctx.Value(tracespec.TracingKey))
	assert.Equal(t, getSpan(span).TraceId(), getSpan(childSpan).TraceId())
	assert.Equal(t, "0.1.1", getSpan(childSpan).SpanId())
	assert.Equal(t, serviceName, childSpan.(*Span).serviceName)
	assert.Equal(t, operationName, childSpan.(*Span).operationName)
	assert.Equal(t, clientFlag, childSpan.(*Span).flag)
}

func TestClientSpan_WithoutTrace(t *testing.T) {
	ctx, span := StartClientSpan(context.Background(), "entrance", "operation")
	defer span.Finish()
	assert.Equal(t, emptyNoopSpan, span)
	assert.Equal(t, context.Background(), ctx)
}

func TestServerSpan(t *testing.T) {
	ctx, span := StartServerSpan(context.Background(), nil, "service", "operation")
	defer span.Finish()
	assert.Equal(t, span, ctx.Value(tracespec.TracingKey))

	const serviceName = "authorization"
	const operationName = "verification"
	ctx, childSpan := span.Fork(ctx, serviceName, operationName)
	defer childSpan.Finish()

	assert.Equal(t, childSpan, ctx.Value(tracespec.TracingKey))
	assert.Equal(t, getSpan(span).TraceId(), getSpan(childSpan).TraceId())
	assert.Equal(t, "0.1", getSpan(childSpan).SpanId())
	assert.Equal(t, serviceName, childSpan.(*Span).serviceName)
	assert.Equal(t, operationName, childSpan.(*Span).operationName)
	assert.Equal(t, clientFlag, childSpan.(*Span).flag)
}

func TestServerSpan_WithCarrier(t *testing.T) {
	md := metadata.New(map[string]string{
		traceIdKey: "a",
		spanIdKey:  "0.1",
	})
	ctx, span := StartServerSpan(context.Background(), grpcCarrier(md), "service", "operation")
	defer span.Finish()
	assert.Equal(t, span, ctx.Value(tracespec.TracingKey))

	const serviceName = "authorization"
	const operationName = "verification"
	ctx, childSpan := span.Fork(ctx, serviceName, operationName)
	defer childSpan.Finish()

	assert.Equal(t, childSpan, ctx.Value(tracespec.TracingKey))
	assert.Equal(t, getSpan(span).TraceId(), getSpan(childSpan).TraceId())
	assert.Equal(t, "0.1.1", getSpan(childSpan).SpanId())
	assert.Equal(t, serviceName, childSpan.(*Span).serviceName)
	assert.Equal(t, operationName, childSpan.(*Span).operationName)
	assert.Equal(t, clientFlag, childSpan.(*Span).flag)
}

func TestSpan_Follow(t *testing.T) {
	tests := []struct {
		span       string
		expectSpan string
	}{
		{
			"0.1",
			"0.2",
		},
		{
			"0",
			"1",
		},
		{
			"a",
			"a",
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			md := metadata.New(map[string]string{
				traceIdKey: "a",
				spanIdKey:  test.span,
			})
			ctx, span := StartServerSpan(context.Background(), grpcCarrier(md),
				"service", "operation")
			defer span.Finish()
			assert.Equal(t, span, ctx.Value(tracespec.TracingKey))

			const serviceName = "authorization"
			const operationName = "verification"
			ctx, childSpan := span.Follow(ctx, serviceName, operationName)
			defer childSpan.Finish()

			assert.Equal(t, childSpan, ctx.Value(tracespec.TracingKey))
			assert.Equal(t, getSpan(span).TraceId(), getSpan(childSpan).TraceId())
			assert.Equal(t, test.expectSpan, getSpan(childSpan).SpanId())
			assert.Equal(t, serviceName, childSpan.(*Span).serviceName)
			assert.Equal(t, operationName, childSpan.(*Span).operationName)
			assert.Equal(t, span.(*Span).flag, childSpan.(*Span).flag)
		})
	}
}

func TestSpan_Visit(t *testing.T) {
	var run bool
	span := newServerSpan(nil, "service", "operation")
	span.Visit(func(key, val string) bool {
		assert.True(t, len(key) > 0)
		assert.True(t, len(val) > 0)
		run = true
		return true
	})
	assert.True(t, run)
}

func getSpan(span tracespec.Trace) tracespec.Trace {
	return span.(*Span)
}
