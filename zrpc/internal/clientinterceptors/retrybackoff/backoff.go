package retrybackoff

import (
	"context"
	"math/rand"
	"time"
)

type BackoffFunc func(ctx context.Context, attempt int) time.Duration

func BackoffLinearWithJitter(waitBetween time.Duration, jitterFraction float64) BackoffFunc {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return func(ctx context.Context, attempt int) time.Duration {
		multiplier := jitterFraction * (r.Float64()*2 - 1)
		return time.Duration(float64(waitBetween) * (1 + multiplier))
	}
}
