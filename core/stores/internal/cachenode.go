package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/mathx"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/syncx"
)

const (
	notFoundPlaceholder = "*"
	// make the expiry unstable to avoid lots of cached items expire at the same time
	// make the unstable expiry to be [0.95, 1.05] * seconds
	expiryDeviation = 0.05
)

// indicates there is no such value associate with the key
var errPlaceholder = errors.New("placeholder")

type cacheNode struct {
	rds            *redis.Redis
	expiry         time.Duration
	notFoundExpiry time.Duration
	barrier        syncx.SharedCalls
	r              *rand.Rand
	lock           *sync.Mutex
	unstableExpiry mathx.Unstable
	stat           *CacheStat
	errNotFound    error
}

func NewCacheNode(rds *redis.Redis, barrier syncx.SharedCalls, st *CacheStat,
	errNotFound error, opts ...Option) Cache {
	o := newOptions(opts...)
	return cacheNode{
		rds:            rds,
		expiry:         o.Expiry,
		notFoundExpiry: o.NotFoundExpiry,
		barrier:        barrier,
		r:              rand.New(rand.NewSource(time.Now().UnixNano())),
		lock:           new(sync.Mutex),
		unstableExpiry: mathx.NewUnstable(expiryDeviation),
		stat:           st,
		errNotFound:    errNotFound,
	}
}

func (c cacheNode) DelCache(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	if _, err := c.rds.Del(keys...); err != nil {
		logx.Errorf("failed to clear cache with keys: %q, error: %v", formatKeys(keys), err)
		c.asyncRetryDelCache(keys...)
	}

	return nil
}

func (c cacheNode) GetCache(key string, v interface{}) error {
	if err := c.doGetCache(key, v); err == errPlaceholder {
		return c.errNotFound
	} else {
		return err
	}
}

func (c cacheNode) SetCache(key string, v interface{}) error {
	return c.SetCacheWithExpire(key, v, c.aroundDuration(c.expiry))
}

func (c cacheNode) SetCacheWithExpire(key string, v interface{}, expire time.Duration) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.rds.Setex(key, string(data), int(expire.Seconds()))
}

func (c cacheNode) String() string {
	return c.rds.Addr
}

func (c cacheNode) Take(v interface{}, key string, query func(v interface{}) error) error {
	return c.doTake(v, key, query, func(v interface{}) error {
		return c.SetCache(key, v)
	})
}

func (c cacheNode) TakeWithExpire(v interface{}, key string,
	query func(v interface{}, expire time.Duration) error) error {
	expire := c.aroundDuration(c.expiry)
	return c.doTake(v, key, func(v interface{}) error {
		return query(v, expire)
	}, func(v interface{}) error {
		return c.SetCacheWithExpire(key, v, expire)
	})
}

func (c cacheNode) aroundDuration(duration time.Duration) time.Duration {
	return c.unstableExpiry.AroundDuration(duration)
}

func (c cacheNode) asyncRetryDelCache(keys ...string) {
	AddCleanTask(func() error {
		_, err := c.rds.Del(keys...)
		return err
	}, keys...)
}

func (c cacheNode) doGetCache(key string, v interface{}) error {
	c.stat.IncrementTotal()
	data, err := c.rds.Get(key)
	if err != nil {
		c.stat.IncrementMiss()
		return err
	}

	if len(data) == 0 {
		c.stat.IncrementMiss()
		return c.errNotFound
	}

	c.stat.IncrementHit()
	if data == notFoundPlaceholder {
		return errPlaceholder
	}

	return c.processCache(key, data, v)
}

func (c cacheNode) doTake(v interface{}, key string, query func(v interface{}) error,
	cacheVal func(v interface{}) error) error {
	val, fresh, err := c.barrier.DoEx(key, func() (interface{}, error) {
		if err := c.doGetCache(key, v); err != nil {
			if err == errPlaceholder {
				return nil, c.errNotFound
			} else if err != c.errNotFound {
				// why we just return the error instead of query from db,
				// because we don't allow the disaster pass to the dbs.
				// fail fast, in case we bring down the dbs.
				return nil, err
			}

			if err = query(v); err == c.errNotFound {
				if err = c.setCacheWithNotFound(key); err != nil {
					logx.Error(err)
				}

				return nil, c.errNotFound
			} else if err != nil {
				c.stat.IncrementDbFails()
				return nil, err
			}

			if err = cacheVal(v); err != nil {
				logx.Error(err)
			}
		}

		return json.Marshal(v)
	})
	if err != nil {
		return err
	}
	if fresh {
		return nil
	} else {
		// got the result from previous ongoing query
		c.stat.IncrementTotal()
		c.stat.IncrementHit()
	}

	return json.Unmarshal(val.([]byte), v)
}

func (c cacheNode) processCache(key string, data string, v interface{}) error {
	err := json.Unmarshal([]byte(data), v)
	if err == nil {
		return nil
	}

	report := fmt.Sprintf("unmarshal cache, node: %s, key: %s, value: %s, error: %v",
		c.rds.Addr, key, data, err)
	logx.Error(report)
	stat.Report(report)
	if _, e := c.rds.Del(key); e != nil {
		logx.Errorf("delete invalid cache, node: %s, key: %s, value: %s, error: %v",
			c.rds.Addr, key, data, e)
	}

	// returns errNotFound to reload the value by the given queryFn
	return c.errNotFound
}

func (c cacheNode) setCacheWithNotFound(key string) error {
	return c.rds.Setex(key, notFoundPlaceholder, int(c.aroundDuration(c.notFoundExpiry).Seconds()))
}
