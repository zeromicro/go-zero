package redis

import (
	"context"
	_ "embed"
	"errors"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	red "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
)

const (
	randomLen       = 16
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
)

var (
	//go:embed lockscript.lua
	lockLuaScript string
	lockScript    = NewScript(lockLuaScript)

	//go:embed delscript.lua
	delLuaScript string
	delScript    = NewScript(delLuaScript)
)

// A RedisLock is a redis lock.
type RedisLock struct {
	store   *Redis
	seconds uint32
	key     string
	id      string
}

func init() {
	rand.NewSource(time.Now().UnixNano())
}

// NewRedisLock returns a RedisLock.
func NewRedisLock(store *Redis, key string) *RedisLock {
	return &RedisLock{
		store: store,
		key:   key,
		id:    stringx.Randn(randomLen),
	}
}

// Acquire acquires the lock.
func (rl *RedisLock) Acquire() (bool, error) {
	return rl.AcquireCtx(context.Background())
}

// AcquireCtx acquires the lock with the given ctx.
func (rl *RedisLock) AcquireCtx(ctx context.Context) (bool, error) {
	seconds := atomic.LoadUint32(&rl.seconds)
	resp, err := rl.store.ScriptRunCtx(ctx, lockScript, []string{rl.key}, []string{
		rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
	})
	if errors.Is(err, red.Nil) {
		return false, nil
	} else if err != nil {
		logx.Errorf("Error on acquiring lock for %s, %s", rl.key, err.Error())
		return false, err
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}

	logx.Errorf("Unknown reply when acquiring lock for %s: %v", rl.key, resp)
	return false, nil
}

// Release releases the lock.
func (rl *RedisLock) Release() (bool, error) {
	return rl.ReleaseCtx(context.Background())
}

// ReleaseCtx releases the lock with the given ctx.
func (rl *RedisLock) ReleaseCtx(ctx context.Context) (bool, error) {
	resp, err := rl.store.ScriptRunCtx(ctx, delScript, []string{rl.key}, []string{rl.id})
	if err != nil {
		return false, err
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return reply == 1, nil
}

// SetExpire sets the expiration.
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}
