package logx

import (
	"context"
	"log"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	traceKey = "trace"
	spanKey  = "span"
)

func TestTraceLog(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	WithContext(ctx).(*traceLogger).write(&buf, levelInfo, testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
}

func TestTraceError(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	errorLog = newLogWriter(log.New(&buf, "", flags))
	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	l := WithContext(ctx).(*traceLogger)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Error(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
	buf.Reset()
	l.WithDuration(time.Second).Errorf(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
	buf.Reset()
	l.WithDuration(time.Second).Errorv(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
}

func TestTraceInfo(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	infoLog = newLogWriter(log.New(&buf, "", flags))
	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	l := WithContext(ctx).(*traceLogger)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
	buf.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
	buf.Reset()
	l.WithDuration(time.Second).Infov(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
}

func TestTraceSlow(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	slowLog = newLogWriter(log.New(&buf, "", flags))
	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	l := WithContext(ctx).(*traceLogger)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Slow(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
	buf.Reset()
	l.WithDuration(time.Second).Slowf(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
	buf.Reset()
	l.WithDuration(time.Second).Slowv(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
}

func TestTraceWithoutContext(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	infoLog = newLogWriter(log.New(&buf, "", flags))
	l := WithContext(context.Background()).(*traceLogger)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	assert.False(t, strings.Contains(buf.String(), traceKey))
	assert.False(t, strings.Contains(buf.String(), spanKey))
	buf.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	assert.False(t, strings.Contains(buf.String(), traceKey))
	assert.False(t, strings.Contains(buf.String(), spanKey))
}
