package threading

import (
	"runtime"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskRunner_Schedule(t *testing.T) {
	times := 100
	pool := NewTaskRunner(runtime.NumCPU())

	var counter int32
	for i := 0; i < times; i++ {
		pool.Schedule(func() {
			atomic.AddInt32(&counter, 1)
		})
	}

	pool.Wait()

	assert.Equal(t, times, int(counter))
}

func TestTaskRunner_ScheduleImmediately(t *testing.T) {
	cpus := runtime.NumCPU()
	times := cpus * 2
	pool := NewTaskRunner(cpus)

	var counter int32
	for i := 0; i < times; i++ {
		err := pool.ScheduleImmediately(func() {
			atomic.AddInt32(&counter, 1)
			time.Sleep(time.Millisecond * 100)
		})
		if i < cpus {
			assert.Nil(t, err)
		} else {
			assert.ErrorIs(t, err, ErrTaskRunnerBusy)
		}
	}

	pool.Wait()

	assert.Equal(t, cpus, int(counter))
}

func BenchmarkRoutinePool(b *testing.B) {
	queue := NewTaskRunner(runtime.NumCPU())
	for i := 0; i < b.N; i++ {
		queue.Schedule(func() {
		})
	}
}
