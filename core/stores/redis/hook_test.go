package redis

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	red "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace/tracetest"
	tracesdk "go.opentelemetry.io/otel/trace"
)

func TestHookProcessCase1(t *testing.T) {
	tracetest.NewInMemoryExporter(t)

	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	err := durationHook.ProcessHook(func(ctx context.Context, cmd red.Cmder) error {
		assert.Equal(t, "redis", tracesdk.SpanFromContext(ctx).(interface{ Name() string }).Name())
		return nil
	})(context.Background(), red.NewCmd(context.Background()))
	if err != nil {
		t.Fatal(err)
	}

	assert.False(t, strings.Contains(buf.String(), "slow"))
}

func TestHookProcessCase2(t *testing.T) {
	tracetest.NewInMemoryExporter(t)

	w, restore := injectLog()
	defer restore()

	err := durationHook.ProcessHook(func(ctx context.Context, cmd red.Cmder) error {
		assert.Equal(t, "redis", tracesdk.SpanFromContext(ctx).(interface{ Name() string }).Name())
		time.Sleep(slowThreshold.Load() + time.Millisecond)
		return nil
	})(context.Background(), red.NewCmd(context.Background()))
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, strings.Contains(w.String(), "slow"))
	assert.True(t, strings.Contains(w.String(), "trace"))
	assert.True(t, strings.Contains(w.String(), "span"))
}

func TestHookProcessPipelineCase1(t *testing.T) {
	tracetest.NewInMemoryExporter(t)
	writer := log.Writer()
	var buf strings.Builder
	log.SetOutput(&buf)
	defer log.SetOutput(writer)

	err := durationHook.ProcessPipelineHook(func(ctx context.Context, cmds []red.Cmder) error {
		return nil
	})(context.Background(), nil)
	assert.NoError(t, err)

	err = durationHook.ProcessPipelineHook(func(ctx context.Context, cmds []red.Cmder) error {
		assert.Equal(t, "redis", tracesdk.SpanFromContext(ctx).(interface{ Name() string }).Name())
		return nil
	})(context.Background(), []red.Cmder{
		red.NewCmd(context.Background()),
	})
	assert.NoError(t, err)

	assert.False(t, strings.Contains(buf.String(), "slow"))
}

func TestHookProcessPipelineCase2(t *testing.T) {
	tracetest.NewInMemoryExporter(t)

	w, restore := injectLog()
	defer restore()

	err := durationHook.ProcessPipelineHook(func(ctx context.Context, cmds []red.Cmder) error {
		assert.Equal(t, "redis", tracesdk.SpanFromContext(ctx).(interface{ Name() string }).Name())
		time.Sleep(slowThreshold.Load() + time.Millisecond)
		return nil
	})(context.Background(), []red.Cmder{
		red.NewCmd(context.Background()),
	})
	assert.NoError(t, err)

	assert.True(t, strings.Contains(w.String(), "slow"))
	assert.True(t, strings.Contains(w.String(), "trace"))
	assert.True(t, strings.Contains(w.String(), "span"))
}

func TestLogDuration(t *testing.T) {
	w, restore := injectLog()
	defer restore()

	logDuration(context.Background(), []red.Cmder{
		red.NewCmd(context.Background(), "get", "foo"),
	}, 1*time.Second)
	assert.True(t, strings.Contains(w.String(), "get foo"))

	logDuration(context.Background(), []red.Cmder{
		red.NewCmd(context.Background(), "get", "foo"),
		red.NewCmd(context.Background(), "set", "bar", 0),
	}, 1*time.Second)
	assert.True(t, strings.Contains(w.String(), `get foo\nset bar 0`))
}

func injectLog() (r *strings.Builder, restore func()) {
	var buf strings.Builder
	w := logx.NewWriter(&buf)
	o := logx.Reset()
	logx.SetWriter(w)

	return &buf, func() {
		logx.Reset()
		logx.SetWriter(o)
	}
}
