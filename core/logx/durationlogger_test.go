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

func TestWithDurationError(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Error("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationErrorf(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Errorf("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationErrorv(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Errorv("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationErrorw(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Errorw("foo", Field("foo", "bar"))
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "foo"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "bar"), builder.String())
}

func TestWithDurationInfo(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Info("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationInfoConsole(t *testing.T) {
	old := atomic.LoadUint32(&encoding)
	atomic.StoreUint32(&encoding, plainEncodingType)
	defer func() {
		atomic.StoreUint32(&encoding, old)
	}()

	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Info("foo")
	assert.True(t, strings.Contains(builder.String(), "ms"), builder.String())
}

func TestWithDurationInfof(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Infof("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationInfov(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Infov("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationInfow(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Infow("foo", Field("foo", "bar"))
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "foo"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "bar"), builder.String())
}

func TestWithDurationWithContextInfow(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)

	otp := otel.GetTracerProvider()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.AlwaysSample()))
	otel.SetTracerProvider(tp)
	defer otel.SetTracerProvider(otp)

	ctx, _ := tp.Tracer("foo").Start(context.Background(), "bar")
	WithDuration(time.Second).WithContext(ctx).Infow("foo", Field("foo", "bar"))
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "foo"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "bar"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "trace"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "span"), builder.String())
}

func TestWithDurationSlow(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).Slow("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationSlowf(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).WithDuration(time.Hour).Slowf("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationSlowv(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).WithDuration(time.Hour).Slowv("foo")
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
}

func TestWithDurationSloww(t *testing.T) {
	var builder strings.Builder
	log.SetOutput(&builder)
	WithDuration(time.Second).WithDuration(time.Hour).Sloww("foo", Field("foo", "bar"))
	assert.True(t, strings.Contains(builder.String(), "duration"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "foo"), builder.String())
	assert.True(t, strings.Contains(builder.String(), "bar"), builder.String())
}
