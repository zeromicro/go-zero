package logx

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestTraceLog(t *testing.T) {
	SetLevel(InfoLevel)
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	WithContext(ctx).Info(testlog)
	validate(t, w.String(), true, true)
}

func TestTraceError(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)
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
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Errorf(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Errorv(testlog)
	fmt.Println(w.String())
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Errorw(testlog, Field("foo", "bar"))
	validate(t, w.String(), true, true)
	assert.True(t, strings.Contains(w.String(), "foo"), w.String())
	assert.True(t, strings.Contains(w.String(), "bar"), w.String())
}

func TestTraceInfo(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	l := WithContext(ctx)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Infov(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Infow(testlog, Field("foo", "bar"))
	validate(t, w.String(), true, true)
	assert.True(t, strings.Contains(w.String(), "foo"), w.String())
	assert.True(t, strings.Contains(w.String(), "bar"), w.String())
}

func TestTraceInfoConsole(t *testing.T) {
	old := atomic.LoadUint32(&encoding)
	atomic.StoreUint32(&encoding, jsonEncodingType)
	defer func() {
		atomic.StoreUint32(&encoding, old)
	}()

	w := new(mockWriter)
	o := writer.Swap(w)
	defer writer.Store(o)

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	l := WithContext(ctx)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Infov(testlog)
	validate(t, w.String(), true, true)
}

func TestTraceSlow(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	l := WithContext(ctx)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Slow(testlog)
	assert.True(t, strings.Contains(w.String(), traceKey))
	assert.True(t, strings.Contains(w.String(), spanKey))
	w.Reset()
	l.WithDuration(time.Second).Slowf(testlog)
	fmt.Println("buf:", w.String())
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Slowv(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Sloww(testlog, Field("foo", "bar"))
	validate(t, w.String(), true, true)
	assert.True(t, strings.Contains(w.String(), "foo"), w.String())
	assert.True(t, strings.Contains(w.String(), "bar"), w.String())
}

func TestTraceWithoutContext(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	l := WithContext(context.Background())
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	validate(t, w.String(), false, false)
	w.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	validate(t, w.String(), false, false)
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
