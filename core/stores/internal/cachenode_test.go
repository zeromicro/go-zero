package internal

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
)

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
		errNotFound:    errors.New("any"),
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
		errNotFound:    errors.New("any"),
	}
	s.Set("any", "value")
	var str string
	assert.NotNil(t, cn.GetCache("any", &str))
	assert.Equal(t, "", str)
	_, err = s.Get("any")
	assert.Equal(t, miniredis.ErrKeyNotFound, err)
}
