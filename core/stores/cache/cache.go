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
	// Cache interface is used to define the cache implementation.
	Cache interface {
		Del(keys ...string) error
		Get(key string, v interface{}) error
		IsNotFound(err error) bool
		Set(key string, v interface{}) error
		SetWithExpire(key string, v interface{}, expire time.Duration) error
		Take(v interface{}, key string, query func(v interface{}) error) error
		TakeWithExpire(v interface{}, key string, query func(v interface{}, expire time.Duration) error) error
	}

	cacheCluster struct {
		dispatcher  *hash.ConsistentHash
		errNotFound error
	}
)

// New returns a Cache.
func New(c ClusterConf, barrier syncx.SingleFlight, st *Stat, errNotFound error,
	opts ...Option) Cache {
	if len(c) == 0 || TotalWeights(c) <= 0 {
		log.Fatal("no cache nodes")
	}

	if len(c) == 1 {
		return NewNode(c[0].NewRedis(), barrier, st, errNotFound, opts...)
	}

	dispatcher := hash.NewConsistentHash()
	for _, node := range c {
		cn := NewNode(node.NewRedis(), barrier, st, errNotFound, opts...)
		dispatcher.AddWithWeight(cn, node.Weight)
	}

	return cacheCluster{
		dispatcher:  dispatcher,
		errNotFound: errNotFound,
	}
}

func (cc cacheCluster) Del(keys ...string) error {
	switch len(keys) {
	case 0:
		return nil
	case 1:
		key := keys[0]
		c, ok := cc.dispatcher.Get(key)
		if !ok {
			return cc.errNotFound
		}

		return c.(Cache).Del(key)
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
			if err := c.(Cache).Del(ks...); err != nil {
				be.Add(err)
			}
		}

		return be.Err()
	}
}

func (cc cacheCluster) Get(key string, v interface{}) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).Get(key, v)
}

func (cc cacheCluster) IsNotFound(err error) bool {
	return err == cc.errNotFound
}

func (cc cacheCluster) Set(key string, v interface{}) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).Set(key, v)
}

func (cc cacheCluster) SetWithExpire(key string, v interface{}, expire time.Duration) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).SetWithExpire(key, v, expire)
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
