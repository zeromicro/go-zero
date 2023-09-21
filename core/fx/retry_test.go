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

func TestRetryCtx(t *testing.T) {
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
}
