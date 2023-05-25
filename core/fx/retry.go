package fx

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/errorx"
)

const defaultRetryTimes = 3

var errTimeout = errors.New("retry timeout")

type (
	// RetryOption defines the method to customize DoWithRetry.
	RetryOption func(*retryOptions)

	retryOptions struct {
		times    int
		interval time.Duration
		timeout  time.Duration
	}
)

// DoWithRetry runs fn, and retries if failed. Default to retry 3 times.
// Note that if the fn function accesses global variables outside the function and performs modification operations,
// it is best to lock them, otherwise there may be data race issues
func DoWithRetry(fn func() error, opts ...RetryOption) error {
	return retry(fn, opts...)
}

// DoWithRetryCtx runs fn, and retries if failed. Default to retry 3 times.
// fn retryCount indicates the current number of retries,starting from 0
// Note that if the fn function accesses global variables outside the function and performs modification operations,
// it is best to lock them, otherwise there may be data race issues
func DoWithRetryCtx(fn func(ctx context.Context, retryCount int) error, opts ...RetryOption) error {
	return retry(fn, opts...)
}

func retry(fn interface{}, opts ...RetryOption) error {
	options := newRetryOptions()
	for _, opt := range opts {
		opt(options)
	}

	sign := make(chan error, 1)
	var berr errorx.BatchError
	var cancelFunc context.CancelFunc
	ctx := context.Background()
	if options.timeout > 0 {
		ctx, cancelFunc = context.WithTimeout(ctx, options.timeout)
		defer cancelFunc()
	}

	for i := 0; i < options.times; i++ {
		go func(retryCount int) {
			switch f := fn.(type) {
			case func() error:
				sign <- f()
			case func(ctx context.Context, retryCount int) error:
				sign <- f(ctx, retryCount)
			}
		}(i)

		select {
		case err := <-sign:
			if err != nil {
				berr.Add(err)
			} else {
				return nil
			}
		case <-ctx.Done():
			berr.Add(errTimeout)
			return berr.Err()
		}

		if options.interval > 0 {
			select {
			case <-ctx.Done():
				berr.Add(errTimeout)
				return berr.Err()
			case <-time.After(options.interval):
			}
		}
	}

	return berr.Err()
}

// WithRetry customize a DoWithRetry call with given retry times.
func WithRetry(times int) RetryOption {
	return func(options *retryOptions) {
		options.times = times
	}
}

func WithInterval(interval time.Duration) RetryOption {
	return func(options *retryOptions) {
		options.interval = interval
	}
}

func WithTimeout(timeout time.Duration) RetryOption {
	return func(options *retryOptions) {
		options.timeout = timeout
	}
}

func newRetryOptions() *retryOptions {
	return &retryOptions{
		times:    defaultRetryTimes,
		interval: 0,
		timeout:  0,
	}
}
