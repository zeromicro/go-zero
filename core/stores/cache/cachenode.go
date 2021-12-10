package cache

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/tal-tech/go-zero/core/jsonx"
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
	barrier        syncx.SingleFlight
	r              *rand.Rand
	lock           *sync.Mutex
	unstableExpiry mathx.Unstable
	stat           *Stat
	errNotFound    error
}

// NewNode returns a cacheNode.
// rds is the underlying redis node or cluster.
// barrier is the barrier that maybe shared with other cache nodes on cache cluster.
// st is used to stat the cache.
// errNotFound defines the error that returned on cache not found.
// opts are the options that customize the cacheNode.
func NewNode(rds *redis.Redis, barrier syncx.SingleFlight, st *Stat,
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

// Del deletes cached values with keys.
func (c cacheNode) Del(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	if len(keys) > 1 && c.rds.Type == redis.ClusterType {
		for _, key := range keys {
			if _, err := c.rds.Del(key); err != nil {
				logx.Errorf("failed to clear cache with key: %q, error: %v", key, err)
				c.asyncRetryDelCache(key)
			}
		}
	} else {
		if _, err := c.rds.Del(keys...); err != nil {
			logx.Errorf("failed to clear cache with keys: %q, error: %v", formatKeys(keys), err)
			c.asyncRetryDelCache(keys...)
		}
	}

	return nil
}

// Get gets the cache with key and fills into v.
func (c cacheNode) Get(key string, v interface{}) error {
	err := c.doGetCache(key, v)
	if err == errPlaceholder {
		return c.errNotFound
	}

	return err
}

// IsNotFound checks if the given error is the defined errNotFound.
func (c cacheNode) IsNotFound(err error) bool {
	return err == c.errNotFound
}

// Set sets the cache with key and v, using c.expiry.
func (c cacheNode) Set(key string, v interface{}) error {
	return c.SetWithExpire(key, v, c.aroundDuration(c.expiry))
}

// SetWithExpire sets the cache with key and v, using given expire.
func (c cacheNode) SetWithExpire(key string, v interface{}, expire time.Duration) error {
	data, err := jsonx.Marshal(v)
	if err != nil {
		return err
	}

	return c.rds.Setex(key, string(data), int(expire.Seconds()))
}

// String returns a string that represents the cacheNode.
func (c cacheNode) String() string {
	return c.rds.Addr
}

// Take takes the result from cache first, if not found,
// query from DB and set cache using c.expiry, then return the result.
func (c cacheNode) Take(v interface{}, key string, query func(v interface{}) error) error {
	return c.doTake(v, key, query, func(v interface{}) error {
		return c.Set(key, v)
	})
}

// TakeWithExpire takes the result from cache first, if not found,
// query from DB and set cache using given expire, then return the result.
func (c cacheNode) TakeWithExpire(v interface{}, key string, query func(v interface{},
	expire time.Duration) error) error {
	expire := c.aroundDuration(c.expiry)
	return c.doTake(v, key, func(v interface{}) error {
		return query(v, expire)
	}, func(v interface{}) error {
		return c.SetWithExpire(key, v, expire)
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

		return jsonx.Marshal(v)
	})
	if err != nil {
		return err
	}
	if fresh {
		return nil
	}

	// got the result from previous ongoing query
	c.stat.IncrementTotal()
	c.stat.IncrementHit()

	return jsonx.Unmarshal(val.([]byte), v)
}

func (c cacheNode) processCache(key, data string, v interface{}) error {
	err := jsonx.Unmarshal([]byte(data), v)
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
