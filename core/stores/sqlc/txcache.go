package sqlc

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
)

type (
	actionCollector interface {
		AddAction(action func(ctx context.Context) error)
	}

	txCache struct {
		cache.Cache
		collector actionCollector
	}
)

func (c *txCache) Set(key string, val any) error {
	c.collector.AddAction(func(ctx context.Context) error {
		return c.Cache.Set(key, val)
	})

	return nil
}

func (c *txCache) SetCtx(ctx context.Context, key string, val any) error {
	c.collector.AddAction(func(ctx context.Context) error {
		return c.Cache.SetCtx(ctx, key, val)
	})

	return nil
}

func (c *txCache) SetWithExpire(key string, val any, expire time.Duration) error {
	fmt.Println("txCache SetWithExpire")
	c.collector.AddAction(func(ctx context.Context) error {
		return c.Cache.SetWithExpire(key, val, expire)
	})

	return nil
}

func (c *txCache) SetWithExpireCtx(ctx context.Context, key string, val any, expire time.Duration) error {
	fmt.Println("txCache SetWithExpireCtx")
	c.collector.AddAction(func(ctx context.Context) error {
		return c.Cache.SetWithExpireCtx(ctx, key, val, expire)
	})

	return nil
}

func (c *txCache) TakeCtx(ctx context.Context, val any, key string, query func(val any) error) error {
	fmt.Println("txCache Take")
	return c.Cache.TakeCtx(context.Background(), val, key, query)
}
