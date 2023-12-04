package cache

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mathx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
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
	return c.DelCtx(context.Background(), keys...)
}

// DelCtx deletes cached values with keys.
func (c cacheNode) DelCtx(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	logger := logx.WithContext(ctx)
	if len(keys) > 1 && c.rds.Type == redis.ClusterType {
		for _, key := range keys {
			if _, err := c.rds.DelCtx(ctx, key); err != nil {
				logger.Errorf("failed to clear cache with key: %q, error: %v", key, err)
				c.asyncRetryDelCache(key)
			}
		}
	} else if _, err := c.rds.DelCtx(ctx, keys...); err != nil {
		logger.Errorf("failed to clear cache with keys: %q, error: %v", formatKeys(keys), err)
		c.asyncRetryDelCache(keys...)
	}

	return nil
}

// Get gets the cache with key and fills into v.
func (c cacheNode) Get(key string, val any) error {
	return c.GetCtx(context.Background(), key, val)
}

// GetCtx gets the cache with key and fills into v.
func (c cacheNode) GetCtx(ctx context.Context, key string, val any) error {
	err := c.doGetCache(ctx, key, val)
	if errors.Is(err, errPlaceholder) {
		return c.errNotFound
	}

	return err
}

// IsNotFound checks if the given error is the defined errNotFound.
func (c cacheNode) IsNotFound(err error) bool {
	return errors.Is(err, c.errNotFound)
}

// Set sets the cache with key and v, using c.expiry.
func (c cacheNode) Set(key string, val any) error {
	return c.SetCtx(context.Background(), key, val)
}

// SetCtx sets the cache with key and v, using c.expiry.
func (c cacheNode) SetCtx(ctx context.Context, key string, val any) error {
	return c.SetWithExpireCtx(ctx, key, val, c.aroundDuration(c.expiry))
}

// SetWithExpire sets the cache with key and v, using given expire.
func (c cacheNode) SetWithExpire(key string, val any, expire time.Duration) error {
	return c.SetWithExpireCtx(context.Background(), key, val, expire)
}

// SetWithExpireCtx sets the cache with key and v, using given expire.
func (c cacheNode) SetWithExpireCtx(ctx context.Context, key string, val any,
	expire time.Duration) error {
	data, err := jsonx.Marshal(val)
	if err != nil {
		return err
	}

	return c.rds.SetexCtx(ctx, key, string(data), int(math.Ceil(expire.Seconds())))
}

// String returns a string that represents the cacheNode.
func (c cacheNode) String() string {
	return c.rds.Addr
}

// Take takes the result from cache first, if not found,
// query from DB and set cache using c.expiry, then return the result.
func (c cacheNode) Take(val any, key string, query func(val any) error) error {
	return c.TakeCtx(context.Background(), val, key, query)
}

// TakeCtx takes the result from cache first, if not found,
// query from DB and set cache using c.expiry, then return the result.
func (c cacheNode) TakeCtx(ctx context.Context, val any, key string,
	query func(val any) error) error {
	return c.doTake(ctx, val, key, query, func(v any) error {
		return c.SetCtx(ctx, key, v)
	})
}

// TakeWithExpire takes the result from cache first, if not found,
// query from DB and set cache using given expire, then return the result.
func (c cacheNode) TakeWithExpire(val any, key string, query func(val any,
	expire time.Duration) error) error {
	return c.TakeWithExpireCtx(context.Background(), val, key, query)
}

// TakeWithExpireCtx takes the result from cache first, if not found,
// query from DB and set cache using given expire, then return the result.
func (c cacheNode) TakeWithExpireCtx(ctx context.Context, val any, key string,
	query func(val any, expire time.Duration) error) error {
	expire := c.aroundDuration(c.expiry)
	return c.doTake(ctx, val, key, func(v any) error {
		return query(v, expire)
	}, func(v any) error {
		return c.SetWithExpireCtx(ctx, key, v, expire)
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

func (c cacheNode) doGetCache(ctx context.Context, key string, v any) error {
	c.stat.IncrementTotal()
	data, err := c.rds.GetCtx(ctx, key)
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

	return c.processCache(ctx, key, data, v)
}

func (c cacheNode) doTake(ctx context.Context, v any, key string,
	query func(v any) error, cacheVal func(v any) error) error {
	logger := logx.WithContext(ctx)
	val, fresh, err := c.barrier.DoEx(key, func() (any, error) {
		if err := c.doGetCache(ctx, key, v); err != nil {
			if errors.Is(err, errPlaceholder) {
				return nil, c.errNotFound
			} else if !errors.Is(err, c.errNotFound) {
				// why we just return the error instead of query from db,
				// because we don't allow the disaster pass to the dbs.
				// fail fast, in case we bring down the dbs.
				return nil, err
			}

			if err = query(v); errors.Is(err, c.errNotFound) {
				if err = c.setCacheWithNotFound(ctx, key); err != nil {
					logger.Error(err)
				}

				return nil, c.errNotFound
			} else if err != nil {
				c.stat.IncrementDbFails()
				return nil, err
			}

			if err = cacheVal(v); err != nil {
				logger.Error(err)
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

	// got the result from previous ongoing query.
	// why not call IncrementTotal at the beginning of this function?
	// because a shared error is returned, and we don't want to count.
	// for example, if the db is down, the query will be failed, we count
	// the shared errors with one db failure.
	c.stat.IncrementTotal()
	c.stat.IncrementHit()

	return jsonx.Unmarshal(val.([]byte), v)
}

func (c cacheNode) processCache(ctx context.Context, key, data string, v any) error {
	err := jsonx.Unmarshal([]byte(data), v)
	if err == nil {
		return nil
	}

	report := fmt.Sprintf("unmarshal cache, node: %s, key: %s, value: %s, error: %v",
		c.rds.Addr, key, data, err)
	logger := logx.WithContext(ctx)
	logger.Error(report)
	stat.Report(report)
	if _, e := c.rds.DelCtx(ctx, key); e != nil {
		logger.Errorf("delete invalid cache, node: %s, key: %s, value: %s, error: %v",
			c.rds.Addr, key, data, e)
	}

	// returns errNotFound to reload the value by the given queryFn
	return c.errNotFound
}

func (c cacheNode) setCacheWithNotFound(ctx context.Context, key string) error {
	seconds := int(math.Ceil(c.aroundDuration(c.notFoundExpiry).Seconds()))
	_, err := c.rds.SetnxExCtx(ctx, key, notFoundPlaceholder, seconds)
	return err
}
