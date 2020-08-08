package logx

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/trace/tracespec"
)

const (
	mockTraceId = "mock-trace-id"
	mockSpanId  = "mock-span-id"
)

var mock tracespec.Trace = new(mockTrace)

func TestTraceLog(t *testing.T) {
	var buf strings.Builder
	ctx := context.WithValue(context.Background(), tracespec.TracingKey, mock)
	WithContext(ctx).(tracingEntry).write(&buf, levelInfo, testlog)
	assert.True(t, strings.Contains(buf.String(), mockTraceId))
	assert.True(t, strings.Contains(buf.String(), mockSpanId))
}

type mockTrace struct{}

func (t mockTrace) TraceId() string {
	return mockTraceId
}

func (t mockTrace) SpanId() string {
	return mockSpanId
}

func (t mockTrace) Finish() {
}

func (t mockTrace) Fork(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	return nil, nil
}

func (t mockTrace) Follow(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	return nil, nil
}

func (t mockTrace) Visit(fn func(key string, val string) bool) {
}
