package redis

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	red "github.com/go-redis/redis"
	"github.com/tal-tech/go-zero/core/logx"
)

const (
	letters     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
	randomLen       = 16
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
)

// A RedisLock is a redis lock.
type RedisLock struct {
	store   *Redis
	seconds uint32
	key     string
	id      string
	lockSHA string
	delSHA  string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewRedisLock returns a RedisLock.
func NewRedisLock(store *Redis, key string) *RedisLock {
	lockSHA, err := store.ScriptLoad(lockCommand)
	if err != nil {
		logx.Error("Error on failed to load lua script in lockCommand")
		return nil
	}
	delSHA, err := store.ScriptLoad(delCommand)
	if err != nil {
		logx.Error("Error on failed to load lua script in delCommand")
		return nil
	}

	return &RedisLock{
		store:   store,
		key:     key,
		id:      randomStr(randomLen),
		lockSHA: lockSHA,
		delSHA:  delSHA,
	}
}

// Acquire acquires the lock.
func (rl *RedisLock) Acquire() (bool, error) {
	seconds := atomic.LoadUint32(&rl.seconds)
	resp, err := rl.store.EvalSha(rl.lockSHA, []string{rl.key}, []string{
		rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
	})
	if err == red.Nil {
		return false, nil
	} else if err != nil {
		if isLuaScriptDone(err) {
			if err = rl.reloadScript("lock", lockCommand); err != nil {
				logx.Errorf("Error on reload lockscript: %s", err)
				return false, err
			}
		} else {
			logx.Errorf("Error on acquiring lock for %s, %s", rl.key, err.Error())
			return false, err
		}
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
	resp, err := rl.store.EvalSha(rl.delSHA, []string{rl.key}, []string{rl.id})
	if err == red.Nil {
		return false, nil
	} else if err != nil {
		if isLuaScriptDone(err) {
			if err = rl.reloadScript("del", lockCommand); err != nil {
				logx.Errorf("Error on reload lockscript: %s", err)
				return false, err
			}
		} else {
			logx.Errorf("Error on acquiring lock for %s, %s", rl.key, err.Error())
			return false, err
		}
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}

	return reply == 1, nil
}

// SetExpire sets the expire.
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}

func (rl *RedisLock) reloadScript(loadType string, script string) error {
	if rl == nil {
		return errors.New("redislock instance was not initialized")
	}

	load, err := rl.store.ScriptLoad(script)
	if err != nil {
		logx.Errorf("Error on failed to load lua script in lockCommand")
		return err
	}
	switch loadType {
	case "lock":
		rl.lockSHA = load
	case "del":
		rl.delSHA = load
	default:
		return errors.New("script initialization is not supported")
	}
	return nil
}

func isLuaScriptDone(err error) bool {
	return strings.Contains(err.Error(), "NOSCRIPT")
}

func randomStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
