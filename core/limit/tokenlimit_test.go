package limit

import (
	"context"
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

func TestTokenLimit_WithCtx(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, redis.New(s.Addr()), "tokenlimit")
	defer s.Close()

	ctx, cancel := context.WithCancel(context.Background())
	ok := l.AllowCtx(ctx)
	assert.True(t, ok)

	cancel()
	for i := 0; i < total; i++ {
		ok := l.AllowCtx(ctx)
		assert.False(t, ok)
		assert.False(t, l.monitorStarted)
	}
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
	store := redistest.CreateRedis(t)

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

func TestTokenLimit_SharedBurst(t *testing.T) {
	tests := []struct {
		name   string
		rate   int
		burst  int
		expect time.Duration
	}{
		{
			name:   "small burst",
			rate:   100,
			burst:  10,
			expect: time.Second,
		},
		{
			name:   "minimum burst",
			rate:   100,
			burst:  1,
			expect: time.Second,
		},
		{
			name:   "ttl boundary",
			rate:   20,
			burst:  10,
			expect: time.Second,
		},
		{
			name:   "fractional ttl",
			rate:   3,
			burst:  4,
			expect: 2 * time.Second,
		},
		{
			name:   "large burst",
			rate:   5,
			burst:  10,
			expect: 4 * time.Second,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := miniredis.RunT(t)
			store := redis.New(s.Addr())
			l := NewTokenLimiter(test.rate, test.burst, store, "tokenlimit")
			other := NewTokenLimiter(test.rate, test.burst, store, "tokenlimit")
			now := time.Now()

			assert.True(t, l.AllowN(now, test.burst))
			assert.Equal(t, test.expect, s.TTL(l.tokenKey))
			assert.Equal(t, test.expect, s.TTL(l.timestampKey))
			assert.False(t, other.AllowN(now, 1))

			now = now.Add(10 * time.Second)
			assert.True(t, other.AllowN(now, test.burst))
			assert.False(t, l.AllowN(now, 1))
		})
	}
}

func TestTokenLimit_TakeBurst(t *testing.T) {
	store := redistest.CreateRedis(t)

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
