package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/proc"
	"github.com/zeromicro/go-zero/core/timex"
)

func TestNextDelay(t *testing.T) {
	tests := []struct {
		name   string
		input  time.Duration
		output time.Duration
		ok     bool
	}{
		{
			name:   "second",
			input:  time.Second,
			output: time.Second * 5,
			ok:     true,
		},
		{
			name:   "5 seconds",
			input:  time.Second * 5,
			output: time.Minute,
			ok:     true,
		},
		{
			name:   "minute",
			input:  time.Minute,
			output: time.Minute * 5,
			ok:     true,
		},
		{
			name:   "5 minutes",
			input:  time.Minute * 5,
			output: time.Hour,
			ok:     true,
		},
		{
			name:   "hour",
			input:  time.Hour,
			output: 0,
			ok:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			old := timingWheel.Load()
			ticker := timex.NewFakeTicker()
			tw, err := collection.NewTimingWheelWithTicker(
				time.Millisecond, timingWheelSlots, func(key, value any) {
					clean(key, value)
				}, ticker)
			timingWheel.Store(tw)
			assert.NoError(t, err)
			t.Cleanup(func() {
				timingWheel.Store(old)
			})

			next, ok := nextDelay(test.input)
			assert.Equal(t, test.ok, ok)
			assert.Equal(t, test.output, next)
			proc.Shutdown()
		})
	}
}
