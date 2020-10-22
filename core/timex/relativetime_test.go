package timex

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRelativeTime(t *testing.T) {
	time.Sleep(time.Millisecond)
	now := Now()
	assert.True(t, now > 0)
	time.Sleep(time.Millisecond)
	assert.True(t, Since(now) > 0)
}

func TestRelativeTime_Time(t *testing.T) {
	diff := time.Until(Time())
	if diff > 0 {
		assert.True(t, diff < time.Second)
	} else {
		assert.True(t, -diff < time.Second)
	}
}
