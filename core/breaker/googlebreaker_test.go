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
	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	testBuckets  = 10
	testInterval = time.Millisecond * 10
)

func init() {
	stat.SetReporter(nil)
}

func getGoogleBreaker() *googleBreaker {
	st := collection.NewRollingWindow[int64, *bucket](func() *bucket {
		return new(bucket)
	}, testBuckets, testInterval)
	return &googleBreaker{
		stat:     st,
		k:        5,
		proba:    mathx.NewProba(),
		lastPass: syncx.NewAtomicDuration(),
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

func TestGoogleBreakerRecover(t *testing.T) {
	st := collection.NewRollingWindow[int64, *bucket](func() *bucket {
		return new(bucket)
	}, testBuckets*2, testInterval)
	b := &googleBreaker{
		stat:     st,
		k:        k,
		proba:    mathx.NewProba(),
		lastPass: syncx.NewAtomicDuration(),
	}
	for i := 0; i < testBuckets; i++ {
		for j := 0; j < 100; j++ {
			b.stat.Add(1)
		}
		time.Sleep(testInterval)
	}
	for i := 0; i < testBuckets; i++ {
		for j := 0; j < 100; j++ {
			b.stat.Add(0)
		}
		time.Sleep(testInterval)
	}
	verify(t, func() bool {
		return b.accept() == nil
	})
}

func TestGoogleBreakerSelfRejectionNotCounted(t *testing.T) {
	// a long window, so no bucket expires while the test is running.
	st := collection.NewRollingWindow[int64, *bucket](func() *bucket {
		return new(bucket)
	}, testBuckets, time.Second)
	b := &googleBreaker{
		stat:     st,
		k:        k,
		proba:    mathx.NewProba(),
		lastPass: syncx.NewAtomicDuration(),
	}
	markFailed(b, 1000)

	var passed int
	before := b.history().total
	for i := 0; i < 100; i++ {
		promise, err := b.allow()
		if err != nil {
			assert.ErrorIs(t, err, ErrServiceUnavailable)
			continue
		}

		passed++
		promise.Reject()
	}

	for i := 0; i < 100; i++ {
		err := b.doReq(func() error {
			return errors.New("service down")
		}, nil, defaultAcceptable)
		if !errors.Is(err, ErrServiceUnavailable) {
			passed++
		}
	}

	// rejections by the breaker itself are not backend outcomes,
	// they must not be counted in the statistics.
	assert.Equal(t, before+int64(passed), b.history().total)
}

func TestGoogleBreakerRecoverAfterSelfRejection(t *testing.T) {
	calls := map[string]func(b *googleBreaker, healthy bool) bool{
		"doReq": func(b *googleBreaker, healthy bool) bool {
			return b.doReq(func() error {
				if healthy {
					return nil
				}

				return errors.New("service down")
			}, nil, defaultAcceptable) == nil
		},
		"allow": func(b *googleBreaker, healthy bool) bool {
			promise, err := b.allow()
			if err != nil {
				return false
			}

			if healthy {
				promise.Accept()
			} else {
				promise.Reject()
			}

			return true
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			// same bucket count as production, so the recovery damping
			// (buckets-workingBuckets)/buckets behaves like production.
			// without the fix the breaker needs more than half a window to recover.
			timeout := testInterval * buckets / 4

			// retries tolerate a slow machine.
			var recovered bool
			for i := 0; i < 3 && !recovered; i++ {
				recovered = recoversWithin(call, timeout)
			}

			assert.True(t, recovered, "breaker should recover within %v", timeout)
		})
	}
}

func recoversWithin(call func(b *googleBreaker, healthy bool) bool, timeout time.Duration) bool {
	st := collection.NewRollingWindow[int64, *bucket](func() *bucket {
		return new(bucket)
	}, buckets, testInterval)
	b := &googleBreaker{
		stat:     st,
		k:        k,
		proba:    mathx.NewProba(),
		lastPass: syncx.NewAtomicDuration(),
	}

	// backend is down under heavy load: real failures open the breaker,
	// then the breaker itself rejects most of the requests.
	for i := 0; i < buckets; i++ {
		for j := 0; j < 1000; j++ {
			call(b, false)
		}

		time.Sleep(testInterval)
	}

	// backend recovers, the breaker must pass traffic again long before
	// the forced probe (forcePassDuration is 1s).
	var consecutive int
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) && consecutive < 50 {
		if call(b, true) {
			consecutive++
		} else {
			consecutive = 0
		}
	}

	return consecutive >= 50
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

func TestGoogleBreakerMoreFallingBuckets(t *testing.T) {
	t.Parallel()

	t.Run("more falling buckets", func(t *testing.T) {
		b := getGoogleBreaker()

		func() {
			stopChan := time.After(testInterval * 6)
			for {
				time.Sleep(time.Millisecond)
				select {
				case <-stopChan:
					return
				default:
					assert.Error(t, b.doReq(func() error {
						return errors.New("foo")
					}, func(err error) error {
						return err
					}, func(err error) bool {
						return err == nil
					}))
				}
			}
		}()

		var count int
		for i := 0; i < 100; i++ {
			if errors.Is(b.doReq(func() error {
				return ErrServiceUnavailable
			}, nil, defaultAcceptable), ErrServiceUnavailable) {
				count++
			}
		}
		assert.True(t, count > 90)
	})
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
	sleep := testInterval
	t.Run("accepts == total", func(t *testing.T) {
		b := getGoogleBreaker()
		markSuccessWithDuration(b, 10, sleep/2)
		result := b.history()
		assert.Equal(t, int64(10), result.accepts)
		assert.Equal(t, int64(10), result.total)
	})

	t.Run("fail == total", func(t *testing.T) {
		b := getGoogleBreaker()
		markFailedWithDuration(b, 10, sleep/2)
		result := b.history()
		assert.Equal(t, int64(0), result.accepts)
		assert.Equal(t, int64(10), result.total)
	})

	t.Run("accepts = 1/2 * total, fail = 1/2 * total", func(t *testing.T) {
		b := getGoogleBreaker()
		markFailedWithDuration(b, 5, sleep/2)
		markSuccessWithDuration(b, 5, sleep/2)
		result := b.history()
		assert.Equal(t, int64(5), result.accepts)
		assert.Equal(t, int64(10), result.total)
	})

	t.Run("auto reset rolling counter", func(t *testing.T) {
		b := getGoogleBreaker()
		time.Sleep(testInterval * testBuckets)
		result := b.history()
		assert.Equal(t, int64(0), result.accepts)
		assert.Equal(t, int64(0), result.total)
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
