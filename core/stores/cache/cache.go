package cache

import (
	"fmt"
	"log"
	"time"

	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/core/hash"
	"github.com/tal-tech/go-zero/core/syncx"
)

type (
	Cache interface {
		Inc(key string, count *int64) error
		IncBy(key string, increment int64, count *int64) error
		Count(key string) (int64, error)
		Counts(keys ...string) ([]int64, error)

		DelCache(keys ...string) error
		GetCache(key string, v interface{}) error
		SetCache(key string, v interface{}) error
		SetCacheWithExpire(key string, v interface{}, expire time.Duration) error
		Take(v interface{}, key string, query func(v interface{}) error) error
		TakeWithExpire(v interface{}, key string, query func(v interface{}, expire time.Duration) error) error
	}

	cacheCluster struct {
		dispatcher  *hash.ConsistentHash
		errNotFound error
	}
)

func NewCache(c ClusterConf, barrier syncx.SharedCalls, st *CacheStat, errNotFound error,
	opts ...Option) Cache {
	if len(c) == 0 || TotalWeights(c) <= 0 {
		log.Fatal("no cache nodes")
	}

	if len(c) == 1 {
		return NewCacheNode(c[0].NewRedis(), barrier, st, errNotFound, opts...)
	}

	dispatcher := hash.NewConsistentHash()
	for _, node := range c {
		cn := NewCacheNode(node.NewRedis(), barrier, st, errNotFound, opts...)
		dispatcher.AddWithWeight(cn, node.Weight)
	}

	return cacheCluster{
		dispatcher:  dispatcher,
		errNotFound: errNotFound,
	}
}

func (cc cacheCluster) DelCache(keys ...string) error {
	switch len(keys) {
	case 0:
		return nil
	case 1:
		key := keys[0]
		c, ok := cc.dispatcher.Get(key)
		if !ok {
			return cc.errNotFound
		}

		return c.(Cache).DelCache(key)
	default:
		var be errorx.BatchError
		nodes := make(map[interface{}][]string)
		for _, key := range keys {
			c, ok := cc.dispatcher.Get(key)
			if !ok {
				be.Add(fmt.Errorf("key %q not found", key))
				continue
			}

			nodes[c] = append(nodes[c], key)
		}
		for c, ks := range nodes {
			if err := c.(Cache).DelCache(ks...); err != nil {
				be.Add(err)
			}
		}

		return be.Err()
	}
}

func (cc cacheCluster) GetCache(key string, v interface{}) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).GetCache(key, v)
}

func (cc cacheCluster) SetCache(key string, v interface{}) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).SetCache(key, v)
}

func (cc cacheCluster) SetCacheWithExpire(key string, v interface{}, expire time.Duration) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).SetCacheWithExpire(key, v, expire)
}

func (cc cacheCluster) Take(v interface{}, key string, query func(v interface{}) error) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).Take(v, key, query)
}

func (cc cacheCluster) TakeWithExpire(v interface{}, key string,
	query func(v interface{}, expire time.Duration) error) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).TakeWithExpire(v, key, query)
}



func (cc cacheCluster) Inc(key string, count *int64) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).Inc(key, count)
}

func (cc cacheCluster) IncBy(key string, increment int64, count *int64) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).IncBy(key, increment, count)
}

func (cc cacheCluster) Count(key string) (int64, error) {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return 0, cc.errNotFound
	}

	return c.(Cache).Count(key)
}

func (cc cacheCluster) Counts(keys ...string) ([]int64, error) {
	switch len(keys) {
	case 0:
		return nil, cc.errNotFound
	default:
		var be errorx.BatchError
		ret := make([]int64,  len(keys))
		for index, key := range keys {
			c, ok := cc.dispatcher.Get(key)
			if !ok {
				ret[index] = 0
				be.Add(fmt.Errorf("key %q not found", key))
				continue
			}
			count, err := c.(Cache).Count(key)
			if err != nil {
				ret[index] = 0
				be.Add(err)
			} else {
				ret[index] = count
			}
		}
		return ret, be.Err()
	}
}