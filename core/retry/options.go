package retry

import (
	"github.com/tal-tech/go-zero/core/retry/backoff"
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
func WithBackoff(backoffFunc backoff.Func) *CallOption {
	return &CallOption{apply: func(o *options) {
		o.backoffFunc = backoffFunc
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
