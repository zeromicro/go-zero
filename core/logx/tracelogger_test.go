package logx

import (
	"context"
	"log"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/trace/tracespec"
)

const (
	mockTraceId = "mock-trace-id"
	mockSpanId  = "mock-span-id"
)

var mock tracespec.Trace = new(mockTrace)

func TestTraceLog(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	ctx := context.WithValue(context.Background(), tracespec.TracingKey, mock)
	WithContext(ctx).(*traceLogger).write(&buf, levelInfo, testlog)
	assert.True(t, strings.Contains(buf.String(), mockTraceId))
	assert.True(t, strings.Contains(buf.String(), mockSpanId))
}

func TestTraceError(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	errorLog = newLogWriter(log.New(&buf, "", flags))
	ctx := context.WithValue(context.Background(), tracespec.TracingKey, mock)
	l := WithContext(ctx).(*traceLogger)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Error(testlog)
	assert.True(t, strings.Contains(buf.String(), mockTraceId))
	assert.True(t, strings.Contains(buf.String(), mockSpanId))
	buf.Reset()
	l.WithDuration(time.Second).Errorf(testlog)
	assert.True(t, strings.Contains(buf.String(), mockTraceId))
	assert.True(t, strings.Contains(buf.String(), mockSpanId))
}

func TestTraceInfo(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	infoLog = newLogWriter(log.New(&buf, "", flags))
	ctx := context.WithValue(context.Background(), tracespec.TracingKey, mock)
	l := WithContext(ctx).(*traceLogger)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	assert.True(t, strings.Contains(buf.String(), mockTraceId))
	assert.True(t, strings.Contains(buf.String(), mockSpanId))
	buf.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	assert.True(t, strings.Contains(buf.String(), mockTraceId))
	assert.True(t, strings.Contains(buf.String(), mockSpanId))
}

func TestTraceSlow(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	slowLog = newLogWriter(log.New(&buf, "", flags))
	ctx := context.WithValue(context.Background(), tracespec.TracingKey, mock)
	l := WithContext(ctx).(*traceLogger)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Slow(testlog)
	assert.True(t, strings.Contains(buf.String(), mockTraceId))
	assert.True(t, strings.Contains(buf.String(), mockSpanId))
	buf.Reset()
	l.WithDuration(time.Second).Slowf(testlog)
	assert.True(t, strings.Contains(buf.String(), mockTraceId))
	assert.True(t, strings.Contains(buf.String(), mockSpanId))
}

func TestTraceWithoutContext(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	infoLog = newLogWriter(log.New(&buf, "", flags))
	l := WithContext(context.Background()).(*traceLogger)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	assert.False(t, strings.Contains(buf.String(), mockTraceId))
	assert.False(t, strings.Contains(buf.String(), mockSpanId))
	buf.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	assert.False(t, strings.Contains(buf.String(), mockTraceId))
	assert.False(t, strings.Contains(buf.String(), mockSpanId))
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
