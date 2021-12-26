package fx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithPanic(t *testing.T) {
	assert.Panics(t, func() {
		_ = DoWithTimeout(func() error {
			panic("hello")
		}, time.Millisecond*50)
	})
}

func TestWithTimeout(t *testing.T) {
	assert.Equal(t, ErrTimeout, DoWithTimeout(func() error {
		time.Sleep(time.Millisecond * 50)
		return nil
	}, time.Millisecond))
}

func TestWithoutTimeout(t *testing.T) {
	assert.Nil(t, DoWithTimeout(func() error {
		return nil
	}, time.Millisecond*50))
}

func TestWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Millisecond * 10)
		cancel()
	}()
	err := DoWithTimeout(func() error {
		time.Sleep(time.Minute)
		return nil
	}, time.Second, WithContext(ctx))
	assert.Equal(t, ErrCanceled, err)
}

func TestDoWithTimeout(t *testing.T) {

	t.Run("cancel", func(t *testing.T) {
		ctx, cancelFunc := context.WithCancel(context.Background())
		cancelFunc()
		err := DoWithTimeout(func() error {
			return nil
		}, time.Minute, WithContext(ctx))
		assert.ErrorIs(t, context.Canceled, err)
	})

	t.Run("deadlineExceeded", func(t *testing.T) {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancelFunc()
		time.Sleep(time.Second * 2)
		err := DoWithTimeout(func() error {
			return nil
		}, time.Minute, WithContext(ctx))
		assert.ErrorIs(t, context.DeadlineExceeded, err)
	})

}
