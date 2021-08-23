package threading

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/lang"
)

func TestWorkerGroup(t *testing.T) {
	m := make(map[string]lang.PlaceholderType)
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())
	group := NewWorkerGroup(func() {
		lock.Lock()
		m[fmt.Sprint(RoutineId())] = lang.Placeholder
		lock.Unlock()
		wg.Done()
	}, runtime.NumCPU())
	go group.Start()
	wg.Wait()
	assert.Equal(t, runtime.NumCPU(), len(m))
}
