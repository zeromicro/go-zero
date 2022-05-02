package logx

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestWithDurationError(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Error("foo")
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
}

func TestWithDurationErrorf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Errorf("foo")
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
}

func TestWithDurationErrorv(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Errorv("foo")
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
}

func TestWithDurationErrorw(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Errorw("foo", Field("foo", "bar"))
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
	assert.True(t, strings.Contains(w.String(), "foo"), w.String())
	assert.True(t, strings.Contains(w.String(), "bar"), w.String())
}

func TestWithDurationInfo(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Info("foo")
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
}

func TestWithDurationInfoConsole(t *testing.T) {
	old := atomic.LoadUint32(&encoding)
	atomic.StoreUint32(&encoding, plainEncodingType)
	defer func() {
		atomic.StoreUint32(&encoding, old)
	}()

	w := new(mockWriter)
	o := writer.Swap(w)
	defer writer.Store(o)

	WithDuration(time.Second).Info("foo")
	assert.True(t, strings.Contains(w.String(), "ms"), w.String())
}

func TestWithDurationInfof(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Infof("foo")
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
}

func TestWithDurationInfov(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Infov("foo")
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
}

func TestWithDurationInfow(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Infow("foo", Field("foo", "bar"))
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
	assert.True(t, strings.Contains(w.String(), "foo"), w.String())
	assert.True(t, strings.Contains(w.String(), "bar"), w.String())
}

func TestWithDurationWithContextInfow(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	WithDuration(time.Second).WithContext(ctx).Infow("foo", Field("foo", "bar"))
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
	assert.True(t, strings.Contains(w.String(), "foo"), w.String())
	assert.True(t, strings.Contains(w.String(), "bar"), w.String())
	assert.True(t, strings.Contains(w.String(), "trace"), w.String())
	assert.True(t, strings.Contains(w.String(), "span"), w.String())
}

func TestWithDurationSlow(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).Slow("foo")
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
}

func TestWithDurationSlowf(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).WithDuration(time.Hour).Slowf("foo")
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
}

func TestWithDurationSlowv(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).WithDuration(time.Hour).Slowv("foo")
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
}

func TestWithDurationSloww(t *testing.T) {
	w := new(mockWriter)
	old := writer.Swap(w)
	defer writer.Store(old)

	WithDuration(time.Second).WithDuration(time.Hour).Sloww("foo", Field("foo", "bar"))
	assert.True(t, strings.Contains(w.String(), "duration"), w.String())
	assert.True(t, strings.Contains(w.String(), "foo"), w.String())
	assert.True(t, strings.Contains(w.String(), "bar"), w.String())
}
