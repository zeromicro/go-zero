package limit

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
)

func init() {
	logx.Disable()
}

func TestTokenLimit_Rescue(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, redis.New(s.Addr()), "tokenlimit")
	s.Close()

	var allowed int
	for i := 0; i < total; i++ {
		time.Sleep(time.Second / time.Duration(total))
		if i == total>>1 {
			assert.Nil(t, s.Restart())
		}
		if l.Allow() {
			allowed++
		}

		// make sure start monitor more than once doesn't matter
		l.startMonitor()
	}

	assert.True(t, allowed >= burst+rate)
}

func TestTokenLimit_Take(t *testing.T) {
	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, store, "tokenlimit")
	var allowed int
	for i := 0; i < total; i++ {
		time.Sleep(time.Second / time.Duration(total))
		if l.Allow() {
			allowed++
		}
	}

	assert.True(t, allowed >= burst+rate)
}

func TestTokenLimit_TakeBurst(t *testing.T) {
	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, store, "tokenlimit")
	var allowed int
	for i := 0; i < total; i++ {
		if l.Allow() {
			allowed++
		}
	}

	assert.True(t, allowed >= burst)
}

func TestTokenLimit_AllowN(t *testing.T) {
	testCases := []struct {
		name         string
		rate         int
		burst        int
		want_allowed int
	}{
		{
			name:         "rate10000",
			rate:         10000,
			burst:        20000,
			want_allowed: 20011,
		},
		{
			name:         "rate5",
			rate:         5,
			burst:        10,
			want_allowed: 15,
		},
	}
	const total = 11

	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	for _, ts := range testCases {
		t.Run(ts.name, func(t *testing.T) {
			got := 0
			now := time.Unix(1642216309, 0)
			lmt := NewTokenLimiter(ts.rate, ts.burst, store, ts.name)

			// take all token
			if lmt.AllowN(now, ts.burst) {
				got += ts.burst
			}

			for i := 0; i < total; i++ {
				now = now.Add(100 * time.Millisecond)
				if lmt.AllowN(now, 1) {
					got++
				}
			}

			assert.Equal(t, ts.want_allowed, got)
		})
	}
}
