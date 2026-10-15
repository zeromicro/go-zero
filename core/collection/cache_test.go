package collection

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var errDummy = errors.New("dummy")

func TestCacheSet(t *testing.T) {
	cache, err := NewCache(time.Second*2, WithName("any"))
	assert.Nil(t, err)

	cache.Set("first", "first element")
	cache.SetWithExpire("second", "second element", time.Second*3)

	value, ok := cache.Get("first")
	assert.True(t, ok)
	assert.Equal(t, "first element", value)
	value, ok = cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, "second element", value)
}

func TestCacheDel(t *testing.T) {
	cache, err := NewCache(time.Second * 2)
	assert.Nil(t, err)

	cache.Set("first", "first element")
	cache.Set("second", "second element")
	cache.Del("first")

	_, ok := cache.Get("first")
	assert.False(t, ok)
	value, ok := cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, "second element", value)
}

func TestCacheTake(t *testing.T) {
	cache, err := NewCache(time.Second * 2)
	assert.Nil(t, err)

	var count int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			cache.Take("first", func() (any, error) {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Millisecond * 100)
				return "first element", nil
			})
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, cache.size())
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}

func TestCacheTakeExists(t *testing.T) {
	cache, err := NewCache(time.Second * 2)
	assert.Nil(t, err)

	var count int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			cache.Set("first", "first element")
			cache.Take("first", func() (any, error) {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Millisecond * 100)
				return "first element", nil
			})
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, cache.size())
	assert.Equal(t, int32(0), atomic.LoadInt32(&count))
}

func TestCacheTakeError(t *testing.T) {
	cache, err := NewCache(time.Second * 2)
	assert.Nil(t, err)

	var count int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			_, err := cache.Take("first", func() (any, error) {
				atomic.AddInt32(&count, 1)
				time.Sleep(time.Millisecond * 100)
				return "", errDummy
			})
			assert.Equal(t, errDummy, err)
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, 0, cache.size())
	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}

func TestCacheWithLruEvicts(t *testing.T) {
	cache, err := NewCache(time.Minute, WithLimit(3))
	assert.Nil(t, err)

	cache.Set("first", "first element")
	cache.Set("second", "second element")
	cache.Set("third", "third element")
	cache.Set("fourth", "fourth element")

	_, ok := cache.Get("first")
	assert.False(t, ok)
	value, ok := cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, "second element", value)
	value, ok = cache.Get("third")
	assert.True(t, ok)
	assert.Equal(t, "third element", value)
	value, ok = cache.Get("fourth")
	assert.True(t, ok)
	assert.Equal(t, "fourth element", value)
}

func TestCacheWithLruEvicted(t *testing.T) {
	cache, err := NewCache(time.Minute, WithLimit(3))
	assert.Nil(t, err)

	cache.Set("first", "first element")
	cache.Set("second", "second element")
	cache.Set("third", "third element")
	cache.Set("fourth", "fourth element")

	_, ok := cache.Get("first")
	assert.False(t, ok)
	value, ok := cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, "second element", value)
	cache.Set("fifth", "fifth element")
	cache.Set("sixth", "sixth element")
	_, ok = cache.Get("third")
	assert.False(t, ok)
	_, ok = cache.Get("fourth")
	assert.False(t, ok)
	value, ok = cache.Get("second")
	assert.True(t, ok)
	assert.Equal(t, "second element", value)
}

func BenchmarkCache(b *testing.B) {
	cache, err := NewCache(time.Second*5, WithLimit(100000))
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < 10000; i++ {
		for j := 0; j < 10; j++ {
			index := strconv.Itoa(i*10000 + j)
			cache.Set("key:"+index, "value:"+index)
		}
	}

	time.Sleep(time.Second * 5)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < b.N; i++ {
				index := strconv.Itoa(i % 10000)
				cache.Get("key:" + index)
				if i%100 == 0 {
					cache.Set("key1:"+index, "value1:"+index)
				}
			}
		}
	})
}

func TestCacheExpiryConcurrentWrites(t *testing.T) {
	c, err := NewCache(time.Second)
	assert.Nil(t, err)
	t.Cleanup(c.timingWheel.Stop)

	const count = 2000
	keys := make([]string, count)
	for i := range keys {
		keys[i] = "key-" + strconv.Itoa(i)
	}

	var wg sync.WaitGroup
	deadline := time.Now().Add(3 * time.Second)
	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(offset int) {
			defer wg.Done()
			for i := offset; time.Now().Before(deadline); i++ {
				c.Set(keys[i%count], 1)
			}
		}(worker * 251)
	}
	wg.Wait()

	time.Sleep(3 * time.Second)
	var survivors []string
	for _, key := range keys {
		if _, ok := c.Get(key); ok {
			survivors = append(survivors, key)
		}
	}
	assert.Empty(t, survivors, "cache retained entries beyond TTL")

	time.Sleep(2 * time.Second)
	orphans, scheduled := 0, 0
	for _, key := range keys {
		if _, ok := c.Get(key); ok {
			if _, hasTimer := c.timingWheel.timers.Get(key); hasTimer {
				scheduled++
			} else {
				orphans++
			}
		}
	}
	assert.Equal(t, 0, orphans, "cache retained entries with no expiry timer")
	assert.Equal(t, 0, scheduled, "cache retained entries beyond TTL")
}

func TestCacheExpiryGet(t *testing.T) {
	c, err := NewCache(time.Millisecond * 100)
	assert.Nil(t, err)
	t.Cleanup(c.timingWheel.Stop)

	c.Set("k1", "v1")
	val, ok := c.Get("k1")
	assert.True(t, ok)
	assert.Equal(t, "v1", val)

	time.Sleep(time.Millisecond * 200)
	_, ok = c.Get("k1")
	assert.False(t, ok)
	assert.Equal(t, 0, c.size())
}

