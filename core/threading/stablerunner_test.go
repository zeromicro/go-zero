package threading

import (
	"math/rand"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStableRunner(t *testing.T) {
	size := bufSize * 2
	rand.NewSource(time.Now().UnixNano())
	runner := NewStableRunner(func(v int) float64 {
		if v == 0 {
			time.Sleep(time.Millisecond * 100)
		} else {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
		}
		return float64(v) + 0.5
	})

	var waitGroup sync.WaitGroup
	waitGroup.Add(1)
	go func() {
		for i := 0; i < size; i++ {
			assert.NoError(t, runner.Push(i))
		}
		runner.Wait()
		waitGroup.Done()
	}()

	values := make([]float64, size)
	for i := 0; i < size; i++ {
		var err error
		values[i], err = runner.Get()
		assert.NoError(t, err)
		time.Sleep(time.Millisecond)
	}

	assert.True(t, sort.Float64sAreSorted(values))
	waitGroup.Wait()

	assert.Equal(t, ErrRunnerClosed, runner.Push(1))
	_, err := runner.Get()
	assert.Equal(t, ErrRunnerClosed, err)
}

func FuzzStableRunner(f *testing.F) {
	rand.NewSource(time.Now().UnixNano())
	f.Add(uint64(bufSize))
	f.Fuzz(func(t *testing.T, n uint64) {
		runner := NewStableRunner(func(v int) float64 {
			if v == 0 {
				time.Sleep(time.Millisecond * 100)
			} else {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
			}
			return float64(v) + 0.5
		})

		go func() {
			for i := 0; i < int(n); i++ {
				assert.NoError(t, runner.Push(i))
			}
		}()

		values := make([]float64, n)
		for i := 0; i < int(n); i++ {
			var err error
			values[i], err = runner.Get()
			assert.NoError(t, err)
		}

		runner.Wait()
		assert.True(t, sort.Float64sAreSorted(values))

		// make sure returning errors after runner is closed
		assert.Equal(t, ErrRunnerClosed, runner.Push(1))
		_, err := runner.Get()
		assert.Equal(t, ErrRunnerClosed, err)
	})
}

func BenchmarkStableRunner(b *testing.B) {
	runner := NewStableRunner(func(v int) float64 {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)))
		return float64(v) + 0.5
	})

	for i := 0; i < b.N; i++ {
		_ = runner.Push(i)
		_, _ = runner.Get()
	}
}
