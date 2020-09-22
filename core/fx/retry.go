package fx

import "github.com/tal-tech/go-zero/core/errorx"

const defaultRetryTimes = 3

type (
	RetryOption func(*retryOptions)

	retryOptions struct {
		times int
	}
)

func DoWithRetries(fn func() error, opts ...RetryOption) error {
	var options = newRetryOptions()
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

func WithRetries(times int) RetryOption {
	return func(options *retryOptions) {
		options.times = times
	}
}

func newRetryOptions() *retryOptions {
	return &retryOptions{
		times: defaultRetryTimes,
	}
}
