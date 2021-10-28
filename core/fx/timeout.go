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
	// Deprecated
	ErrCanceled = context.Canceled
	// ErrTimeout is the error returned when the context's deadline passes.
	// Deprecated
	ErrTimeout = context.DeadlineExceeded
)

// DoOption defines the method to customize a DoWithTimeout call.
// Deprecated
type DoOption func() context.Context

// DoWithTimeout runs fn with timeout control.
// Deprecated: Use ExecuteWithTimeout instead
func DoWithTimeout(fn func() error, timeout time.Duration, opts ...DoOption) error {
	parentCtx := context.Background()
	for _, opt := range opts {
		parentCtx = opt()
	}
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	// create channel with buffer size 1 to avoid goroutine leak
	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
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
// Deprecated
func WithContext(ctx context.Context) DoOption {
	return func() context.Context {
		return ctx
	}
}

// ExecuteWithTimeout runs fn with timeout control.
func ExecuteWithTimeout(ctx context.Context, fn func() error) (err error) {
	// create channel with buffer size 1 to avoid goroutine leak
	done := make(chan error, 1)
	panicChan := make(chan interface{}, 1)
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
	case err = <-done:
		return
	case <-ctx.Done():
		return ctx.Err()
	}
}
