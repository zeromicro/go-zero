package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/core/hash"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/syncx"
)

type mockedNode struct {
	vals        map[string][]byte
	errNotFound error
}

func (mc *mockedNode) DelCache(keys ...string) error {
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

func (mc *mockedNode) GetCache(key string, v interface{}) error {
	bs, ok := mc.vals[key]
	if ok {
		return json.Unmarshal(bs, v)
	}

	return mc.errNotFound
}

func (mc *mockedNode) SetCache(key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	mc.vals[key] = data
	return nil
}

func (mc *mockedNode) SetCacheWithExpire(key string, v interface{}, expire time.Duration) error {
	return mc.SetCache(key, v)
}

func (mc *mockedNode) Take(v interface{}, key string, query func(v interface{}) error) error {
	if _, ok := mc.vals[key]; ok {
		return mc.GetCache(key, v)
	}

	if err := query(v); err != nil {
		return err
	}

	return mc.SetCache(key, v)
}

func (mc *mockedNode) TakeWithExpire(v interface{}, key string, query func(v interface{}, expire time.Duration) error) error {
	return mc.Take(v, key, func(v interface{}) error {
		return query(v, 0)
	})
}

func TestCache_SetDel(t *testing.T) {
	const total = 1000
	r1 := miniredis.NewMiniRedis()
	assert.Nil(t, r1.Start())
	defer r1.Close()
	r2 := miniredis.NewMiniRedis()
	assert.Nil(t, r2.Start())
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
	c := NewCache(conf, syncx.NewSharedCalls(), NewCacheStat("mock"), errPlaceholder)
	for i := 0; i < total; i++ {
		if i%2 == 0 {
			assert.Nil(t, c.SetCache(fmt.Sprintf("key/%d", i), i))
		} else {
			assert.Nil(t, c.SetCacheWithExpire(fmt.Sprintf("key/%d", i), i, 0))
		}
	}
	for i := 0; i < total; i++ {
		var v int
		assert.Nil(t, c.GetCache(fmt.Sprintf("key/%d", i), &v))
		assert.Equal(t, i, v)
	}
	for i := 0; i < total; i++ {
		assert.Nil(t, c.DelCache(fmt.Sprintf("key/%d", i)))
	}
	for i := 0; i < total; i++ {
		var v int
		assert.Equal(t, errPlaceholder, c.GetCache(fmt.Sprintf("key/%d", i), &v))
		assert.Equal(t, 0, v)
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
		assert.Nil(t, c.SetCache(strconv.Itoa(i), i))
	}

	counts := make(map[int]int)
	for i, m := range maps {
		counts[i] = len(m)
	}
	entropy := calcEntropy(counts, total)
	assert.True(t, len(counts) > 1)
	assert.True(t, entropy > .95, fmt.Sprintf("entropy should be greater than 0.95, but got %.2f", entropy))

	for i := 0; i < total; i++ {
		var v int
		assert.Nil(t, c.GetCache(strconv.Itoa(i), &v))
		assert.Equal(t, i, v)
	}

	for i := 0; i < total/10; i++ {
		assert.Nil(t, c.DelCache(strconv.Itoa(i*10), strconv.Itoa(i*10+1), strconv.Itoa(i*10+2)))
		assert.Nil(t, c.DelCache(strconv.Itoa(i*10+9)))
	}

	var count int
	for i := 0; i < total/10; i++ {
		var val int
		if i%2 == 0 {
			assert.Nil(t, c.Take(&val, strconv.Itoa(i*10), func(v interface{}) error {
				*v.(*int) = i
				count++
				return nil
			}))
		} else {
			assert.Nil(t, c.TakeWithExpire(&val, strconv.Itoa(i*10), func(v interface{}, expire time.Duration) error {
				*v.(*int) = i
				count++
				return nil
			}))
		}
		assert.Equal(t, i, val)
	}
	assert.Equal(t, total/10, count)
}

func calcEntropy(m map[int]int, total int) float64 {
	var entropy float64

	for _, v := range m {
		proba := float64(v) / float64(total)
		entropy -= proba * math.Log2(proba)
	}

	return entropy / math.Log2(float64(len(m)))
}
