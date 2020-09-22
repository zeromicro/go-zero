package internal

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
