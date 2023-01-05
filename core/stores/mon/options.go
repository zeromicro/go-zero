package mon

import (
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

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

// WithTimeout sets the timeout for the mongo client.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.SetTimeout(timeout)
	}
}
