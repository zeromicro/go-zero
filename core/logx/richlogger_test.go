package logx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, span := tp.Tracer("trace-id").Start(context.Background(), "span-id")
	defer span.End()

	WithContext(ctx).Info(testlog)
	validate(t, w.String(), true, true)
}

func TestTraceDebug(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, span := tp.Tracer("foo").Start(context.Background(), "bar")
	defer span.End()

	l := WithContext(ctx)
	SetLevel(DebugLevel)
	l.WithDuration(time.Second).Debug(testlog)
	assert.True(t, strings.Contains(w.String(), traceKey))
	assert.True(t, strings.Contains(w.String(), spanKey))
	w.Reset()
	l.WithDuration(time.Second).Debugf(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Debugfn(func() any {
		return testlog
	})
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Debugv(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Debugv(testobj)
	validateContentType(t, w.String(), map[string]any{}, true, true)
	w.Reset()
	l.WithDuration(time.Second).Debugw(testlog, Field("foo", "bar"))
	validate(t, w.String(), true, true)
	assert.True(t, strings.Contains(w.String(), "foo"), w.String())
	assert.True(t, strings.Contains(w.String(), "bar"), w.String())
}

func TestTraceError(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, span := tp.Tracer("trace-id").Start(context.Background(), "span-id")
	defer span.End()

	var nilCtx context.Context
	l := WithContext(context.Background())
	l = l.WithContext(nilCtx)
	l = l.WithContext(ctx)
	SetLevel(ErrorLevel)
	l.WithDuration(time.Second).Error(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Errorf(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Errorfn(func() any {
		return testlog
	})
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Errorv(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Errorv(testobj)
	validateContentType(t, w.String(), map[string]any{}, true, true)
	w.Reset()
	l.WithDuration(time.Second).Errorw(testlog, Field("basket", "ball"))
	validate(t, w.String(), true, true)
	assert.True(t, strings.Contains(w.String(), "basket"), w.String())
	assert.True(t, strings.Contains(w.String(), "ball"), w.String())
}

func TestTraceInfo(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, span := tp.Tracer("trace-id").Start(context.Background(), "span-id")
	defer span.End()

	SetLevel(InfoLevel)
	l := WithContext(ctx)
	l.WithDuration(time.Second).Info(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Infofn(func() any {
		return testlog
	})
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Infov(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Infov(testobj)
	validateContentType(t, w.String(), map[string]any{}, true, true)
	w.Reset()
	l.WithDuration(time.Second).Infow(testlog, Field("basket", "ball"))
	validate(t, w.String(), true, true)
	assert.True(t, strings.Contains(w.String(), "basket"), w.String())
	assert.True(t, strings.Contains(w.String(), "ball"), w.String())
}

func TestTraceInfoConsole(t *testing.T) {
	old := atomic.SwapUint32(&encoding, jsonEncodingType)
	defer atomic.StoreUint32(&encoding, old)

	w := new(mockWriter)
	o := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(o)
	}()

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, span := tp.Tracer("trace-id").Start(context.Background(), "span-id")
	defer span.End()

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
	l.WithDuration(time.Second).Infov(testobj)
	validateContentType(t, w.String(), map[string]any{}, true, true)
}

func TestTraceSlow(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, span := tp.Tracer("trace-id").Start(context.Background(), "span-id")
	defer span.End()

	l := WithContext(ctx)
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Slow(testlog)
	assert.True(t, strings.Contains(w.String(), traceKey))
	assert.True(t, strings.Contains(w.String(), spanKey))
	w.Reset()
	l.WithDuration(time.Second).Slowf(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Slowfn(func() any {
		return testlog
	})
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Slowv(testlog)
	validate(t, w.String(), true, true)
	w.Reset()
	l.WithDuration(time.Second).Slowv(testobj)
	validateContentType(t, w.String(), map[string]any{}, true, true)
	w.Reset()
	l.WithDuration(time.Second).Sloww(testlog, Field("basket", "ball"))
	validate(t, w.String(), true, true)
	assert.True(t, strings.Contains(w.String(), "basket"), w.String())
	assert.True(t, strings.Contains(w.String(), "ball"), w.String())
}

func TestTraceWithoutContext(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	l := WithContext(context.Background())
	SetLevel(InfoLevel)
	l.WithDuration(time.Second).Info(testlog)
	validate(t, w.String(), false, false)
	w.Reset()
	l.WithDuration(time.Second).Infof(testlog)
	validate(t, w.String(), false, false)
}

func TestLogWithFields(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	ctx := ContextWithFields(context.Background(), Field("foo", "bar"))
	l := WithContext(ctx)
	SetLevel(InfoLevel)
	l.Info(testlog)

	var val mockValue
	assert.Nil(t, json.Unmarshal([]byte(w.String()), &val))
	assert.Equal(t, "bar", val.Foo)
}

func TestLogWithCallerSkip(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	l := WithCallerSkip(1).WithCallerSkip(0)
	p := func(v string) {
		l.Info(v)
	}

	file, line := getFileLine()
	p(testlog)
	assert.True(t, w.Contains(fmt.Sprintf("%s:%d", file, line+1)))

	w.Reset()
	l = WithCallerSkip(0).WithCallerSkip(1)
	file, line = getFileLine()
	p(testlog)
	assert.True(t, w.Contains(fmt.Sprintf("%s:%d", file, line+1)))
}

func TestLogWithCallerSkipCopy(t *testing.T) {
	log1 := WithCallerSkip(2)
	log2 := log1.WithCallerSkip(3)
	log3 := log2.WithCallerSkip(-1)
	assert.Equal(t, 2, log1.(*richLogger).callerSkip)
	assert.Equal(t, 3, log2.(*richLogger).callerSkip)
	assert.Equal(t, 3, log3.(*richLogger).callerSkip)
}

func TestLogWithContextCopy(t *testing.T) {
	c1 := context.Background()
	c2 := context.WithValue(context.Background(), "foo", "bar")
	log1 := WithContext(c1)
	log2 := log1.WithContext(c2)
	assert.Equal(t, c1, log1.(*richLogger).ctx)
	assert.Equal(t, c2, log2.(*richLogger).ctx)
}

func TestLogWithDurationCopy(t *testing.T) {
	log1 := WithContext(context.Background())
	log2 := log1.WithDuration(time.Second)
	assert.Empty(t, log1.(*richLogger).fields)
	assert.Equal(t, 1, len(log2.(*richLogger).fields))

	var w mockWriter
	old := writer.Swap(&w)
	defer writer.Store(old)
	log2.Info("hello")
	assert.Contains(t, w.String(), `"duration":"1000.0ms"`)
}

func TestLogWithFieldsCopy(t *testing.T) {
	log1 := WithContext(context.Background())
	log2 := log1.WithFields(Field("foo", "bar"))
	log3 := log1.WithFields()
	assert.Empty(t, log1.(*richLogger).fields)
	assert.Equal(t, 1, len(log2.(*richLogger).fields))
	assert.Equal(t, log1, log3)
	assert.Empty(t, log3.(*richLogger).fields)

	var w mockWriter
	old := writer.Swap(&w)
	defer writer.Store(old)

	log2.Info("hello")
	assert.Contains(t, w.String(), `"foo":"bar"`)
}

func TestLoggerWithFields(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	writer.lock.RLock()
	defer func() {
		writer.lock.RUnlock()
		writer.Store(old)
	}()

	l := WithContext(context.Background()).WithFields(Field("foo", "bar"))
	l.Info(testlog)

	var val mockValue
	assert.Nil(t, json.Unmarshal([]byte(w.String()), &val))
	assert.Equal(t, "bar", val.Foo)
}

func validate(t *testing.T, body string, expectedTrace, expectedSpan bool) {
	var val mockValue
	dec := json.NewDecoder(strings.NewReader(body))

	for {
		var doc mockValue
		err := dec.Decode(&doc)
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			continue
		}

		val = doc
	}

	assert.Equal(t, expectedTrace, len(val.Trace) > 0, body)
	assert.Equal(t, expectedSpan, len(val.Span) > 0, body)
}

func validateContentType(t *testing.T, body string, expectedType any, expectedTrace, expectedSpan bool) {
	var val mockValue
	dec := json.NewDecoder(strings.NewReader(body))

	for {
		var doc mockValue
		err := dec.Decode(&doc)
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			continue
		}

		val = doc
	}

	assert.IsType(t, expectedType, val.Content, body)
	assert.Equal(t, expectedTrace, len(val.Trace) > 0, body)
	assert.Equal(t, expectedSpan, len(val.Span) > 0, body)
}

type mockValue struct {
	Trace   string `json:"trace"`
	Span    string `json:"span"`
	Foo     string `json:"foo"`
	Content any    `json:"content"`
}
