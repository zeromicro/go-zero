package mon

import (
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

const defaultTimeout = time.Second * 3

var slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

type (
	options = mopt.ClientOptions

	// Option defines the method to customize a mongo model.
	Option func(opts *options)
)

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

func defaultTimeoutOption() Option {
	return func(opts *options) {
		opts.SetTimeout(defaultTimeout)
	}
}

// WithTimeout set the mon client operation timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.SetTimeout(timeout)
	}
}
