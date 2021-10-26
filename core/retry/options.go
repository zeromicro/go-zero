package retry

import (
	"google.golang.org/grpc/codes"
	"time"
)

// WithDisable disables the retry behaviour on this call, or this interceptor.
//
// Its semantically the same to `WithMax`
func WithDisable() *CallOption {
	return WithMax(0)
}

// WithMax sets the maximum number of retries on this call, or this interceptor.
func WithMax(maxRetries int) *CallOption {
	return &CallOption{apply: func(options *options) {
		options.max = maxRetries
	}}
}

// WithBackoff sets the `BackoffFunc` used to control time between retries.
func WithBackoff(backoffFunc func(attempt int) time.Duration) *CallOption {
	return &CallOption{apply: func(o *options) {
		o.backoffFunc = func(attempt int) time.Duration {
			return backoffFunc(attempt)
		}
	}}
}

// WithCodes Allow code to be retried.
func WithCodes(retryCodes ...codes.Code) *CallOption {
	return &CallOption{apply: func(o *options) {
		o.codes = retryCodes
	}}
}

// WithPerRetryTimeout timeout for each retry
func WithPerRetryTimeout(timeout time.Duration) *CallOption {
	return &CallOption{apply: func(o *options) {
		o.perCallTimeout = timeout
	}}
}
