package mongo

import (
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
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
