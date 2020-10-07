package cache

import (
	"errors"
	"math/rand"
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
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

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
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

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
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

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
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

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
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

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
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer s.Close()

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
