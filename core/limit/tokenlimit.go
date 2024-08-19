package limit

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	xrate "golang.org/x/time/rate"
)

const (
	tokenFormat     = "{%s}.tokens"
	timestampFormat = "{%s}.ts"
	pingInterval    = time.Millisecond * 100
)

var (
	//go:embed tokenscript.lua
	tokenLuaScript string
	tokenScript    = redis.NewScript(tokenLuaScript)
)

// A TokenLimiter controls how frequently events are allowed to happen with in one second.
type TokenLimiter struct {
	rate           int
	burst          int
	store          *redis.Redis
	tokenKey       string
	timestampKey   string
	rescueLock     sync.Mutex
	redisAlive     uint32
	monitorStarted bool
	rescueLimiter  *xrate.Limiter
}

// NewTokenLimiter returns a new TokenLimiter that allows events up to rate and permits
// bursts of at most burst tokens.
func NewTokenLimiter(rate, burst int, store *redis.Redis, key string) *TokenLimiter {
	tokenKey := fmt.Sprintf(tokenFormat, key)
	timestampKey := fmt.Sprintf(timestampFormat, key)

	return &TokenLimiter{
		rate:          rate,
		burst:         burst,
		store:         store,
		tokenKey:      tokenKey,
		timestampKey:  timestampKey,
		redisAlive:    1,
		rescueLimiter: xrate.NewLimiter(xrate.Every(time.Second/time.Duration(rate)), burst),
	}
}

// Allow is shorthand for AllowN(time.Now(), 1).
func (lim *TokenLimiter) Allow() bool {
	return lim.AllowN(time.Now(), 1)
}

// AllowCtx is shorthand for AllowNCtx(ctx,time.Now(), 1) with incoming context.
func (lim *TokenLimiter) AllowCtx(ctx context.Context) bool {
	return lim.AllowNCtx(ctx, time.Now(), 1)
}

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate.
// Otherwise, use Reserve or Wait.
func (lim *TokenLimiter) AllowN(now time.Time, n int) bool {
	return lim.reserveN(context.Background(), now, n)
}

// AllowNCtx reports whether n events may happen at time now with incoming context.
// Use this method if you intend to drop / skip events that exceed the rate.
// Otherwise, use Reserve or Wait.
func (lim *TokenLimiter) AllowNCtx(ctx context.Context, now time.Time, n int) bool {
	return lim.reserveN(ctx, now, n)
}

func (lim *TokenLimiter) reserveN(ctx context.Context, now time.Time, n int) bool {
	if atomic.LoadUint32(&lim.redisAlive) == 0 {
		return lim.rescueLimiter.AllowN(now, n)
	}

	resp, err := lim.store.ScriptRunCtx(ctx,
		tokenScript,
		[]string{
			lim.tokenKey,
			lim.timestampKey,
		},
		[]string{
			strconv.Itoa(lim.rate),
			strconv.Itoa(lim.burst),
			strconv.FormatInt(now.Unix(), 10),
			strconv.Itoa(n),
		})
	// redis allowed == false
	// Lua boolean false -> r Nil bulk reply
	if errors.Is(err, redis.Nil) {
		return false
	}
	if errorx.In(err, context.DeadlineExceeded, context.Canceled) {
		logx.Errorf("fail to use rate limiter: %s", err)
		return false
	}
	if err != nil {
		logx.Errorf("fail to use rate limiter: %s, use in-process limiter for rescue", err)
		lim.startMonitor()
		return lim.rescueLimiter.AllowN(now, n)
	}

	code, ok := resp.(int64)
	if !ok {
		logx.Errorf("fail to eval redis script: %v, use in-process limiter for rescue", resp)
		lim.startMonitor()
		return lim.rescueLimiter.AllowN(now, n)
	}

	// redis allowed == true
	// Lua boolean true -> r integer reply with value of 1
	return code == 1
}

func (lim *TokenLimiter) startMonitor() {
	lim.rescueLock.Lock()
	defer lim.rescueLock.Unlock()

	if lim.monitorStarted {
		return
	}

	lim.monitorStarted = true
	atomic.StoreUint32(&lim.redisAlive, 0)

	go lim.waitForRedis()
}

func (lim *TokenLimiter) waitForRedis() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		lim.rescueLock.Lock()
		lim.monitorStarted = false
		lim.rescueLock.Unlock()
	}()

	for range ticker.C {
		if lim.store.Ping() {
			atomic.StoreUint32(&lim.redisAlive, 1)
			return
		}
	}
}
