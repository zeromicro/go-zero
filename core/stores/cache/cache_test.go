package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/hash"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	"github.com/zeromicro/go-zero/core/syncx"
)

var _ Cache = (*mockedNode)(nil)

type mockedNode struct {
	vals        map[string][]byte
	errNotFound error
}

func (mc *mockedNode) Del(keys ...string) error {
	return mc.DelCtx(context.Background(), keys...)
}

func (mc *mockedNode) DelCtx(_ context.Context, keys ...string) error {
	var be errorx.BatchError

	for _, key := range keys {
		if _, ok := mc.vals[key]; !ok {
			be.Add(mc.errNotFound)
		} else {
			delete(mc.vals, key)
		}
	}

	return be.Err()
}

func (mc *mockedNode) Get(key string, val any) error {
	return mc.GetCtx(context.Background(), key, val)
}

func (mc *mockedNode) GetCtx(ctx context.Context, key string, val any) error {
	bs, ok := mc.vals[key]
	if ok {
		return json.Unmarshal(bs, val)
	}

	return mc.errNotFound
}

func (mc *mockedNode) IsNotFound(err error) bool {
	return errors.Is(err, mc.errNotFound)
}

func (mc *mockedNode) Set(key string, val any) error {
	return mc.SetCtx(context.Background(), key, val)
}

func (mc *mockedNode) SetCtx(ctx context.Context, key string, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	mc.vals[key] = data
	return nil
}

func (mc *mockedNode) SetWithExpire(key string, val any, expire time.Duration) error {
	return mc.SetWithExpireCtx(context.Background(), key, val, expire)
}

func (mc *mockedNode) SetWithExpireCtx(ctx context.Context, key string, val any, expire time.Duration) error {
	return mc.Set(key, val)
}

func (mc *mockedNode) Take(val any, key string, query func(val any) error) error {
	return mc.TakeCtx(context.Background(), val, key, query)
}

func (mc *mockedNode) TakeCtx(ctx context.Context, val any, key string, query func(val any) error) error {
	if _, ok := mc.vals[key]; ok {
		return mc.GetCtx(ctx, key, val)
	}

	if err := query(val); err != nil {
		return err
	}

	return mc.SetCtx(ctx, key, val)
}

func (mc *mockedNode) TakeWithExpire(val any, key string, query func(val any, expire time.Duration) error) error {
	return mc.TakeWithExpireCtx(context.Background(), val, key, query)
}

func (mc *mockedNode) TakeWithExpireCtx(ctx context.Context, val any, key string, query func(val any, expire time.Duration) error) error {
	return mc.Take(val, key, func(val any) error {
		return query(val, 0)
	})
}

func TestCache_SetDel(t *testing.T) {
	t.Run("test set del", func(t *testing.T) {
		const total = 1000
		r1 := redistest.CreateRedis(t)
		r2 := redistest.CreateRedis(t)
		conf := ClusterConf{
			{
				RedisConf: redis.RedisConf{
					Host: r1.Addr,
					Type: redis.NodeType,
				},
				Weight: 100,
			},
			{
				RedisConf: redis.RedisConf{
					Host: r2.Addr,
					Type: redis.NodeType,
				},
				Weight: 100,
			},
		}
		c := New(conf, syncx.NewSingleFlight(), NewStat("mock"), errPlaceholder)
		for i := 0; i < total; i++ {
			if i%2 == 0 {
				assert.Nil(t, c.Set(fmt.Sprintf("key/%d", i), i))
			} else {
				assert.Nil(t, c.SetWithExpire(fmt.Sprintf("key/%d", i), i, 0))
			}
		}
		for i := 0; i < total; i++ {
			var val int
			assert.Nil(t, c.Get(fmt.Sprintf("key/%d", i), &val))
			assert.Equal(t, i, val)
		}
		assert.Nil(t, c.Del())
		for i := 0; i < total; i++ {
			assert.Nil(t, c.Del(fmt.Sprintf("key/%d", i)))
		}
		assert.Nil(t, c.Del("a", "b", "c"))
		for i := 0; i < total; i++ {
			var val int
			assert.True(t, c.IsNotFound(c.Get(fmt.Sprintf("key/%d", i), &val)))
			assert.Equal(t, 0, val)
		}
	})

	t.Run("test set del error", func(t *testing.T) {
		r1, err := miniredis.Run()
		assert.NoError(t, err)
		defer r1.Close()

		r2, err := miniredis.Run()
		assert.NoError(t, err)
		defer r2.Close()

		conf := ClusterConf{
			{
				RedisConf: redis.RedisConf{
					Host: r1.Addr(),
					Type: redis.NodeType,
				},
				Weight: 100,
			},
			{
				RedisConf: redis.RedisConf{
					Host: r2.Addr(),
					Type: redis.NodeType,
				},
				Weight: 100,
			},
		}
		c := New(conf, syncx.NewSingleFlight(), NewStat("mock"), errPlaceholder)
		r1.SetError("mock error")
		r2.SetError("mock error")
		assert.NoError(t, c.Del("a", "b", "c"))
	})
}

func TestCache_OneNode(t *testing.T) {
	const total = 1000
	r := redistest.CreateRedis(t)
	conf := ClusterConf{
		{
			RedisConf: redis.RedisConf{
				Host: r.Addr,
				Type: redis.NodeType,
			},
			Weight: 100,
		},
	}
	c := New(conf, syncx.NewSingleFlight(), NewStat("mock"), errPlaceholder)
	for i := 0; i < total; i++ {
		if i%2 == 0 {
			assert.Nil(t, c.Set(fmt.Sprintf("key/%d", i), i))
		} else {
			assert.Nil(t, c.SetWithExpire(fmt.Sprintf("key/%d", i), i, 0))
		}
	}
	for i := 0; i < total; i++ {
		var val int
		assert.Nil(t, c.Get(fmt.Sprintf("key/%d", i), &val))
		assert.Equal(t, i, val)
	}
	assert.Nil(t, c.Del())
	for i := 0; i < total; i++ {
		assert.Nil(t, c.Del(fmt.Sprintf("key/%d", i)))
	}
	for i := 0; i < total; i++ {
		var val int
		assert.True(t, c.IsNotFound(c.Get(fmt.Sprintf("key/%d", i), &val)))
		assert.Equal(t, 0, val)
	}
}

func TestCache_Balance(t *testing.T) {
	const (
		numNodes = 100
		total    = 10000
	)
	dispatcher := hash.NewConsistentHash()
	maps := make([]map[string][]byte, numNodes)
	for i := 0; i < numNodes; i++ {
		maps[i] = map[string][]byte{
			strconv.Itoa(i): []byte(strconv.Itoa(i)),
		}
	}
	for i := 0; i < numNodes; i++ {
		dispatcher.AddWithWeight(&mockedNode{
			vals:        maps[i],
			errNotFound: errPlaceholder,
		}, 100)
	}

	c := cacheCluster{
		dispatcher:  dispatcher,
		errNotFound: errPlaceholder,
	}
	for i := 0; i < total; i++ {
		assert.Nil(t, c.Set(strconv.Itoa(i), i))
	}

	counts := make(map[int]int)
	for i, m := range maps {
		counts[i] = len(m)
	}
	entropy := calcEntropy(counts, total)
	assert.True(t, len(counts) > 1)
	assert.True(t, entropy > .95, fmt.Sprintf("entropy should be greater than 0.95, but got %.2f", entropy))

	for i := 0; i < total; i++ {
		var val int
		assert.Nil(t, c.Get(strconv.Itoa(i), &val))
		assert.Equal(t, i, val)
	}

	for i := 0; i < total/10; i++ {
		assert.Nil(t, c.Del(strconv.Itoa(i*10), strconv.Itoa(i*10+1), strconv.Itoa(i*10+2)))
		assert.Nil(t, c.Del(strconv.Itoa(i*10+9)))
	}

	var count int
	for i := 0; i < total/10; i++ {
		var val int
		if i%2 == 0 {
			assert.Nil(t, c.Take(&val, strconv.Itoa(i*10), func(val any) error {
				*val.(*int) = i
				count++
				return nil
			}))
		} else {
			assert.Nil(t, c.TakeWithExpire(&val, strconv.Itoa(i*10), func(val any, expire time.Duration) error {
				*val.(*int) = i
				count++
				return nil
			}))
		}
		assert.Equal(t, i, val)
	}
	assert.Equal(t, total/10, count)
}

func TestCacheNoNode(t *testing.T) {
	dispatcher := hash.NewConsistentHash()
	c := cacheCluster{
		dispatcher:  dispatcher,
		errNotFound: errPlaceholder,
	}
	assert.NotNil(t, c.Del("foo"))
	assert.NotNil(t, c.Del("foo", "bar", "any"))
	assert.NotNil(t, c.Get("foo", nil))
	assert.NotNil(t, c.Set("foo", nil))
	assert.NotNil(t, c.SetWithExpire("foo", nil, time.Second))
	assert.NotNil(t, c.Take(nil, "foo", func(val any) error {
		return nil
	}))
	assert.NotNil(t, c.TakeWithExpire(nil, "foo", func(val any, duration time.Duration) error {
		return nil
	}))
}

func calcEntropy(m map[int]int, total int) float64 {
	var entropy float64

	for _, val := range m {
		proba := float64(val) / float64(total)
		entropy -= proba * math.Log2(proba)
	}

	return entropy / math.Log2(float64(len(m)))
}
