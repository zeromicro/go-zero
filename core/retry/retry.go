package retry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var (
	// Permanent is a wrapper function when you don't want retry
	// with some inner error.
	Permanent = backoff.Permanent
)

type NotifyWithTimes func(err error, duration time.Duration, times int)

// PolicyType denotes if the back off delay should be constant or exponential.
type PolicyType int

const (
	// PolicyConstant is a backoff policy that always returns the same backoff delay.
	PolicyConstant PolicyType = iota + 1
	// PolicyExponential is a backoff implementation that increases the backoff period
	// for each retry attempt using a randomization function that grows exponentially.
	PolicyExponential
)

// Config encapsulates the back off policy configuration.
type Config struct {
	Policy PolicyType

	// Constant back off
	Duration time.Duration

	// Exponential back off
	InitialInterval     time.Duration
	RandomizationFactor float32
	Multiplier          float32
	MaxInterval         time.Duration
	MaxElapsedTime      time.Duration

	// Additional options
	MaxRetries int64
}

// DefaultConfig represents the default configuration for a
// `Config`.
func DefaultConfig() Config {
	return Config{
		Policy:              PolicyExponential,
		Duration:            5 * time.Second,
		InitialInterval:     backoff.DefaultInitialInterval,     // 500ms
		RandomizationFactor: backoff.DefaultRandomizationFactor, // 0.5
		Multiplier:          backoff.DefaultMultiplier,          // 1.5
		MaxInterval:         backoff.DefaultMaxInterval,         // 60s
		MaxElapsedTime:      backoff.DefaultMaxElapsedTime,      // 15min
		MaxRetries:          10,
	}
}

// DefaultConfigWithNoRetry represents the default configuration with `MaxRetries` set to 0.
// This may be useful for those brokers which can handles retries on its own.
func DefaultConfigWithNoRetry() Config {
	c := DefaultConfig()
	c.MaxRetries = 0

	return c
}

func (c *Config) NewBackOff() backoff.BackOff {
	var b backoff.BackOff
	switch c.Policy {
	case PolicyConstant:
		b = backoff.NewConstantBackOff(c.Duration)
	case PolicyExponential:
		eb := backoff.NewExponentialBackOff()
		eb.InitialInterval = c.InitialInterval
		eb.RandomizationFactor = float64(c.RandomizationFactor)
		eb.Multiplier = float64(c.Multiplier)
		eb.MaxInterval = c.MaxInterval
		eb.MaxElapsedTime = c.MaxElapsedTime
		b = eb
	}

	if c.MaxRetries >= 0 {
		b = backoff.WithMaxRetries(b, uint64(c.MaxRetries))
	}

	return b
}

// NewBackOffWithContext add context deadline control for backoff
func (c *Config) NewBackOffWithContext(ctx context.Context) backoff.BackOff {
	b := c.NewBackOff()

	return backoff.WithContext(b, ctx)
}

// NotifyRecover is a wrapper around backoff.RetryNotify that adds another callback for when an operation
// previously failed but has since recovered. The main purpose of this wrapper is to call `notify` only when
// the operations fail the first time and `recovered` when it finally succeeds. This can be helpful in limiting
// log messages to only the events that operators need to be alerted on.
func NotifyRecover(operation backoff.Operation, b backoff.BackOff, notify NotifyWithTimes, recovered func(times int), verbose bool) error {
	var notified bool

	i := 0

	return backoff.RetryNotify(func() error {
		err := operation()

		if err == nil && notified {
			notified = false
			recovered(i)
		}

		return err
	}, b, func(err error, d time.Duration) {
		i++

		if !verbose && notified {
			return
		}

		notify(err, d, i)
		notified = true
	})
}
