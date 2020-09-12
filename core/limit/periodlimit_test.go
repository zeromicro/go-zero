package limit

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stores/redis"
)

func TestPeriodLimit_Take(t *testing.T) {
	testPeriodLimit(t)
}

func TestPeriodLimit_TakeWithAlign(t *testing.T) {
	testPeriodLimit(t, Align())
}

func TestPeriodLimit_RedisUnavailable(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)

	const (
		seconds = 1
		total   = 100
		quota   = 5
	)
	l := NewPeriodLimit(seconds, quota, redis.NewRedis(s.Addr(), redis.NodeType), "periodlimit")
	s.Close()
	val, err := l.Take("first")
	assert.NotNil(t, err)
	assert.Equal(t, 0, val)
}

func testPeriodLimit(t *testing.T, opts ...LimitOption) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

	const (
		seconds = 1
		total   = 100
		quota   = 5
	)
	l := NewPeriodLimit(seconds, quota, redis.NewRedis(s.Addr(), redis.NodeType), "periodlimit", opts...)
	var allowed, hitQuota, overQuota int
	for i := 0; i < total; i++ {
		val, err := l.Take("first")
		if err != nil {
			t.Error(err)
		}
		switch val {
		case Allowed:
			allowed++
		case HitQuota:
			hitQuota++
		case OverQuota:
			overQuota++
		default:
			t.Error("unknown status")
		}
	}

	assert.Equal(t, quota-1, allowed)
	assert.Equal(t, 1, hitQuota)
	assert.Equal(t, total-quota, overQuota)
}
