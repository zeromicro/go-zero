package cache

import (
	"context"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/collection"
	"github.com/zeromicro/go-zero/core/jsonx"
	"github.com/zeromicro/go-zero/core/logx"
)

type (
	// MultiLevelCacheOption defines the method to customize a multiLevelCache.
	MultiLevelCacheOption func(opts *multiLevelCacheOptions)

	multiLevelCacheOptions struct {
		localExpire time.Duration
		localLimit  int
	}

	multiLevelCache struct {
		local       *collection.Cache
		remote      Cache
		errNotFound error
	}
)

const (
	defaultLocalExpire = time.Minute * 5
	defaultLocalLimit  = 10000
	// placeholder for not-found entries in local cache, must match cacheNode's placeholder.
	localNotFoundPlaceholder = "*"
)

// NewMultiLevelCache returns a Cache that combines an in-memory collection.Cache (L1)
// with a remote cache.Cache (L2, typically Redis-backed).
//
// On reads, L1 is checked first for minimal latency; on a miss, L2 is queried
// and the result is promoted into L1. Writes and deletes are applied to both layers.
//
// Usage with sqlc:
//
//	mlc, err := cache.NewMultiLevelCache(remoteCache, sql.ErrNoRows)
//	conn := sqlc.NewConnWithCache(db, mlc)
func NewMultiLevelCache(remote Cache, errNotFound error,
	opts ...MultiLevelCacheOption) (Cache, error) {
	o := multiLevelCacheOptions{
		localExpire: defaultLocalExpire,
		localLimit:  defaultLocalLimit,
	}
	for _, opt := range opts {
		opt(&o)
	}

	local, err := collection.NewCache(o.localExpire,
		collection.WithLimit(o.localLimit),
		collection.WithName("multilevel"))
	if err != nil {
		return nil, err
	}

	return &multiLevelCache{
		local:       local,
		remote:      remote,
		errNotFound: errNotFound,
	}, nil
}

// WithLocalExpire customizes the in-memory cache expiry duration.
func WithLocalExpire(d time.Duration) MultiLevelCacheOption {
	return func(opts *multiLevelCacheOptions) {
		if d > 0 {
			opts.localExpire = d
		}
	}
}

// WithLocalLimit customizes the in-memory cache max entry count.
func WithLocalLimit(limit int) MultiLevelCacheOption {
	return func(opts *multiLevelCacheOptions) {
		if limit > 0 {
			opts.localLimit = limit
		}
	}
}

// Del deletes cached values with keys.
func (mc *multiLevelCache) Del(keys ...string) error {
	return mc.DelCtx(context.Background(), keys...)
}

// DelCtx deletes cached values with keys.
func (mc *multiLevelCache) DelCtx(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		mc.local.Del(key)
	}
	return mc.remote.DelCtx(ctx, keys...)
}

// Get gets the cache with key and fills into v.
func (mc *multiLevelCache) Get(key string, val any) error {
	return mc.GetCtx(context.Background(), key, val)
}

// GetCtx gets the cache with key and fills into v.
// It checks the local in-memory cache first, then falls back to the remote cache.
func (mc *multiLevelCache) GetCtx(ctx context.Context, key string, val any) error {
	if data, ok := mc.local.Get(key); ok {
		bs, ok := data.([]byte)
		if ok {
			if string(bs) == localNotFoundPlaceholder {
				return mc.errNotFound
			}
			return jsonx.Unmarshal(bs, val)
		}
	}

	// L1 miss, try L2
	if err := mc.remote.GetCtx(ctx, key, val); err != nil {
		return err
	}

	// Promote to L1
	mc.promoteToLocal(key, val)
	return nil
}

// IsNotFound checks if the given error is the defined errNotFound.
func (mc *multiLevelCache) IsNotFound(err error) bool {
	return errors.Is(err, mc.errNotFound)
}

// Set sets the cache with key and v, using default expiry.
func (mc *multiLevelCache) Set(key string, val any) error {
	return mc.SetCtx(context.Background(), key, val)
}

// SetCtx sets the cache with key and v, using default expiry.
func (mc *multiLevelCache) SetCtx(ctx context.Context, key string, val any) error {
	if err := mc.remote.SetCtx(ctx, key, val); err != nil {
		return err
	}

	mc.promoteToLocal(key, val)
	return nil
}

// SetWithExpire sets the cache with key and v, using given expire.
func (mc *multiLevelCache) SetWithExpire(key string, val any, expire time.Duration) error {
	return mc.SetWithExpireCtx(context.Background(), key, val, expire)
}

// SetWithExpireCtx sets the cache with key and v, using given expire.
func (mc *multiLevelCache) SetWithExpireCtx(ctx context.Context, key string, val any,
	expire time.Duration) error {
	if err := mc.remote.SetWithExpireCtx(ctx, key, val, expire); err != nil {
		return err
	}

	mc.promoteToLocal(key, val)
	return nil
}

// Take takes the result from cache first, if not found,
// query from DB and set cache using default expiry, then return the result.
func (mc *multiLevelCache) Take(val any, key string, query func(val any) error) error {
	return mc.TakeCtx(context.Background(), val, key, query)
}

// TakeCtx takes the result from cache first, if not found,
// query from DB and set cache using default expiry, then return the result.
func (mc *multiLevelCache) TakeCtx(ctx context.Context, val any, key string,
	query func(val any) error) error {
	// Check L1 first
	if data, ok := mc.local.Get(key); ok {
		bs, ok := data.([]byte)
		if ok {
			if string(bs) == localNotFoundPlaceholder {
				return mc.errNotFound
			}
			if err := jsonx.Unmarshal(bs, val); err == nil {
				return nil
			}
			// unmarshal failed, evict bad entry and fall through
			mc.local.Del(key)
		}
	}

	// L1 miss, delegate to L2
	if err := mc.remote.TakeCtx(ctx, val, key, query); err != nil {
		if mc.IsNotFound(err) {
			mc.local.Set(key, []byte(localNotFoundPlaceholder))
		}
		return err
	}

	// Promote to L1
	mc.promoteToLocal(key, val)
	return nil
}

// TakeWithExpire takes the result from cache first, if not found,
// query from DB and set cache using given expire, then return the result.
func (mc *multiLevelCache) TakeWithExpire(val any, key string,
	query func(val any, expire time.Duration) error) error {
	return mc.TakeWithExpireCtx(context.Background(), val, key, query)
}

// TakeWithExpireCtx takes the result from cache first, if not found,
// query from DB and set cache using given expire, then return the result.
func (mc *multiLevelCache) TakeWithExpireCtx(ctx context.Context, val any, key string,
	query func(val any, expire time.Duration) error) error {
	// Check L1 first
	if data, ok := mc.local.Get(key); ok {
		bs, ok := data.([]byte)
		if ok {
			if string(bs) == localNotFoundPlaceholder {
				return mc.errNotFound
			}
			if err := jsonx.Unmarshal(bs, val); err == nil {
				return nil
			}
			mc.local.Del(key)
		}
	}

	// L1 miss, delegate to L2
	if err := mc.remote.TakeWithExpireCtx(ctx, val, key, query); err != nil {
		if mc.IsNotFound(err) {
			mc.local.Set(key, []byte(localNotFoundPlaceholder))
		}
		return err
	}

	// Promote to L1
	mc.promoteToLocal(key, val)
	return nil
}

func (mc *multiLevelCache) promoteToLocal(key string, val any) {
	data, err := jsonx.Marshal(val)
	if err != nil {
		logx.Errorf("multilevel cache: failed to marshal value for key %q: %v", key, err)
		return
	}
	mc.local.Set(key, data)
}
