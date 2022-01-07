package logx

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/timex"
)

func TestLimitedExecutor_logOrDiscard(t *testing.T) {
	tests := []struct {
		name      string
		threshold time.Duration
		lastTime  time.Duration
		discarded uint32
		executed  bool
	}{
		{
			name:     "nil executor",
			executed: true,
		},
		{
			name:      "regular",
			threshold: time.Hour,
			lastTime:  timex.Now(),
			discarded: 10,
			executed:  false,
		},
		{
			name:      "slow",
			threshold: time.Duration(1),
			lastTime:  -1000,
			discarded: 10,
			executed:  true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			executor := newLimitedExecutor(0)
			executor.threshold = test.threshold
			executor.discarded = test.discarded
			executor.lastTime.Set(test.lastTime)

			var run int32
			executor.logOrDiscard(func() {
				atomic.AddInt32(&run, 1)
			})
			if test.executed {
				assert.Equal(t, int32(1), atomic.LoadInt32(&run))
			} else {
				assert.Equal(t, int32(0), atomic.LoadInt32(&run))
				assert.Equal(t, test.discarded+1, atomic.LoadUint32(&executor.discarded))
			}
		})
	}
}
