package logx

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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
	validate(t, buf.String(), true, true)
}

func TestTraceError(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	errorLog = newLogWriter(log.New(&buf, "", flags))
	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	l := WithContext(context.Background())
	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	var nilCtx context.Context
	l = l.WithContext(nilCtx)
	l = l.WithContext(ctx)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Error(testlog)
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Errorf(testlog)
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Errorv(testlog)
	fmt.Println(buf.String())
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Errorw(testlog, Field("foo", "bar"))
	validate(t, buf.String(), true, true)
	assert.True(t, strings.Contains(buf.String(), "foo"), buf.String())
	assert.True(t, strings.Contains(buf.String(), "bar"), buf.String())
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
	l := WithContext(ctx)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Infov(testlog)
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Infow(testlog, Field("foo", "bar"))
	validate(t, buf.String(), true, true)
	assert.True(t, strings.Contains(buf.String(), "foo"), buf.String())
	assert.True(t, strings.Contains(buf.String(), "bar"), buf.String())
}

func TestTraceInfoConsole(t *testing.T) {
	old := atomic.LoadUint32(&encoding)
	atomic.StoreUint32(&encoding, jsonEncodingType)
	defer func() {
		atomic.StoreUint32(&encoding, old)
	}()

	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	infoLog = newLogWriter(log.New(&buf, "", flags))
	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	l := WithContext(ctx)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Infov(testlog)
	validate(t, buf.String(), true, true)
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
	l := WithContext(ctx)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Slow(testlog)
	assert.True(t, strings.Contains(buf.String(), traceKey))
	assert.True(t, strings.Contains(buf.String(), spanKey))
	buf.Reset()
	l.WithDuration(time.Second).Slowf(testlog)
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Slowv(testlog)
	validate(t, buf.String(), true, true)
	buf.Reset()
	l.WithDuration(time.Second).Sloww(testlog, Field("foo", "bar"))
	validate(t, buf.String(), true, true)
	assert.True(t, strings.Contains(buf.String(), "foo"), buf.String())
	assert.True(t, strings.Contains(buf.String(), "bar"), buf.String())
}

func TestTraceWithoutContext(t *testing.T) {
	var buf mockWriter
	atomic.StoreUint32(&initialized, 1)
	infoLog = newLogWriter(log.New(&buf, "", flags))
	l := WithContext(context.Background())
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	validate(t, buf.String(), false, false)
	buf.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	validate(t, buf.String(), false, false)
}

func validate(t *testing.T, body string, expectedTrace, expectedSpan bool) {
	var val mockValue
	assert.Nil(t, json.Unmarshal([]byte(body), &val))
	assert.Equal(t, expectedTrace, len(val.Trace) > 0)
	assert.Equal(t, expectedSpan, len(val.Span) > 0)
}

type mockValue struct {
	Trace string `json:"trace"`
	Span  string `json:"span"`
}
