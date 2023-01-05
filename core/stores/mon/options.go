package mon

import (
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
	mopt "go.mongodb.org/mongo-driver/mongo/options"
)

var slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

type (
	options struct {
		timeout time.Duration
	}

	// Option defines the method to customize a mongo model.
	Option func(opts *options)
)

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

func defaultOptions() *options {
	return &options{
		timeout: defaultTimeout,
	}
}

// WithTimeout sets the timeout for the mongo client.
func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		opts.timeout = timeout
	}
}

func (opts *options) mgoOptions() []*mopt.ClientOptions {
	var mOpts []*mopt.ClientOptions
	if opts.timeout > 0 {
		mOpts = append(mOpts, mopt.Client().SetTimeout(opts.timeout))
	}
	return mOpts
}
