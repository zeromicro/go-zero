package backoff

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitBetween(t *testing.T) {
	fn := Interval(time.Second)
	assert.EqualValues(t, time.Second, fn(1))
}

func TestExponential(t *testing.T) {
	fn := Exponential(time.Second)
	assert.EqualValues(t, time.Second, fn(1))
}

func TestLinearWithJitter(t *testing.T) {
	const rounds = 1000000
	var total time.Duration
	fn := LinearWithJitter(time.Second, 0.5)
	for i := 0; i < rounds; i++ {
		total += fn(1)
	}

	// 0.1% tolerance
	assert.True(t, total/time.Duration(rounds)-time.Second < time.Millisecond)
}
