package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zeromicro/go-zero/core/stringx"
)

func TestRedisLock(t *testing.T) {
	testFn := func(ctx context.Context) func(client *Redis) {
		return func(client *Redis) {
			key := stringx.Rand()
			firstLock := NewRedisLock(client, key)
			firstLock.SetExpire(5)
			firstAcquire, err := firstLock.Acquire()
			assert.Nil(t, err)
			assert.True(t, firstAcquire)

			secondLock := NewRedisLock(client, key)
			secondLock.SetExpire(5)
			againAcquire, err := secondLock.Acquire()
			assert.Nil(t, err)
			assert.False(t, againAcquire)

			release, err := firstLock.Release()
			assert.Nil(t, err)
			assert.True(t, release)

			endAcquire, err := secondLock.Acquire()
			assert.Nil(t, err)
			assert.True(t, endAcquire)
		}
	}

	t.Run("normal", func(t *testing.T) {
		runOnRedis(t, testFn(nil))
	})

	t.Run("withContext", func(t *testing.T) {
		runOnRedis(t, testFn(context.Background()))
	})
}

func TestRedisLock_Expired(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		key := stringx.Rand()
		redisLock := NewRedisLock(client, key)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := redisLock.AcquireCtx(ctx)
		assert.NotNil(t, err)
	})

	runOnRedis(t, func(client *Redis) {
		key := stringx.Rand()
		redisLock := NewRedisLock(client, key)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := redisLock.ReleaseCtx(ctx)
		assert.NotNil(t, err)
	})
}
