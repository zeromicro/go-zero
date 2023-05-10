package fx

import "github.com/zeromicro/go-zero/core/errorx"

const defaultRetryTimes = 3

type (
	// RetryOption defines the method to customize DoWithRetry.
	RetryOption func(*retryOptions)

	retryOptions struct {
		times int
		// errHandler determine if an error is returned
		errHandler func(error) bool
	}
)

// DoWithRetry runs fn, and retries if failed. Default to retry 3 times.
func DoWithRetry(fn func() error, opts ...RetryOption) error {
	options := newRetryOptions()
	for _, opt := range opts {
		opt(options)
	}

	var batchErr errorx.BatchError
	var err error
	for i := 0; i < options.times; i++ {
		if err = fn(); err == nil {
			return nil
		}
		batchErr.Add(err)
		// return if error handler return true
		if options.errHandler != nil && options.errHandler(err) {
			return batchErr.Err()
		}
	}

	return batchErr.Err()
}

// WithRetry customize a DoWithRetry call with given retry times.
func WithRetry(times int) RetryOption {
	return func(options *retryOptions) {
		options.times = times
	}
}

// WithErrHandler customize a error handler
func WithErrHandler(errHandler func(error) bool) RetryOption {
	return func(options *retryOptions) {
		options.errHandler = errHandler
	}
}

func newRetryOptions() *retryOptions {
	return &retryOptions{
		times: defaultRetryTimes,
	}
}
