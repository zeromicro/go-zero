package syncx

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestImmutableResource(t *testing.T) {
	var count int
	r := NewImmutableResource(func() (any, error) {
		count++
		return "hello", nil
	})

	res, err := r.Get()
	assert.Equal(t, "hello", res)
	assert.Equal(t, 1, count)
	assert.Nil(t, err)

	// again
	res, err = r.Get()
	assert.Equal(t, "hello", res)
	assert.Equal(t, 1, count)
	assert.Nil(t, err)
}

func TestImmutableResourceError(t *testing.T) {
	var count int
	r := NewImmutableResource(func() (any, error) {
		count++
		return nil, errors.New("any")
	})

	res, err := r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 1, count)

	// again
	res, err = r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 1, count)

	r.refreshInterval = 0
	time.Sleep(time.Millisecond)
	res, err = r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 2, count)
}

// It's hard to test more than one goroutine fetching the resource at the same time,
// because it's difficult to make more than one goroutine to pass the first read lock
// and wait another to pass the read lock before it gets the write lock.
func TestImmutableResourceConcurrent(t *testing.T) {
	const message = "hello"
	var count int32
	ready := make(chan struct{})
	r := NewImmutableResource(func() (any, error) {
		atomic.AddInt32(&count, 1)
		close(ready)                      // signal that fetch started
		time.Sleep(10 * time.Millisecond) // simulate slow fetch
		return message, nil
	})

	const goroutines = 10
	var wg sync.WaitGroup
	results := make([]any, goroutines)
	errs := make([]error, goroutines)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			res, err := r.Get()
			results[idx] = res
			errs[idx] = err
		}(i)
	}

	// wait for fetch to start
	<-ready

	wg.Wait()

	// fetch should only be called once despite concurrent access
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))

	// all goroutines should eventually get the same result
	for i := 0; i < goroutines; i++ {
		assert.Nil(t, errs[i])
		assert.Equal(t, message, results[i])
	}
}

func TestImmutableResourceErrorRefreshAlways(t *testing.T) {
	var count int
	r := NewImmutableResource(func() (any, error) {
		count++
		return nil, errors.New("any")
	}, WithRefreshIntervalOnFailure(0))

	res, err := r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 1, count)

	// again
	res, err = r.Get()
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "any", err.Error())
	assert.Equal(t, 2, count)
}
