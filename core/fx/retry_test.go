package fx

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	assert.NotNil(t, DoWithRetry(func() error {
		return errors.New("any")
	}))

	times1 := 0
	assert.Nil(t, DoWithRetry(func() error {
		times1++
		if times1 == defaultRetryTimes {
			return nil
		}
		return errors.New("any")
	}))

	times2 := 0
	assert.NotNil(t, DoWithRetry(func() error {
		times2++
		if times2 == defaultRetryTimes+1 {
			return nil
		}
		return errors.New("any")
	}))

	total := 2 * defaultRetryTimes
	times3 := 0
	assert.Nil(t, DoWithRetry(func() error {
		times3++
		if times3 == total {
			return nil
		}
		return errors.New("any")
	}, WithRetry(total)))
}

func TestRetryWithTimeout(t *testing.T) {
	assert.Nil(t, DoWithRetry(func() error {
		return nil
	}, WithTimeout(time.Millisecond*500)))

	times1 := 0
	assert.Nil(t, DoWithRetry(func() error {
		times1++
		if times1 == 1 {
			return errors.New("any ")
		}
		time.Sleep(time.Millisecond * 150)
		return nil
	}, WithTimeout(time.Millisecond*250)))

	total := defaultRetryTimes
	times2 := 0
	assert.Nil(t, DoWithRetry(func() error {
		times2++
		if times2 == total {
			return nil
		}
		time.Sleep(time.Millisecond * 50)
		return errors.New("any")
	}, WithTimeout(time.Millisecond*50*(time.Duration(total)+2))))

	assert.NotNil(t, DoWithRetry(func() error {
		return errors.New("any")
	}, WithTimeout(time.Millisecond*250)))
}

func TestRetryWithInterval(t *testing.T) {
	times1 := 0
	assert.NotNil(t, DoWithRetry(func() error {
		times1++
		if times1 == 1 {
			return errors.New("any")
		}
		time.Sleep(time.Millisecond * 150)
		return nil
	}, WithTimeout(time.Millisecond*250), WithInterval(time.Millisecond*150)))

	times2 := 0
	assert.NotNil(t, DoWithRetry(func() error {
		times2++
		if times2 == 2 {
			return nil
		}
		time.Sleep(time.Millisecond * 150)
		return errors.New("any ")
	}, WithTimeout(time.Millisecond*250), WithInterval(time.Millisecond*150)))

}

func TestRetryWithWithIgnoreErrors(t *testing.T) {
	ignoreErr1 := errors.New("ignore error1")
	ignoreErr2 := errors.New("ignore error2")
	ignoreErrs := []error{ignoreErr1, ignoreErr2}

	assert.Nil(t, DoWithRetry(func() error {
		return ignoreErr1
	}, WithIgnoreErrors(ignoreErrs)))

	assert.Nil(t, DoWithRetry(func() error {
		return ignoreErr2
	}, WithIgnoreErrors(ignoreErrs)))

	assert.NotNil(t, DoWithRetry(func() error {
		return errors.New("any")
	}))
}

func TestRetryCtx(t *testing.T) {
	t.Run("with timeout", func(t *testing.T) {
		assert.NotNil(t, DoWithRetryCtx(context.Background(), func(ctx context.Context, retryCount int) error {
			if retryCount == 0 {
				return errors.New("any")
			}
			time.Sleep(time.Millisecond * 150)
			return nil
		}, WithTimeout(time.Millisecond*250), WithInterval(time.Millisecond*150)))

		assert.NotNil(t, DoWithRetryCtx(context.Background(), func(ctx context.Context, retryCount int) error {
			if retryCount == 1 {
				return nil
			}
			time.Sleep(time.Millisecond * 150)
			return errors.New("any ")
		}, WithTimeout(time.Millisecond*250), WithInterval(time.Millisecond*150)))
	})

	t.Run("with deadline exceeded", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*250))
		defer cancel()

		var times int
		assert.Error(t, DoWithRetryCtx(ctx, func(ctx context.Context, retryCount int) error {
			times++
			time.Sleep(time.Millisecond * 150)
			return errors.New("any")
		}, WithInterval(time.Millisecond*150)))
		assert.Equal(t, 1, times)
	})

	t.Run("with deadline not exceeded", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond*250))
		defer cancel()

		var times int
		assert.NoError(t, DoWithRetryCtx(ctx, func(ctx context.Context, retryCount int) error {
			times++
			if times == defaultRetryTimes {
				return nil
			}

			time.Sleep(time.Millisecond * 50)
			return errors.New("any")
		}))
		assert.Equal(t, defaultRetryTimes, times)
	})
}
