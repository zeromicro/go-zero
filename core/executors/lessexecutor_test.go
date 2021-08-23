package executors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/timex"
)

func TestLessExecutor_DoOrDiscard(t *testing.T) {
	executor := NewLessExecutor(time.Minute)
	assert.True(t, executor.DoOrDiscard(func() {}))
	assert.False(t, executor.DoOrDiscard(func() {}))
	executor.lastTime.Set(timex.Now() - time.Minute - time.Second*30)
	assert.True(t, executor.DoOrDiscard(func() {}))
	assert.False(t, executor.DoOrDiscard(func() {}))
}

func BenchmarkLessExecutor(b *testing.B) {
	exec := NewLessExecutor(time.Millisecond)
	for i := 0; i < b.N; i++ {
		exec.DoOrDiscard(func() {
		})
	}
}
