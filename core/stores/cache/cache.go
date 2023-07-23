package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/hash"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
)

type (
	// Cache interface is used to define the cache implementation.
	Cache interface {
		// Del deletes cached values with keys.
		Del(keys ...string) error
		// DelCtx deletes cached values with keys.
		DelCtx(ctx context.Context, keys ...string) error
		// Get gets the cache with key and fills into v.
		Get(key string, val any) error
		// GetCtx gets the cache with key and fills into v.
		GetCtx(ctx context.Context, key string, val any) error
		// IsNotFound checks if the given error is the defined errNotFound.
		IsNotFound(err error) bool
		// Set sets the cache with key and v, using c.expiry.
		Set(key string, val any) error
		// SetCtx sets the cache with key and v, using c.expiry.
		SetCtx(ctx context.Context, key string, val any) error
		// SetWithExpire sets the cache with key and v, using given expire.
		SetWithExpire(key string, val any, expire time.Duration) error
		// SetWithExpireCtx sets the cache with key and v, using given expire.
		SetWithExpireCtx(ctx context.Context, key string, val any, expire time.Duration) error
		// Take takes the result from cache first, if not found,
		// query from DB and set cache using c.expiry, then return the result.
		Take(val any, key string, query func(val any) error) error
		// TakeCtx takes the result from cache first, if not found,
		// query from DB and set cache using c.expiry, then return the result.
		TakeCtx(ctx context.Context, val any, key string, query func(val any) error) error
		// TakeWithExpire takes the result from cache first, if not found,
		// query from DB and set cache using given expire, then return the result.
		TakeWithExpire(val any, key string, query func(val any, expire time.Duration) error) error
		// TakeWithExpireCtx takes the result from cache first, if not found,
		// query from DB and set cache using given expire, then return the result.
		TakeWithExpireCtx(ctx context.Context, val any, key string,
			query func(val any, expire time.Duration) error) error
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
		return NewNode(redis.MustNewRedis(c[0].RedisConf), barrier, st, errNotFound, opts...)
	}

	dispatcher := hash.NewConsistentHash()
	for _, node := range c {
		cn := NewNode(redis.MustNewRedis(node.RedisConf), barrier, st, errNotFound, opts...)
		dispatcher.AddWithWeight(cn, node.Weight)
	}

	return cacheCluster{
		dispatcher:  dispatcher,
		errNotFound: errNotFound,
	}
}

// Del deletes cached values with keys.
func (cc cacheCluster) Del(keys ...string) error {
	return cc.DelCtx(context.Background(), keys...)
}

// DelCtx deletes cached values with keys.
func (cc cacheCluster) DelCtx(ctx context.Context, keys ...string) error {
	switch len(keys) {
	case 0:
		return nil
	case 1:
		key := keys[0]
		c, ok := cc.dispatcher.Get(key)
		if !ok {
			return cc.errNotFound
		}

		return c.(Cache).DelCtx(ctx, key)
	default:
		var be errorx.BatchError
		nodes := make(map[any][]string)
		for _, key := range keys {
			c, ok := cc.dispatcher.Get(key)
			if !ok {
				be.Add(fmt.Errorf("key %q not found", key))
				continue
			}

			nodes[c] = append(nodes[c], key)
		}
		for c, ks := range nodes {
			if err := c.(Cache).DelCtx(ctx, ks...); err != nil {
				be.Add(err)
			}
		}

		return be.Err()
	}
}

// Get gets the cache with key and fills into v.
func (cc cacheCluster) Get(key string, val any) error {
	return cc.GetCtx(context.Background(), key, val)
}

// GetCtx gets the cache with key and fills into v.
func (cc cacheCluster) GetCtx(ctx context.Context, key string, val any) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).GetCtx(ctx, key, val)
}

// IsNotFound checks if the given error is the defined errNotFound.
func (cc cacheCluster) IsNotFound(err error) bool {
	return errors.Is(err, cc.errNotFound)
}

// Set sets the cache with key and v, using c.expiry.
func (cc cacheCluster) Set(key string, val any) error {
	return cc.SetCtx(context.Background(), key, val)
}

// SetCtx sets the cache with key and v, using c.expiry.
func (cc cacheCluster) SetCtx(ctx context.Context, key string, val any) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).SetCtx(ctx, key, val)
}

// SetWithExpire sets the cache with key and v, using given expire.
func (cc cacheCluster) SetWithExpire(key string, val any, expire time.Duration) error {
	return cc.SetWithExpireCtx(context.Background(), key, val, expire)
}

// SetWithExpireCtx sets the cache with key and v, using given expire.
func (cc cacheCluster) SetWithExpireCtx(ctx context.Context, key string, val any, expire time.Duration) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).SetWithExpireCtx(ctx, key, val, expire)
}

// Take takes the result from cache first, if not found,
// query from DB and set cache using c.expiry, then return the result.
func (cc cacheCluster) Take(val any, key string, query func(val any) error) error {
	return cc.TakeCtx(context.Background(), val, key, query)
}

// TakeCtx takes the result from cache first, if not found,
// query from DB and set cache using c.expiry, then return the result.
func (cc cacheCluster) TakeCtx(ctx context.Context, val any, key string, query func(val any) error) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).TakeCtx(ctx, val, key, query)
}

// TakeWithExpire takes the result from cache first, if not found,
// query from DB and set cache using given expire, then return the result.
func (cc cacheCluster) TakeWithExpire(val any, key string, query func(val any, expire time.Duration) error) error {
	return cc.TakeWithExpireCtx(context.Background(), val, key, query)
}

// TakeWithExpireCtx takes the result from cache first, if not found,
// query from DB and set cache using given expire, then return the result.
func (cc cacheCluster) TakeWithExpireCtx(ctx context.Context, val any, key string, query func(val any, expire time.Duration) error) error {
	c, ok := cc.dispatcher.Get(key)
	if !ok {
		return cc.errNotFound
	}

	return c.(Cache).TakeWithExpireCtx(ctx, val, key, query)
}
