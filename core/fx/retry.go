package fx

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/errorx"
)

const defaultRetryTimes = 3

type (
	// RetryOption defines the method to customize DoWithRetry.
	RetryOption func(*retryOptions)

	retryOptions struct {
		times        int
		interval     time.Duration
		timeout      time.Duration
		ignoreErrors []error
	}
)

// DoWithRetry runs fn, and retries if failed. Default to retry 3 times.
// Note that if the fn function accesses global variables outside the function
// and performs modification operations, it is best to lock them,
// otherwise there may be data race issues
func DoWithRetry(fn func() error, opts ...RetryOption) error {
	return retry(context.Background(), func(errChan chan error, retryCount int) {
		errChan <- fn()
	}, opts...)
}

// DoWithRetryCtx runs fn, and retries if failed. Default to retry 3 times.
// fn retryCount indicates the current number of retries, starting from 0
// Note that if the fn function accesses global variables outside the function
// and performs modification operations, it is best to lock them,
// otherwise there may be data race issues
func DoWithRetryCtx(ctx context.Context, fn func(ctx context.Context, retryCount int) error,
	opts ...RetryOption) error {
	return retry(ctx, func(errChan chan error, retryCount int) {
		errChan <- fn(ctx, retryCount)
	}, opts...)
}

func retry(ctx context.Context, fn func(errChan chan error, retryCount int), opts ...RetryOption) error {
	options := newRetryOptions()
	for _, opt := range opts {
		opt(options)
	}

	var berr errorx.BatchError
	var cancelFunc context.CancelFunc
	if options.timeout > 0 {
		ctx, cancelFunc = context.WithTimeout(ctx, options.timeout)
		defer cancelFunc()
	}

	errChan := make(chan error, 1)
	for i := 0; i < options.times; i++ {
		go fn(errChan, i)

		select {
		case err := <-errChan:
			if err != nil {
				for _, ignoreErr := range options.ignoreErrors {
					if errors.Is(err, ignoreErr) {
						return nil
					}
				}
				berr.Add(err)
			} else {
				return nil
			}
		case <-ctx.Done():
			berr.Add(ctx.Err())
			return berr.Err()
		}

		if options.interval > 0 {
			select {
			case <-ctx.Done():
				berr.Add(ctx.Err())
				return berr.Err()
			case <-time.After(options.interval):
			}
		}
	}

	return berr.Err()
}

// WithIgnoreErrors Ignore the specified errors
func WithIgnoreErrors(ignoreErrors []error) RetryOption {
	return func(options *retryOptions) {
		options.ignoreErrors = ignoreErrors
	}
}

// WithInterval customizes a DoWithRetry call with given interval.
func WithInterval(interval time.Duration) RetryOption {
	return func(options *retryOptions) {
		options.interval = interval
	}
}

// WithRetry customizes a DoWithRetry call with given retry times.
func WithRetry(times int) RetryOption {
	return func(options *retryOptions) {
		options.times = times
	}
}

// WithTimeout customizes a DoWithRetry call with given timeout.
func WithTimeout(timeout time.Duration) RetryOption {
	return func(options *retryOptions) {
		options.timeout = timeout
	}
}

func newRetryOptions() *retryOptions {
	return &retryOptions{
		times: defaultRetryTimes,
	}
}
