package mongo

import "time"

type (
	options struct {
		timeout       time.Duration
		slowThreshold time.Duration
	}

	// Option defines the method to customize a mongo model.
	Option func(opts *options)
)

func defaultOptions() *options {
	return &options{
		timeout:       defaultTimeout,
		slowThreshold: defaultSlowThreshold,
	}
}
