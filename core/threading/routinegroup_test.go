package threading

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx/logtest"
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
	logtest.Discard(t)

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
