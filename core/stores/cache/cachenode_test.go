package cache

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
)

var errTestNotFound = errors.New("not found")

func init() {
	logx.Disable()
	stat.SetReporter(nil)
}

func TestCacheNode_DelCache(t *testing.T) {
	t.Run("del cache", func(t *testing.T) {
		store, clean, err := redistest.CreateRedis()
		assert.Nil(t, err)
		store.Type = redis.ClusterType
		defer clean()

		cn := cacheNode{
			rds:            store,
			r:              rand.New(rand.NewSource(time.Now().UnixNano())),
			lock:           new(sync.Mutex),
			unstableExpiry: mathx.NewUnstable(expiryDeviation),
			stat:           NewStat("any"),
			errNotFound:    errTestNotFound,
		}
		assert.Nil(t, cn.Del())
		assert.Nil(t, cn.Del([]string{}...))
		assert.Nil(t, cn.Del(make([]string, 0)...))
		cn.Set("first", "one")
		assert.Nil(t, cn.Del("first"))
		cn.Set("first", "one")
		cn.Set("second", "two")
		assert.Nil(t, cn.Del("first", "second"))
	})

	t.Run("del cache with errors", func(t *testing.T) {
		old := timingWheel
		ticker := timex.NewFakeTicker()
		var err error
		timingWheel, err = collection.NewTimingWheelWithTicker(
			time.Millisecond, timingWheelSlots, func(key, value interface{}) {
				clean(key, value)
			}, ticker)
		assert.NoError(t, err)
		t.Cleanup(func() {
			timingWheel = old
		})

		r, err := miniredis.Run()
		assert.NoError(t, err)
		defer r.Close()
		r.SetError("mock error")

		node := NewNode(redis.New(r.Addr(), redis.Cluster()), syncx.NewSingleFlight(),
			NewStat("any"), errTestNotFound)
		assert.NoError(t, node.Del("foo", "bar"))
		ticker.Tick()
		runtime.Gosched()
	})
}

func TestCacheNode_DelCacheWithErrors(t *testing.T) {
	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()
	store.Type = redis.ClusterType

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errTestNotFound,
	}
	assert.Nil(t, cn.Del("third", "fourth"))
}

func TestCacheNode_InvalidCache(t *testing.T) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

	cn := cacheNode{
		rds:            redis.New(s.Addr()),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errTestNotFound,
	}
	s.Set("any", "value")
	var str string
	assert.NotNil(t, cn.Get("any", &str))
	assert.Equal(t, "", str)
	_, err = s.Get("any")
	assert.Equal(t, miniredis.ErrKeyNotFound, err)
}

func TestCacheNode_SetWithExpire(t *testing.T) {
	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errors.New("any"),
	}
	assert.NotNil(t, cn.SetWithExpire("key", make(chan int), time.Second))
}

func TestCacheNode_Take(t *testing.T) {
	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cn := NewNode(store, syncx.NewSingleFlight(), NewStat("any"), errTestNotFound,
		WithExpiry(time.Second), WithNotFoundExpiry(time.Second))
	var str string
	err = cn.Take(&str, "any", func(v interface{}) error {
		*v.(*string) = "value"
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "value", str)
	assert.Nil(t, cn.Get("any", &str))
	val, err := store.Get("any")
	assert.Nil(t, err)
	assert.Equal(t, `"value"`, val)
}

func TestCacheNode_TakeBadRedis(t *testing.T) {
	r, err := miniredis.Run()
	assert.NoError(t, err)
	defer r.Close()
	r.SetError("mock error")

	cn := NewNode(redis.New(r.Addr()), syncx.NewSingleFlight(), NewStat("any"),
		errTestNotFound, WithExpiry(time.Second), WithNotFoundExpiry(time.Second))
	var str string
	assert.Error(t, cn.Take(&str, "any", func(v interface{}) error {
		*v.(*string) = "value"
		return nil
	}))
}

func TestCacheNode_TakeNotFound(t *testing.T) {
	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errTestNotFound,
	}
	var str string
	err = cn.Take(&str, "any", func(v interface{}) error {
		return errTestNotFound
	})
	assert.True(t, cn.IsNotFound(err))
	assert.True(t, cn.IsNotFound(cn.Get("any", &str)))
	val, err := store.Get("any")
	assert.Nil(t, err)
	assert.Equal(t, `*`, val)

	store.Set("any", "*")
	err = cn.Take(&str, "any", func(v interface{}) error {
		return nil
	})
	assert.True(t, cn.IsNotFound(err))
	assert.True(t, cn.IsNotFound(cn.Get("any", &str)))

	store.Del("any")
	errDummy := errors.New("dummy")
	err = cn.Take(&str, "any", func(v interface{}) error {
		return errDummy
	})
	assert.Equal(t, errDummy, err)
}

func TestCacheNode_TakeWithExpire(t *testing.T) {
	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errors.New("any"),
	}
	var str string
	err = cn.TakeWithExpire(&str, "any", func(v interface{}, expire time.Duration) error {
		*v.(*string) = "value"
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "value", str)
	assert.Nil(t, cn.Get("any", &str))
	val, err := store.Get("any")
	assert.Nil(t, err)
	assert.Equal(t, `"value"`, val)
}

func TestCacheNode_String(t *testing.T) {
	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errors.New("any"),
	}
	assert.Equal(t, store.Addr, cn.String())
}

func TestCacheValueWithBigInt(t *testing.T) {
	store, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            store,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSingleFlight(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewStat("any"),
		errNotFound:    errors.New("any"),
	}

	const (
		key         = "key"
		value int64 = 323427211229009810
	)

	assert.Nil(t, cn.Set(key, value))
	var val interface{}
	assert.Nil(t, cn.Get(key, &val))
	assert.Equal(t, strconv.FormatInt(value, 10), fmt.Sprintf("%v", val))
}
