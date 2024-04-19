package breaker

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/stat"
)

func init() {
	stat.SetReporter(nil)
}

func TestCircuitBreaker_Allow(t *testing.T) {
	t.Run("allow", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		_, err := b.Allow()
		assert.Nil(t, err)
	})

	t.Run("allow with ctx", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		_, err := b.AllowCtx(context.Background())
		assert.Nil(t, err)
	})

	t.Run("allow with ctx timeout", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()
		time.Sleep(time.Millisecond)
		_, err := b.AllowCtx(ctx)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("allow with ctx cancel", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		for i := 0; i < 100; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			cancel()
			_, err := b.AllowCtx(ctx)
			assert.ErrorIs(t, err, context.Canceled)
		}
		_, err := b.AllowCtx(context.Background())
		assert.NoError(t, err)
	})
}

func TestCircuitBreaker_Do(t *testing.T) {
	t.Run("do", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		err := b.Do(func() error {
			return nil
		})
		assert.Nil(t, err)
	})

	t.Run("do with ctx", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		err := b.DoCtx(context.Background(), func() error {
			return nil
		})
		assert.Nil(t, err)
	})

	t.Run("do with ctx timeout", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()
		time.Sleep(time.Millisecond)
		err := b.DoCtx(ctx, func() error {
			return nil
		})
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("do with ctx cancel", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		for i := 0; i < 100; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			cancel()
			err := b.DoCtx(ctx, func() error {
				return nil
			})
			assert.ErrorIs(t, err, context.Canceled)
		}
		assert.NoError(t, b.DoCtx(context.Background(), func() error {
			return nil
		}))
	})
}

func TestCircuitBreaker_DoWithAcceptable(t *testing.T) {
	t.Run("doWithAcceptable", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		err := b.DoWithAcceptable(func() error {
			return nil
		}, func(err error) bool {
			return true
		})
		assert.Nil(t, err)
	})

	t.Run("doWithAcceptable with ctx", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		err := b.DoWithAcceptableCtx(context.Background(), func() error {
			return nil
		}, func(err error) bool {
			return true
		})
		assert.Nil(t, err)
	})

	t.Run("doWithAcceptable with ctx timeout", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()
		time.Sleep(time.Millisecond)
		err := b.DoWithAcceptableCtx(ctx, func() error {
			return nil
		}, func(err error) bool {
			return true
		})
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("doWithAcceptable with ctx cancel", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		for i := 0; i < 100; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			cancel()
			err := b.DoWithAcceptableCtx(ctx, func() error {
				return nil
			}, func(err error) bool {
				return true
			})
			assert.ErrorIs(t, err, context.Canceled)
		}
		assert.NoError(t, b.DoWithAcceptableCtx(context.Background(), func() error {
			return nil
		}, func(err error) bool {
			return true
		}))
	})
}

func TestCircuitBreaker_DoWithFallback(t *testing.T) {
	t.Run("doWithFallback", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		err := b.DoWithFallback(func() error {
			return nil
		}, func(err error) error {
			return err
		})
		assert.Nil(t, err)
	})

	t.Run("doWithFallback with ctx", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		err := b.DoWithFallbackCtx(context.Background(), func() error {
			return nil
		}, func(err error) error {
			return err
		})
		assert.Nil(t, err)
	})

	t.Run("doWithFallback with ctx timeout", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()
		time.Sleep(time.Millisecond)
		err := b.DoWithFallbackCtx(ctx, func() error {
			return nil
		}, func(err error) error {
			return err
		})
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("doWithFallback with ctx cancel", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		for i := 0; i < 100; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			cancel()
			err := b.DoWithFallbackCtx(ctx, func() error {
				return nil
			}, func(err error) error {
				return err
			})
			assert.ErrorIs(t, err, context.Canceled)
		}
		assert.NoError(t, b.DoWithFallbackCtx(context.Background(), func() error {
			return nil
		}, func(err error) error {
			return err
		}))
	})
}

func TestCircuitBreaker_DoWithFallbackAcceptable(t *testing.T) {
	t.Run("doWithFallbackAcceptable", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		err := b.DoWithFallbackAcceptable(func() error {
			return nil
		}, func(err error) error {
			return err
		}, func(err error) bool {
			return true
		})
		assert.Nil(t, err)
	})

	t.Run("doWithFallbackAcceptable with ctx", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		err := b.DoWithFallbackAcceptableCtx(context.Background(), func() error {
			return nil
		}, func(err error) error {
			return err
		}, func(err error) bool {
			return true
		})
		assert.Nil(t, err)
	})

	t.Run("doWithFallbackAcceptable with ctx timeout", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
		defer cancel()
		time.Sleep(time.Millisecond)
		err := b.DoWithFallbackAcceptableCtx(ctx, func() error {
			return nil
		}, func(err error) error {
			return err
		}, func(err error) bool {
			return true
		})
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("doWithFallbackAcceptable with ctx cancel", func(t *testing.T) {
		b := NewBreaker()
		assert.True(t, len(b.Name()) > 0)
		for i := 0; i < 100; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			cancel()
			err := b.DoWithFallbackAcceptableCtx(ctx, func() error {
				return nil
			}, func(err error) error {
				return err
			}, func(err error) bool {
				return true
			})
			assert.ErrorIs(t, err, context.Canceled)
		}
		assert.NoError(t, b.DoWithFallbackAcceptableCtx(context.Background(), func() error {
			return nil
		}, func(err error) error {
			return err
		}, func(err error) bool {
			return true
		}))
	})
}

func TestLogReason(t *testing.T) {
	b := NewBreaker()
	assert.True(t, len(b.Name()) > 0)

	for i := 0; i < 1000; i++ {
		_ = b.Do(func() error {
			return errors.New(strconv.Itoa(i))
		})
	}
	errs := b.(*circuitBreaker).throttle.(loggedThrottle).errWin
	assert.Equal(t, numHistoryReasons, errs.count)
}

func TestErrorWindow(t *testing.T) {
	tests := []struct {
		name    string
		reasons []string
	}{
		{
			name: "no error",
		},
		{
			name:    "one error",
			reasons: []string{"foo"},
		},
		{
			name:    "two errors",
			reasons: []string{"foo", "bar"},
		},
		{
			name:    "five errors",
			reasons: []string{"first", "second", "third", "fourth", "fifth"},
		},
		{
			name:    "six errors",
			reasons: []string{"first", "second", "third", "fourth", "fifth", "sixth"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var ew errorWindow
			for _, reason := range test.reasons {
				ew.add(reason)
			}
			var reasons []string
			if len(test.reasons) > numHistoryReasons {
				reasons = test.reasons[len(test.reasons)-numHistoryReasons:]
			} else {
				reasons = test.reasons
			}
			for _, reason := range reasons {
				assert.True(t, strings.Contains(ew.String(), reason), fmt.Sprintf("actual: %s", ew.String()))
			}
		})
	}
}

func TestPromiseWithReason(t *testing.T) {
	tests := []struct {
		name   string
		reason string
		expect string
	}{
		{
			name: "success",
		},
		{
			name:   "success",
			reason: "fail",
			expect: "fail",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			promise := promiseWithReason{
				promise: new(mockedPromise),
				errWin:  new(errorWindow),
			}
			if len(test.reason) == 0 {
				promise.Accept()
			} else {
				promise.Reject(test.reason)
			}

			assert.True(t, strings.Contains(promise.errWin.String(), test.expect))
		})
	}
}

func BenchmarkGoogleBreaker(b *testing.B) {
	br := NewBreaker()
	for i := 0; i < b.N; i++ {
		_ = br.Do(func() error {
			return nil
		})
	}
}

type mockedPromise struct{}

func (m *mockedPromise) Accept() {
}

func (m *mockedPromise) Reject() {
}
