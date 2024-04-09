package threading

import (
	"errors"
	"io"
	"log"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutineGroupRun(t *testing.T) {
	var count int32
	group := NewRoutineGroup()
	for i := 0; i < 3; i++ {
		group.Run(func() {
			atomic.AddInt32(&count, 1)
		})
	}

	group.Wait()

	assert.Equal(t, int32(3), count)
}

func TestRoutingGroupRunSafe(t *testing.T) {
	log.SetOutput(io.Discard)

	var count int32
	group := NewRoutineGroup()
	var once sync.Once
	for i := 0; i < 3; i++ {
		group.RunSafe(func() {
			once.Do(func() {
				panic("")
			})
			atomic.AddInt32(&count, 1)
		})
	}

	group.Wait()

	assert.Equal(t, int32(2), count)
}

func TestRoutineErrGroupRunErrorNil(t *testing.T) {
	var count int32
	group := NewRoutineErrGroup()
	group.SetLimit(3)

	for i := 0; i < 3; i++ {
		group.Run(func() error {
			atomic.AddInt32(&count, 1)
			return nil
		})
	}

	err := group.Wait()

	assert.Equal(t, int32(3), count)
	assert.NoError(t, err)
}

func TestRoutineErrGroupRun(t *testing.T) {
	var count int32
	group := NewRoutineErrGroup()
	group.SetLimit(3)

	for i := 0; i < 3; i++ {
		group.Run(func() error {
			atomic.AddInt32(&count, 1)
			return errors.New("error")
		})
	}

	err := group.Wait()

	assert.Equal(t, int32(3), count)
	assert.Error(t, err)
	assert.EqualError(t, err, "error\nerror\nerror")
}

func TestRoutingErrGroupRunSafe(t *testing.T) {
	log.SetOutput(io.Discard)

	var count int32
	group := NewRoutineErrGroup()
	group.SetLimit(3)

	var once sync.Once
	for i := 0; i < 3; i++ {
		i := i
		group.RunSafe(func() error {
			once.Do(func() {
				panic("")
			})
			atomic.AddInt32(&count, 1)
			if i == 1 {
				return errors.New("error")
			}
			return nil
		})
	}

	err := group.Wait()

	assert.Equal(t, int32(2), count)
	assert.Error(t, err)
	assert.EqualError(t, err, "error")
}
