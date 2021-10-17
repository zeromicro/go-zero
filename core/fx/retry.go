package fx

import "github.com/zeromicro/go-zero/core/errorx"

const defaultRetryTimes = 3

type (
	// RetryOption defines the method to customize DoWithRetry.
	RetryOption func(*retryOptions)

	retryOptions struct {
		times int
	}
)

// DoWithRetry runs fn, and retries if failed. Default to retry 3 times.
func DoWithRetry(fn func() error, opts ...RetryOption) error {
	options := newRetryOptions()
	for _, opt := range opts {
		opt(options)
	}

	var berr errorx.BatchError
	for i := 0; i < options.times; i++ {
		if err := fn(); err != nil {
			berr.Add(err)
		} else {
			return nil
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

func newRetryOptions() *retryOptions {
	return &retryOptions{
		times: defaultRetryTimes,
	}
}
