package cache

import (
	"time"

	"zero/core/stores/internal"
)

type Option = internal.Option

func WithExpiry(expiry time.Duration) Option {
	return func(o *internal.Options) {
		o.Expiry = expiry
	}
}

func WithNotFoundExpiry(expiry time.Duration) Option {
	return func(o *internal.Options) {
		o.NotFoundExpiry = expiry
	}
}
