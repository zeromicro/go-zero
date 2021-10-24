package retrybackoff

import (
	"context"
	"math/rand"
	"time"
)

type BackoffFunc func(ctx context.Context, attempt uint) time.Duration

func BackoffLinearWithJitter(waitBetween time.Duration, jitterFraction float64) BackoffFunc {
	return func(ctx context.Context, attempt uint) time.Duration {
		multiplier := jitterFraction * (rand.Float64()*2 - 1)
		return time.Duration(float64(waitBetween) * (1 + multiplier))
	}
}
