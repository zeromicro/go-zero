package redis

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	red "github.com/go-redis/redis"
	"github.com/tal-tech/go-zero/core/breaker"
	"github.com/tal-tech/go-zero/core/mapping"
)

const (
	ClusterType = "cluster"
	NodeType    = "node"
	Nil         = red.Nil

	blockingQueryTimeout = 5 * time.Second
	readWriteTimeout     = 2 * time.Second

	slowThreshold = time.Millisecond * 100
)

var ErrNilNode = errors.New("nil redis node")

type (
	Pair struct {
		Key   string
		Score int64
	}

	// Redis defines a redis node/cluster. It is thread-safe.
	Redis struct {
		Addr string
		Type string
		Pass string
		brk  breaker.Breaker
	}

	RedisNode interface {
		red.Cmdable
	}

	// GeoLocation is used with GeoAdd to add geospatial location.
	GeoLocation = red.GeoLocation
	// GeoRadiusQuery is used with GeoRadius to query geospatial index.
	GeoRadiusQuery = red.GeoRadiusQuery
	GeoPos         = red.GeoPos

	Pipeliner = red.Pipeliner

	// Z represents sorted set member.
	Z      = red.Z
	ZStore = red.ZStore

	IntCmd   = red.IntCmd
	FloatCmd = red.FloatCmd
)

func NewRedis(redisAddr, redisType string, redisPass ...string) *Redis {
	var pass string
	for _, v := range redisPass {
		pass = v
	}

	return &Redis{
		Addr: redisAddr,
		Type: redisType,
		Pass: pass,
		brk:  breaker.NewBreaker(),
	}
}

// BitCount is redis bitcount command implementation.
func (s *Redis) BitCount(key string, start, end int64) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.BitCount(key, &red.BitCount{
			Start: start,
			End:   end,
		}).Result()
		return err
	}, acceptable)

	return
}

// BitOpAnd is redis bit operation (and) command implementation.
func (s *Redis) BitOpAnd(destKey string, keys ...string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.BitOpAnd(destKey, keys...).Result()
		return err
	}, acceptable)

	return
}

// BitOpNot is redis bit operation (not) command implementation.
func (s *Redis) BitOpNot(destKey, key string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.BitOpNot(destKey, key).Result()
		return err
	}, acceptable)

	return
}

// BitOpOr is redis bit operation (or) command implementation.
func (s *Redis) BitOpOr(destKey string, keys ...string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.BitOpOr(destKey, keys...).Result()
		return err
	}, acceptable)

	return
}

// BitOpXor is redis bit operation (xor) command implementation.
func (s *Redis) BitOpXor(destKey string, keys ...string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.BitOpXor(destKey, keys...).Result()
		return err
	}, acceptable)

	return
}

// BitPos is redis bitpos command implementation.
func (s *Redis) BitPos(key string, bit int64, start, end int64) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.BitPos(key, bit, start, end).Result()
		return err
	}, acceptable)

	return
}

// Blpop uses passed in redis connection to execute blocking queries.
// Doesn't benefit from pooling redis connections of blocking queries
func (s *Redis) Blpop(redisNode RedisNode, key string) (string, error) {
	if redisNode == nil {
		return "", ErrNilNode
	}

	vals, err := redisNode.BLPop(blockingQueryTimeout, key).Result()
	if err != nil {
		return "", err
	}

	if len(vals) < 2 {
		return "", fmt.Errorf("no value on key: %s", key)
	}

	return vals[1], nil
}

func (s *Redis) BlpopEx(redisNode RedisNode, key string) (string, bool, error) {
	if redisNode == nil {
		return "", false, ErrNilNode
	}

	vals, err := redisNode.BLPop(blockingQueryTimeout, key).Result()
	if err != nil {
		return "", false, err
	}

	if len(vals) < 2 {
		return "", false, fmt.Errorf("no value on key: %s", key)
	}

	return vals[1], true, nil
}

func (s *Redis) Del(keys ...string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.Del(keys...).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Eval(script string, keys []string, args ...interface{}) (val interface{}, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.Eval(script, keys, args...).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Exists(key string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.Exists(key).Result()
		if err != nil {
			return err
		}

		val = v == 1
		return nil
	}, acceptable)

	return
}

func (s *Redis) Expire(key string, seconds int) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.Expire(key, time.Duration(seconds)*time.Second).Err()
	}, acceptable)
}

func (s *Redis) Expireat(key string, expireTime int64) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.ExpireAt(key, time.Unix(expireTime, 0)).Err()
	}, acceptable)
}

func (s *Redis) GeoAdd(key string, geoLocation ...*GeoLocation) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.GeoAdd(key, geoLocation...).Result()
		if err != nil {
			return err
		}

		val = v
		return nil
	}, acceptable)
	return
}

func (s *Redis) GeoDist(key string, member1, member2, unit string) (val float64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.GeoDist(key, member1, member2, unit).Result()
		if err != nil {
			return err
		}

		val = v
		return nil
	}, acceptable)
	return
}

func (s *Redis) GeoHash(key string, members ...string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.GeoHash(key, members...).Result()
		if err != nil {
			return err
		}

		val = v
		return nil
	}, acceptable)
	return
}

func (s *Redis) GeoRadius(key string, longitude, latitude float64, query *GeoRadiusQuery) (val []GeoLocation, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.GeoRadius(key, longitude, latitude, query).Result()
		if err != nil {
			return err
		}

		val = v
		return nil
	}, acceptable)
	return
}
func (s *Redis) GeoRadiusByMember(key, member string, query *GeoRadiusQuery) (val []GeoLocation, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.GeoRadiusByMember(key, member, query).Result()
		if err != nil {
			return err
		}

		val = v
		return nil
	}, acceptable)
	return
}

func (s *Redis) GeoPos(key string, members ...string) (val []*GeoPos, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.GeoPos(key, members...).Result()
		if err != nil {
			return err
		}

		val = v
		return nil
	}, acceptable)
	return
}

func (s *Redis) Get(key string) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		if val, err = conn.Get(key).Result(); err == red.Nil {
			return nil
		} else if err != nil {
			return err
		} else {
			return nil
		}
	}, acceptable)

	return
}

func (s *Redis) GetBit(key string, offset int64) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.GetBit(key, offset).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Hdel(key, field string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.HDel(key, field).Result()
		if err != nil {
			return err
		}

		val = v == 1
		return nil
	}, acceptable)

	return
}

func (s *Redis) Hexists(key, field string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HExists(key, field).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Hget(key, field string) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HGet(key, field).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Hgetall(key string) (val map[string]string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HGetAll(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Hincrby(key, field string, increment int) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.HIncrBy(key, field, int64(increment)).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Hkeys(key string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HKeys(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Hlen(key string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.HLen(key).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Hmget(key string, fields ...string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.HMGet(key, fields...).Result()
		if err != nil {
			return err
		}

		val = toStrings(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Hset(key, field, value string) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.HSet(key, field, value).Err()
	}, acceptable)
}

func (s *Redis) Hsetnx(key, field, value string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HSetNX(key, field, value).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Hmset(key string, fieldsAndValues map[string]string) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		vals := make(map[string]interface{}, len(fieldsAndValues))
		for k, v := range fieldsAndValues {
			vals[k] = v
		}

		return conn.HMSet(key, vals).Err()
	}, acceptable)
}

func (s *Redis) Hscan(key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		keys, cur, err = conn.HScan(key, cursor, match, count).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Hvals(key string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.HVals(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Incr(key string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.Incr(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Incrby(key string, increment int64) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.IncrBy(key, int64(increment)).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Keys(pattern string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.Keys(pattern).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Llen(key string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.LLen(key).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Lpop(key string) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.LPop(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Lpush(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.LPush(key, values...).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Lrange(key string, start int, stop int) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.LRange(key, int64(start), int64(stop)).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Lrem(key string, count int, value string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.LRem(key, int64(count), value).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Mget(keys ...string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.MGet(keys...).Result()
		if err != nil {
			return err
		}

		val = toStrings(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Persist(key string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.Persist(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Pfadd(key string, values ...interface{}) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.PFAdd(key, values...).Result()
		if err != nil {
			return err
		}

		val = v == 1
		return nil
	}, acceptable)

	return
}

func (s *Redis) Pfcount(key string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.PFCount(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Pfmerge(dest string, keys ...string) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		_, err = conn.PFMerge(dest, keys...).Result()
		return err
	}, acceptable)
}

func (s *Redis) Ping() (val bool) {
	// ignore error, error means false
	_ = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			val = false
			return nil
		}

		v, err := conn.Ping().Result()
		if err != nil {
			val = false
			return nil
		}

		val = v == "PONG"
		return nil
	}, acceptable)

	return
}

func (s *Redis) Pipelined(fn func(Pipeliner) error) (err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		_, err = conn.Pipelined(fn)
		return err

	}, acceptable)

	return
}

func (s *Redis) Rpop(key string) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.RPop(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Rpush(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.RPush(key, values...).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Sadd(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.SAdd(key, values...).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Scan(cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		keys, cur, err = conn.Scan(cursor, match, count).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) SetBit(key string, offset int64, value int) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		_, err = conn.SetBit(key, offset, value).Result()
		return err
	}, acceptable)
}

func (s *Redis) Sscan(key string, cursor uint64, match string, count int64) (keys []string, cur uint64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		keys, cur, err = conn.SScan(key, cursor, match, count).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Scard(key string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SCard(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Set(key string, value string) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.Set(key, value, 0).Err()
	}, acceptable)
}

func (s *Redis) Setex(key, value string, seconds int) error {
	return s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		return conn.Set(key, value, time.Duration(seconds)*time.Second).Err()
	}, acceptable)
}

func (s *Redis) Setnx(key, value string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SetNX(key, value, 0).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) SetnxEx(key, value string, seconds int) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SetNX(key, value, time.Duration(seconds)*time.Second).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Sismember(key string, value interface{}) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SIsMember(key, value).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Srem(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.SRem(key, values...).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Smembers(key string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SMembers(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Spop(key string) (val string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SPop(key).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Srandmember(key string, count int) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SRandMemberN(key, int64(count)).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Sunion(keys ...string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SUnion(keys...).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Sunionstore(destination string, keys ...string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.SUnionStore(destination, keys...).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Sdiff(keys ...string) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.SDiff(keys...).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Sdiffstore(destination string, keys ...string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.SDiffStore(destination, keys...).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Ttl(key string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		duration, err := conn.TTL(key).Result()
		if err != nil {
			return err
		}

		val = int(duration / time.Second)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zadd(key string, score int64, value string) (val bool, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZAdd(key, red.Z{
			Score:  float64(score),
			Member: value,
		}).Result()
		if err != nil {
			return err
		}

		val = v == 1
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zadds(key string, ps ...Pair) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		var zs []red.Z
		for _, p := range ps {
			z := red.Z{Score: float64(p.Score), Member: p.Key}
			zs = append(zs, z)
		}

		v, err := conn.ZAdd(key, zs...).Result()
		if err != nil {
			return err
		}

		val = v
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zcard(key string) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZCard(key).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zcount(key string, start, stop int64) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZCount(key, strconv.FormatInt(start, 10), strconv.FormatInt(stop, 10)).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zincrby(key string, increment int64, field string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZIncrBy(key, float64(increment), field).Result()
		if err != nil {
			return err
		}

		val = int64(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zscore(key string, value string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZScore(key, value).Result()
		if err != nil {
			return err
		}

		val = int64(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zrank(key, field string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.ZRank(key, field).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Zrem(key string, values ...interface{}) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZRem(key, values...).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zremrangebyscore(key string, start, stop int64) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZRemRangeByScore(key, strconv.FormatInt(start, 10),
			strconv.FormatInt(stop, 10)).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zremrangebyrank(key string, start, stop int64) (val int, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZRemRangeByRank(key, start, stop).Result()
		if err != nil {
			return err
		}

		val = int(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zrange(key string, start, stop int64) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.ZRange(key, start, stop).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) ZrangeWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZRangeWithScores(key, start, stop).Result()
		if err != nil {
			return err
		}

		val = toPairs(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) ZRevRangeWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZRevRangeWithScores(key, start, stop).Result()
		if err != nil {
			return err
		}

		val = toPairs(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) ZrangebyscoreWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZRangeByScoreWithScores(key, red.ZRangeBy{
			Min: strconv.FormatInt(start, 10),
			Max: strconv.FormatInt(stop, 10),
		}).Result()
		if err != nil {
			return err
		}

		val = toPairs(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) ZrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) (
	val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		if size <= 0 {
			return nil
		}

		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZRangeByScoreWithScores(key, red.ZRangeBy{
			Min:    strconv.FormatInt(start, 10),
			Max:    strconv.FormatInt(stop, 10),
			Offset: int64(page * size),
			Count:  int64(size),
		}).Result()
		if err != nil {
			return err
		}

		val = toPairs(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zrevrange(key string, start, stop int64) (val []string, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.ZRevRange(key, start, stop).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) ZrevrangebyscoreWithScores(key string, start, stop int64) (val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZRevRangeByScoreWithScores(key, red.ZRangeBy{
			Min: strconv.FormatInt(start, 10),
			Max: strconv.FormatInt(stop, 10),
		}).Result()
		if err != nil {
			return err
		}

		val = toPairs(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) ZrevrangebyscoreWithScoresAndLimit(key string, start, stop int64, page, size int) (
	val []Pair, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		if size <= 0 {
			return nil
		}

		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		v, err := conn.ZRevRangeByScoreWithScores(key, red.ZRangeBy{
			Min:    strconv.FormatInt(start, 10),
			Max:    strconv.FormatInt(stop, 10),
			Offset: int64(page * size),
			Count:  int64(size),
		}).Result()
		if err != nil {
			return err
		}

		val = toPairs(v)
		return nil
	}, acceptable)

	return
}

func (s *Redis) Zrevrank(key string, field string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.ZRevRank(key, field).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) Zunionstore(dest string, store ZStore, keys ...string) (val int64, err error) {
	err = s.brk.DoWithAcceptable(func() error {
		conn, err := getRedis(s)
		if err != nil {
			return err
		}

		val, err = conn.ZUnionStore(dest, store, keys...).Result()
		return err
	}, acceptable)

	return
}

func (s *Redis) String() string {
	return s.Addr
}

func (s *Redis) scriptLoad(script string) (string, error) {
	conn, err := getRedis(s)
	if err != nil {
		return "", err
	}

	return conn.ScriptLoad(script).Result()
}

func acceptable(err error) bool {
	return err == nil || err == red.Nil
}

func getRedis(r *Redis) (RedisNode, error) {
	switch r.Type {
	case ClusterType:
		return getCluster(r.Addr, r.Pass)
	case NodeType:
		return getClient(r.Addr, r.Pass)
	default:
		return nil, fmt.Errorf("redis type '%s' is not supported", r.Type)
	}
}

func toPairs(vals []red.Z) []Pair {
	pairs := make([]Pair, len(vals))
	for i, val := range vals {
		switch member := val.Member.(type) {
		case string:
			pairs[i] = Pair{
				Key:   member,
				Score: int64(val.Score),
			}
		default:
			pairs[i] = Pair{
				Key:   mapping.Repr(val.Member),
				Score: int64(val.Score),
			}
		}
	}
	return pairs
}

func toStrings(vals []interface{}) []string {
	ret := make([]string, len(vals))
	for i, val := range vals {
		if val == nil {
			ret[i] = ""
		} else {
			switch val := val.(type) {
			case string:
				ret[i] = val
			default:
				ret[i] = mapping.Repr(val)
			}
		}
	}
	return ret
}
