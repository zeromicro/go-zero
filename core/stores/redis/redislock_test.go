package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stringx"
)

func TestRedisLock(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
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
	})
}
