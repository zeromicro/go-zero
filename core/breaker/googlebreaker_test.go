package breaker

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/stat"
)

const (
	testBuckets  = 10
	testInterval = time.Millisecond * 10
)

func init() {
	stat.SetReporter(nil)
}

func getGoogleBreaker() *googleBreaker {
	st := collection.NewRollingWindow(testBuckets, testInterval)
	return &googleBreaker{
		stat:  st,
		k:     5,
		proba: mathx.NewProba(),
	}
}

func markSuccessWithDuration(b *googleBreaker, count int, sleep time.Duration) {
	for i := 0; i < count; i++ {
		b.markSuccess()
		time.Sleep(sleep)
	}
}

func markFailedWithDuration(b *googleBreaker, count int, sleep time.Duration) {
	for i := 0; i < count; i++ {
		b.markFailure()
		time.Sleep(sleep)
	}
}

func TestGoogleBreakerClose(t *testing.T) {
	b := getGoogleBreaker()
	markSuccess(b, 80)
	assert.Nil(t, b.accept())
	markSuccess(b, 120)
	assert.Nil(t, b.accept())
}

func TestGoogleBreakerOpen(t *testing.T) {
	b := getGoogleBreaker()
	markSuccess(b, 10)
	assert.Nil(t, b.accept())
	markFailed(b, 100000)
	time.Sleep(testInterval * 2)
	verify(t, func() bool {
		return b.accept() != nil
	})
}

func TestGoogleBreakerFallback(t *testing.T) {
	b := getGoogleBreaker()
	markSuccess(b, 1)
	assert.Nil(t, b.accept())
	markFailed(b, 10000)
	time.Sleep(testInterval * 2)
	verify(t, func() bool {
		return b.doReq(func() error {
			return errors.New("any")
		}, func(err error) error {
			return nil
		}, defaultAcceptable) == nil
	})
}

func TestGoogleBreakerReject(t *testing.T) {
	b := getGoogleBreaker()
	markSuccess(b, 100)
	assert.Nil(t, b.accept())
	markFailed(b, 10000)
	time.Sleep(testInterval)
	assert.Equal(t, ErrServiceUnavailable, b.doReq(func() error {
		return ErrServiceUnavailable
	}, nil, defaultAcceptable))
}

func TestGoogleBreakerAcceptable(t *testing.T) {
	b := getGoogleBreaker()
	errAcceptable := errors.New("any")
	assert.Equal(t, errAcceptable, b.doReq(func() error {
		return errAcceptable
	}, nil, func(err error) bool {
		return errors.Is(err, errAcceptable)
	}))
}

func TestGoogleBreakerNotAcceptable(t *testing.T) {
	b := getGoogleBreaker()
	errAcceptable := errors.New("any")
	assert.Equal(t, errAcceptable, b.doReq(func() error {
		return errAcceptable
	}, nil, func(err error) bool {
		return !errors.Is(err, errAcceptable)
	}))
}

func TestGoogleBreakerPanic(t *testing.T) {
	b := getGoogleBreaker()
	assert.Panics(t, func() {
		_ = b.doReq(func() error {
			panic("fail")
		}, nil, defaultAcceptable)
	})
}

func TestGoogleBreakerHalfOpen(t *testing.T) {
	b := getGoogleBreaker()
	assert.Nil(t, b.accept())
	t.Run("accept single failed/accept", func(t *testing.T) {
		markFailed(b, 10000)
		time.Sleep(testInterval * 2)
		verify(t, func() bool {
			return b.accept() != nil
		})
	})
	t.Run("accept single failed/allow", func(t *testing.T) {
		markFailed(b, 10000)
		time.Sleep(testInterval * 2)
		verify(t, func() bool {
			_, err := b.allow()
			return err != nil
		})
	})
	time.Sleep(testInterval * testBuckets)
	t.Run("accept single succeed", func(t *testing.T) {
		assert.Nil(t, b.accept())
		markSuccess(b, 10000)
		verify(t, func() bool {
			return b.accept() == nil
		})
	})
}

func TestGoogleBreakerSelfProtection(t *testing.T) {
	t.Run("total request < 100", func(t *testing.T) {
		b := getGoogleBreaker()
		markFailed(b, 4)
		time.Sleep(testInterval)
		assert.Nil(t, b.accept())
	})
	t.Run("total request > 100, total < 2 * success", func(t *testing.T) {
		b := getGoogleBreaker()
		size := rand.Intn(10000)
		accepts := size + 1
		markSuccess(b, accepts)
		markFailed(b, size-accepts)
		assert.Nil(t, b.accept())
	})
}

func TestGoogleBreakerHistory(t *testing.T) {
	var b *googleBreaker
	var accepts, total int64

	sleep := testInterval
	t.Run("accepts == total", func(t *testing.T) {
		b = getGoogleBreaker()
		markSuccessWithDuration(b, 10, sleep/2)
		accepts, total = b.history()
		assert.Equal(t, int64(10), accepts)
		assert.Equal(t, int64(10), total)
	})

	t.Run("fail == total", func(t *testing.T) {
		b = getGoogleBreaker()
		markFailedWithDuration(b, 10, sleep/2)
		accepts, total = b.history()
		assert.Equal(t, int64(0), accepts)
		assert.Equal(t, int64(10), total)
	})

	t.Run("accepts = 1/2 * total, fail = 1/2 * total", func(t *testing.T) {
		b = getGoogleBreaker()
		markFailedWithDuration(b, 5, sleep/2)
		markSuccessWithDuration(b, 5, sleep/2)
		accepts, total = b.history()
		assert.Equal(t, int64(5), accepts)
		assert.Equal(t, int64(10), total)
	})

	t.Run("auto reset rolling counter", func(t *testing.T) {
		b = getGoogleBreaker()
		time.Sleep(testInterval * testBuckets)
		accepts, total = b.history()
		assert.Equal(t, int64(0), accepts)
		assert.Equal(t, int64(0), total)
	})
}

func BenchmarkGoogleBreakerAllow(b *testing.B) {
	breaker := getGoogleBreaker()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		_ = breaker.accept()
		if i%2 == 0 {
			breaker.markSuccess()
		} else {
			breaker.markFailure()
		}
	}
}

func BenchmarkGoogleBreakerDoReq(b *testing.B) {
	breaker := getGoogleBreaker()
	b.ResetTimer()
	for i := 0; i <= b.N; i++ {
		_ = breaker.doReq(func() error {
			return nil
		}, nil, defaultAcceptable)
	}
}

func markSuccess(b *googleBreaker, count int) {
	for i := 0; i < count; i++ {
		p, err := b.allow()
		if err != nil {
			break
		}
		p.Accept()
	}
}

func markFailed(b *googleBreaker, count int) {
	for i := 0; i < count; i++ {
		p, err := b.allow()
		if err == nil {
			p.Reject()
		}
	}
}
