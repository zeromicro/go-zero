package syncx

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/lang"
)

const limit = 10

func TestPoolGet(t *testing.T) {
	stack := NewPool(limit, create, destroy)
	ch := make(chan lang.PlaceholderType)

	for i := 0; i < limit; i++ {
		var fail AtomicBool
		go func() {
			v := stack.Get()
			if v.(int) != 1 {
				fail.Set(true)
			}
			ch <- lang.Placeholder
		}()

		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fail()
		}

		if fail.True() {
			t.Fatal("unmatch value")
		}
	}
}

func TestPoolPopTooMany(t *testing.T) {
	stack := NewPool(limit, create, destroy)
	ch := make(chan lang.PlaceholderType, 1)

	for i := 0; i < limit; i++ {
		var wait sync.WaitGroup
		wait.Add(1)
		go func() {
			stack.Get()
			ch <- lang.Placeholder
			wait.Done()
		}()

		wait.Wait()
		select {
		case <-ch:
		default:
			t.Fail()
		}
	}

	var waitGroup, pushWait sync.WaitGroup
	waitGroup.Add(1)
	pushWait.Add(1)
	go func() {
		pushWait.Done()
		stack.Get()
		waitGroup.Done()
	}()

	pushWait.Wait()
	stack.Put(1)
	waitGroup.Wait()
}

func TestPoolPopFirst(t *testing.T) {
	var value int32
	stack := NewPool(limit, func() any {
		return atomic.AddInt32(&value, 1)
	}, destroy)

	for i := 0; i < 100; i++ {
		v := stack.Get().(int32)
		assert.Equal(t, 1, int(v))
		stack.Put(v)
	}
}

func TestPoolWithMaxAge(t *testing.T) {
	var value int32
	stack := NewPool(limit, func() any {
		return atomic.AddInt32(&value, 1)
	}, destroy, WithMaxAge(time.Millisecond))

	v1 := stack.Get().(int32)
	// put nil should not matter
	stack.Put(nil)
	stack.Put(v1)
	time.Sleep(time.Millisecond * 10)
	v2 := stack.Get().(int32)
	assert.NotEqual(t, v1, v2)
}

func TestNewPoolPanics(t *testing.T) {
	assert.Panics(t, func() {
		NewPool(0, create, destroy)
	})
}

func TestPoolDestroyAll(t *testing.T) {
	var destroyed []int
	var destroyCount int32

	destroyFunc := func(item any) {
		destroyed = append(destroyed, item.(int))
		atomic.AddInt32(&destroyCount, 1)
	}

	pool := NewPool(limit, create, destroyFunc)

	// Put some resources into the pool
	pool.Put(10)
	pool.Put(20)
	pool.Put(30)

	// Destroy all resources
	pool.DestroyAll()

	// Verify all resources were destroyed
	assert.Equal(t, int32(3), atomic.LoadInt32(&destroyCount))
	assert.Contains(t, destroyed, 10)
	assert.Contains(t, destroyed, 20)
	assert.Contains(t, destroyed, 30)

	// Verify pool is empty - next Get should create new resource
	val := pool.Get()
	assert.Equal(t, 1, val) // create() returns 1
}

func TestPoolDestroyAllEmpty(t *testing.T) {
	var destroyCount int32
	destroyFunc := func(_ any) {
		atomic.AddInt32(&destroyCount, 1)
	}

	pool := NewPool(limit, create, destroyFunc)

	// DestroyAll on empty pool should not panic
	pool.DestroyAll()

	// No resources should have been destroyed
	assert.Equal(t, int32(0), atomic.LoadInt32(&destroyCount))

	// Pool should still work normally
	val := pool.Get()
	assert.Equal(t, 1, val)
}

func TestPoolDestroyAllWithNilDestroy(t *testing.T) {
	pool := NewPool(limit, create, nil)

	// Put some resources into the pool
	pool.Put(10)
	pool.Put(20)

	// DestroyAll with nil destroy function should not panic
	pool.DestroyAll()

	// Pool should be empty and work normally
	val := pool.Get()
	assert.Equal(t, 1, val)
}

func TestPoolDestroyAllConcurrency(t *testing.T) {
	var destroyCount int32
	var createCount int32

	createFunc := func() any {
		return atomic.AddInt32(&createCount, 1)
	}

	destroyFunc := func(_ any) {
		atomic.AddInt32(&destroyCount, 1)
	}

	pool := NewPool(limit, createFunc, destroyFunc)

	// Add some initial resources
	for i := 0; i < 5; i++ {
		pool.Put(i + 100)
	}

	var wg sync.WaitGroup
	const goroutines = 10

	// Concurrently perform various operations
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			switch id % 4 {
			case 0:
				// DestroyAll
				pool.DestroyAll()
			case 1:
				// Get resources
				val := pool.Get()
				pool.Put(val)
			case 2:
				// Put resources
				pool.Put(id + 1000)
			case 3:
				// Get and don't put back
				pool.Get()
			}
		}(i)
	}

	wg.Wait()

	// Final DestroyAll to clean up
	pool.DestroyAll()

	// Pool should work after concurrent operations
	val := pool.Get()
	assert.NotNil(t, val)
}

func TestPoolDestroyAllWakesWaitingGoroutines(t *testing.T) {
	pool := NewPool(1, create, destroy) // Small pool size

	// Fill the pool
	resource := pool.Get()
	assert.Equal(t, 1, resource)

	var wg sync.WaitGroup
	var gotResource bool

	// Start a goroutine that will wait for a resource
	wg.Add(1)
	go func() {
		defer wg.Done()
		val := pool.Get() // This will block since pool is full
		gotResource = true
		assert.Equal(t, 1, val) // Should get a newly created resource after DestroyAll
	}()

	// Give the goroutine time to start waiting
	time.Sleep(10 * time.Millisecond)

	// DestroyAll should wake up the waiting goroutine
	pool.DestroyAll()

	wg.Wait()
	assert.True(t, gotResource)
}

func create() any {
	return 1
}

func destroy(_ any) {
}
