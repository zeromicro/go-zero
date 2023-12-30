package redis

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	red "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx/logtest"
	ztrace "github.com/zeromicro/go-zero/core/trace"
)

func TestHookProcessCase1(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})
	defer ztrace.StopAgent()

	w := logtest.NewCollector(t)
	hookFunc := durationHook.ProcessHook(func(ctx context.Context, cmd red.Cmder) error {
		return nil
	})

	err := hookFunc(context.Background(), red.NewCmd(context.Background()))
	assert.Nil(t, err)
	str := w.String()
	assert.False(t, strings.Contains(str, "slow"))
}

func TestHookProcessCase2(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})
	defer ztrace.StopAgent()

	w := logtest.NewCollector(t)
	hookFunc := durationHook.ProcessHook(func(ctx context.Context, cmd red.Cmder) error {
		time.Sleep(slowThreshold.Load() + time.Millisecond)
		return nil
	})

	err := hookFunc(context.Background(), red.NewCmd(context.Background(), "foo", "bar"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.True(t, strings.Contains(w.String(), "slow"))
}

func TestHookProcessPipelineCase1(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})
	defer ztrace.StopAgent()

	w := logtest.NewCollector(t)
	hookFunc := durationHook.ProcessPipelineHook(func(ctx context.Context, cmds []red.Cmder) error {
		return nil
	})

	err := hookFunc(context.Background(), []red.Cmder{
		red.NewCmd(context.Background()),
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.False(t, strings.Contains(w.String(), "slow"))
}

func TestHookProcessPipelineCase2(t *testing.T) {
	ztrace.StartAgent(ztrace.Config{
		Name:     "go-zero-test",
		Endpoint: "http://localhost:14268/api/traces",
		Batcher:  "jaeger",
		Sampler:  1.0,
	})
	defer ztrace.StopAgent()

	w := logtest.NewCollector(t)
	hookFunc := durationHook.ProcessPipelineHook(func(ctx context.Context, cmds []red.Cmder) error {
		time.Sleep(slowThreshold.Load() + time.Millisecond)
		return nil
	})

	err := hookFunc(context.Background(), []red.Cmder{
		red.NewCmd(context.Background()),
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.True(t, strings.Contains(w.String(), "slow"))
}

func TestLogDuration(t *testing.T) {
	w := logtest.NewCollector(t)

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

func TestFormatError(t *testing.T) {
	// Test case: err is OpError
	err := &net.OpError{
		Err: mockOpError{},
	}
	assert.Equal(t, "timeout", formatError(err))

	// Test case: err is nil
	assert.Equal(t, "", formatError(nil))

	// Test case: err is red.Nil
	assert.Equal(t, "", formatError(red.Nil))

	// Test case: err is io.EOF
	assert.Equal(t, "eof", formatError(io.EOF))

	// Test case: err is context.DeadlineExceeded
	assert.Equal(t, "context deadline", formatError(context.DeadlineExceeded))

	// Test case: err is breaker.ErrServiceUnavailable
	assert.Equal(t, "breaker", formatError(breaker.ErrServiceUnavailable))

	// Test case: err is unknown
	assert.Equal(t, "unexpected error", formatError(errors.New("some error")))
}

type mockOpError struct {
}

func (mockOpError) Error() string {
	return "mock error"
}

func (mockOpError) Timeout() bool {
	return true
}
