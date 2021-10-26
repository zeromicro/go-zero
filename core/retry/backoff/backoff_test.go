package backoff

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWaitBetween(t *testing.T) {
	fn := Interval(time.Second)
	assert.EqualValues(t, time.Second, fn(1))
}

func TestExponential(t *testing.T) {
	fn := Exponential(time.Second)
	assert.EqualValues(t, time.Second, fn(1))
}
