package fx

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

var (
	// ErrCanceled is the error returned when the context is canceled.
	ErrCanceled = context.Canceled
	// ErrTimeout is the error returned when the context's deadline passes.
	ErrTimeout = context.DeadlineExceeded
)

// DoOption defines the method to customize a DoWithTimeout call.
type DoOption func() context.Context

// DoWithTimeout runs fn with timeout control.
func DoWithTimeout(fn func() error, timeout time.Duration, opts ...DoOption) error {
	parentCtx := context.Background()
	for _, opt := range opts {
		parentCtx = opt()
	}
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	// create channel with buffer size 1 to avoid goroutine leak
	done := make(chan error, 1)
	panicChan := make(chan any, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				// attach call stack to avoid missing in different goroutine
				panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
			}
		}()
		done <- fn()
	}()

	select {
	case p := <-panicChan:
		panic(p)
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// WithContext customizes a DoWithTimeout call with given ctx.
func WithContext(ctx context.Context) DoOption {
	return func() context.Context {
		return ctx
	}
}
