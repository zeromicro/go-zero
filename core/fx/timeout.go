package fx

import (
	"context"
	"time"
)

var (
	ErrCanceled = context.Canceled
	ErrTimeout  = context.DeadlineExceeded
)

type FxOption func() context.Context

func DoWithTimeout(fn func() error, timeout time.Duration, opts ...FxOption) error {
	parentCtx := context.Background()
	for _, opt := range opts {
		parentCtx = opt()
	}
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	done := make(chan error)
	panicChan := make(chan interface{}, 1)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()
		done <- fn()
		close(done)
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

func WithContext(ctx context.Context) FxOption {
	return func() context.Context {
		return ctx
	}
}
