package cache

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mathx"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/syncx"
)

var errTestNotFound = errors.New("not found")

func init() {
	logx.Disable()
	stat.SetReporter(nil)
}

func TestCacheNode_DelCache(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errTestNotFound,
	}
	assert.Nil(t, cn.DelCache())
	assert.Nil(t, cn.DelCache([]string{}...))
	assert.Nil(t, cn.DelCache(make([]string, 0)...))
	cn.SetCache("first", "one")
	assert.Nil(t, cn.DelCache("first"))
	cn.SetCache("first", "one")
	cn.SetCache("second", "two")
	assert.Nil(t, cn.DelCache("first", "second"))
}

func TestCacheNode_InvalidCache(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errTestNotFound,
	}
	s.Set("any", "value")
	var str string
	assert.NotNil(t, cn.GetCache("any", &str))
	assert.Equal(t, "", str)
	_, err = s.Get("any")
	assert.Equal(t, miniredis.ErrKeyNotFound, err)
}

func TestCacheNode_Take(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSharedCalls(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errTestNotFound,
	}
	var str string
	err = cn.Take(&str, "any", func(v interface{}) error {
		*v.(*string) = "value"
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "value", str)
	assert.Nil(t, cn.GetCache("any", &str))
	val, err := s.Get("any")
	assert.Nil(t, err)
	assert.Equal(t, `"value"`, val)
}

func TestCacheNode_TakeNotFound(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSharedCalls(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errTestNotFound,
	}
	var str string
	err = cn.Take(&str, "any", func(v interface{}) error {
		return errTestNotFound
	})
	assert.Equal(t, errTestNotFound, err)
	assert.Equal(t, errTestNotFound, cn.GetCache("any", &str))
	val, err := s.Get("any")
	assert.Nil(t, err)
	assert.Equal(t, `*`, val)

	s.Set("any", "*")
	err = cn.Take(&str, "any", func(v interface{}) error {
		return nil
	})
	assert.Equal(t, errTestNotFound, err)
	assert.Equal(t, errTestNotFound, cn.GetCache("any", &str))

	s.Del("any")
	var errDummy = errors.New("dummy")
	err = cn.Take(&str, "any", func(v interface{}) error {
		return errDummy
	})
	assert.Equal(t, errDummy, err)
}

func TestCacheNode_TakeWithExpire(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSharedCalls(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errors.New("any"),
	}
	var str string
	err = cn.TakeWithExpire(&str, "any", func(v interface{}, expire time.Duration) error {
		*v.(*string) = "value"
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "value", str)
	assert.Nil(t, cn.GetCache("any", &str))
	val, err := s.Get("any")
	assert.Nil(t, err)
	assert.Equal(t, `"value"`, val)
}

func TestCacheNode_String(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSharedCalls(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errors.New("any"),
	}
	assert.Equal(t, s.Addr(), cn.String())
}

func TestCacheValueWithBigInt(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		barrier:        syncx.NewSharedCalls(),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errors.New("any"),
	}

	const (
		key         = "key"
		value int64 = 323427211229009810
	)

	assert.Nil(t, cn.SetCache(key, value))
	var val interface{}
	assert.Nil(t, cn.GetCache(key, &val))
	assert.Equal(t, strconv.FormatInt(value, 10), fmt.Sprintf("%v", val))
}

func TestCacheNode_Inc(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	testKey := "TestCacheNode_Inc"

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errTestNotFound,
	}

	// test count
	c, e := cn.Count(testKey)
	assert.Nil(t, e)
	assert.Equal(t, c, int64(0))

	assert.Nil(t, cn.Inc(testKey, &c))
	assert.Equal(t, c, int64(1))
}

func TestCacheNode_IncBy(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	testKey := "TestCacheNode_IncBy"

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errTestNotFound,
	}

	c, e := cn.Count(testKey)
	assert.Nil(t, e)
	assert.Equal(t, c, int64(0))

	assert.Nil(t, cn.IncBy(testKey, 10, &c))
	assert.Equal(t, c, int64(10))
}

func TestCacheNode_Count(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	testKey := "TestCacheNode_Count"

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errTestNotFound,
	}

	// before
	c, e := cn.Count(testKey)
	assert.Nil(t, e)
	assert.Equal(t, c, int64(0))

	// add increment
	assert.Nil(t, cn.Inc(testKey, &c))
	assert.Equal(t, c, int64(1))

	// after
	c, e = cn.Count(testKey)
	assert.Nil(t, e)
	assert.Equal(t, c, int64(1))
}

func TestCacheNode_Counts(t *testing.T) {
	s, clean, err := createMiniRedis()
	assert.Nil(t, err)
	defer clean()

	testKeys := []string{
		"TestCacheNode_Inc",
		"TestCacheNode_IncBy",
		"TestCacheNode_Count",
		"TestCacheNode_Counts",
	}

	cn := cacheNode{
		rds:            redis.NewRedis(s.Addr(), redis.NodeType),
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           NewCacheStat("any"),
		errNotFound:    errTestNotFound,
	}

	// before
	counts, e := cn.Counts(testKeys...)
	assert.Nil(t, e)
	assert.NotNil(t, counts)
	assert.Equal(t, len(testKeys), len(counts))
	for _, r := range counts {
		assert.Equal(t, r, int64(0))
	}

	// add increment
	var c int64
	for _, key := range testKeys {
		err = cn.Inc(key, &c)
		assert.Nil(t, err)
		assert.Equal(t, c, int64(1))
	}

	// add increment
	for index, key := range testKeys {
		err = cn.IncBy(key, int64(index), &c)
		assert.Nil(t, err)
		assert.Equal(t, c, int64(1+index))
	}

	// after
	counts, e = cn.Counts(testKeys...)
	assert.Nil(t, e)
	assert.NotNil(t, counts)
	assert.Equal(t, len(testKeys), len(counts))

	for index, r := range counts {
		assert.Equal(t, r, int64(1+index))
	}
}
