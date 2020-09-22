package limit

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stores/redis"
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
	l := NewTokenLimiter(rate, burst, redis.NewRedis(s.Addr(), redis.NodeType), "tokenlimit")
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
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, redis.NewRedis(s.Addr(), redis.NodeType), "tokenlimit")
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
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

	const (
		total = 100
		rate  = 5
		burst = 10
	)
	l := NewTokenLimiter(rate, burst, redis.NewRedis(s.Addr(), redis.NodeType), "tokenlimit")
	var allowed int
	for i := 0; i < total; i++ {
		if l.Allow() {
			allowed++
		}
	}

	assert.True(t, allowed >= burst)
}
