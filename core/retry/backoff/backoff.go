package backoff

import (
	"math/rand"
	"time"
)

// Func defines the method to calculate how long to retry.
type Func func(attempt int) time.Duration

// LinearWithJitter waits a set period of time, allowing for jitter (fractional adjustment).
func LinearWithJitter(waitBetween time.Duration, jitterFraction float64) Func {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func(attempt int) time.Duration {
		multiplier := jitterFraction * (r.Float64()*2 - 1)
		return time.Duration(float64(waitBetween) * (1 + multiplier))
	}
}

// Interval it waits for a fixed period of time between calls.
func Interval(interval time.Duration) Func {
	return func(attempt int) time.Duration {
		return interval
	}
}

// Exponential produces increasing intervals for each attempt.
func Exponential(scalar time.Duration) Func {
	return func(attempt int) time.Duration {
		return scalar * time.Duration((1<<attempt)>>1)
	}
}
