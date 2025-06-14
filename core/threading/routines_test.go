package threading

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/lang"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestRoutineId(t *testing.T) {
	assert.True(t, RoutineId() > 0)
}

func TestRunSafe(t *testing.T) {
	log.SetOutput(io.Discard)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	ch := make(chan lang.PlaceholderType)
	go RunSafe(func() {
		defer func() {
			ch <- lang.Placeholder
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestRunSafeCtx(t *testing.T) {
	var buf bytes.Buffer
	logx.SetWriter(logx.NewWriter(&buf))
	ctx := context.Background()
	ch := make(chan lang.PlaceholderType)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	go RunSafeCtx(ctx, func() {
		defer func() {
			ch <- lang.Placeholder
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestGoSafeCtx(t *testing.T) {
	var buf bytes.Buffer
	logx.SetWriter(logx.NewWriter(&buf))
	ctx := context.Background()
	ch := make(chan lang.PlaceholderType)

	i := 0

	defer func() {
		assert.Equal(t, 1, i)
	}()

	GoSafeCtx(ctx, func() {
		defer func() {
			ch <- lang.Placeholder
		}()

		panic("panic")
	})

	<-ch
	i++
}

func TestRunSafeWrap(t *testing.T) {
	logx.Disable()

	t.Run("normal error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			err := errors.New("test err")
			err2 := RunSafeWrap(func() error {
				return err
			})
			assert.Equal(t, err, err2)
		})
	})

	t.Run("no error", func(t *testing.T) {
		assert.NotPanics(t, func() {
			err2 := RunSafeWrap(func() error {
				return nil
			})
			assert.Nil(t, err2)
		})
	})

	t.Run("panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			err2 := RunSafeWrap(func() error {
				panic("test")
			})
			assert.Contains(t, err2.Error(), "panic: test")
		})
	})
}
