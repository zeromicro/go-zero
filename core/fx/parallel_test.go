package fx

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallel(t *testing.T) {
	var count int32
	Parallel(func() {
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&count, 1)
	}, func() {
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&count, 2)
	}, func() {
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&count, 3)
	})
	assert.Equal(t, int32(6), count)
}

func TestParallelErr(t *testing.T) {
	var count int32
	err := ParallelErr(
		func() error {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, 1)
			return errors.New("failed to exec #1")
		},
		func() error {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, 2)
			return errors.New("failed to exec #2")

		},
		func() error {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, 3)
			return nil
		},
	)

	assert.Equal(t, int32(6), count)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to exec #1", "failed to exec #2")
}

func TestParallelErrErrorNil(t *testing.T) {
	var count int32
	err := ParallelErr(
		func() error {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, 1)
			return nil
		},
		func() error {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, 2)
			return nil

		},
		func() error {
			time.Sleep(time.Millisecond * 100)
			atomic.AddInt32(&count, 3)
			return nil
		},
	)

	assert.Equal(t, int32(6), count)
	assert.NoError(t, err)
}
