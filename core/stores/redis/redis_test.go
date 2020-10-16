package redis

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	red "github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestRedis_Exists(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := NewRedis(client.Addr, "").Exists("a")
		assert.NotNil(t, err)
		ok, err := client.Exists("a")
		assert.Nil(t, err)
		assert.False(t, ok)
		assert.Nil(t, client.Set("a", "b"))
		ok, err = client.Exists("a")
		assert.Nil(t, err)
		assert.True(t, ok)
	})
}

func TestRedis_Eval(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := NewRedis(client.Addr, "").Eval(`redis.call("EXISTS", KEYS[1])`, []string{"notexist"})
		assert.NotNil(t, err)
		_, err = client.Eval(`redis.call("EXISTS", KEYS[1])`, []string{"notexist"})
		assert.Equal(t, Nil, err)
		err = client.Set("key1", "value1")
		assert.Nil(t, err)
		_, err = client.Eval(`redis.call("EXISTS", KEYS[1])`, []string{"key1"})
		assert.Equal(t, Nil, err)
		val, err := client.Eval(`return redis.call("EXISTS", KEYS[1])`, []string{"key1"})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
	})
}

func TestRedis_Hgetall(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		_, err := NewRedis(client.Addr, "").Hgetall("a")
		assert.NotNil(t, err)
		vals, err := client.Hgetall("a")
		assert.Nil(t, err)
		assert.EqualValues(t, map[string]string{
			"aa": "aaa",
			"bb": "bbb",
		}, vals)
	})
}

func TestRedis_Hvals(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.NotNil(t, NewRedis(client.Addr, "").Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		_, err := NewRedis(client.Addr, "").Hvals("a")
		assert.NotNil(t, err)
		vals, err := client.Hvals("a")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"aaa", "bbb"}, vals)
	})
}

func TestRedis_Hsetnx(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		_, err := NewRedis(client.Addr, "").Hsetnx("a", "bb", "ccc")
		assert.NotNil(t, err)
		ok, err := client.Hsetnx("a", "bb", "ccc")
		assert.Nil(t, err)
		assert.False(t, ok)
		ok, err = client.Hsetnx("a", "dd", "ddd")
		assert.Nil(t, err)
		assert.True(t, ok)
		vals, err := client.Hvals("a")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"aaa", "bbb", "ddd"}, vals)
	})
}

func TestRedis_HdelHlen(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		_, err := NewRedis(client.Addr, "").Hlen("a")
		assert.NotNil(t, err)
		num, err := client.Hlen("a")
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		val, err := client.Hdel("a", "aa")
		assert.Nil(t, err)
		assert.True(t, val)
		vals, err := client.Hvals("a")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"bbb"}, vals)
	})
}

func TestRedis_HIncrBy(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := NewRedis(client.Addr, "").Hincrby("key", "field", 2)
		assert.NotNil(t, err)
		val, err := client.Hincrby("key", "field", 2)
		assert.Nil(t, err)
		assert.Equal(t, 2, val)
		val, err = client.Hincrby("key", "field", 3)
		assert.Nil(t, err)
		assert.Equal(t, 5, val)
	})
}

func TestRedis_Hkeys(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		_, err := NewRedis(client.Addr, "").Hkeys("a")
		assert.NotNil(t, err)
		vals, err := client.Hkeys("a")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"aa", "bb"}, vals)
	})
}

func TestRedis_Hmget(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		_, err := NewRedis(client.Addr, "").Hmget("a", "aa", "bb")
		assert.NotNil(t, err)
		vals, err := client.Hmget("a", "aa", "bb")
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"aaa", "bbb"}, vals)
		vals, err = client.Hmget("a", "aa", "no", "bb")
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"aaa", "", "bbb"}, vals)
	})
}

func TestRedis_Hmset(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.NotNil(t, NewRedis(client.Addr, "").Hmset("a", nil))
		assert.Nil(t, client.Hmset("a", map[string]string{
			"aa": "aaa",
			"bb": "bbb",
		}))
		vals, err := client.Hmget("a", "aa", "bb")
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"aaa", "bbb"}, vals)
	})
}

func TestRedis_Incr(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := NewRedis(client.Addr, "").Incr("a")
		assert.NotNil(t, err)
		val, err := client.Incr("a")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
		val, err = client.Incr("a")
		assert.Nil(t, err)
		assert.Equal(t, int64(2), val)
	})
}

func TestRedis_IncrBy(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := NewRedis(client.Addr, "").Incrby("a", 2)
		assert.NotNil(t, err)
		val, err := client.Incrby("a", 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(2), val)
		val, err = client.Incrby("a", 3)
		assert.Nil(t, err)
		assert.Equal(t, int64(5), val)
	})
}

func TestRedis_Keys(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.Set("key1", "value1")
		assert.Nil(t, err)
		err = client.Set("key2", "value2")
		assert.Nil(t, err)
		_, err = NewRedis(client.Addr, "").Keys("*")
		assert.NotNil(t, err)
		keys, err := client.Keys("*")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"key1", "key2"}, keys)
	})
}

func TestRedis_HyperLogLog(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		client.Ping()
		r := NewRedis(client.Addr, "")
		_, err := r.Pfadd("key1")
		assert.NotNil(t, err)
		_, err = r.Pfcount("*")
		assert.NotNil(t, err)
		err = r.Pfmerge("*")
		assert.NotNil(t, err)
	})
}

func TestRedis_List(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := NewRedis(client.Addr, "").Lpush("key", "value1", "value2")
		assert.NotNil(t, err)
		val, err := client.Lpush("key", "value1", "value2")
		assert.Nil(t, err)
		assert.Equal(t, 2, val)
		_, err = NewRedis(client.Addr, "").Rpush("key", "value3", "value4")
		assert.NotNil(t, err)
		val, err = client.Rpush("key", "value3", "value4")
		assert.Nil(t, err)
		assert.Equal(t, 4, val)
		_, err = NewRedis(client.Addr, "").Llen("key")
		assert.NotNil(t, err)
		val, err = client.Llen("key")
		assert.Nil(t, err)
		assert.Equal(t, 4, val)
		vals, err := client.Lrange("key", 0, 10)
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"value2", "value1", "value3", "value4"}, vals)
		_, err = NewRedis(client.Addr, "").Lpop("key")
		assert.NotNil(t, err)
		v, err := client.Lpop("key")
		assert.Nil(t, err)
		assert.Equal(t, "value2", v)
		val, err = client.Lpush("key", "value1", "value2")
		assert.Nil(t, err)
		assert.Equal(t, 5, val)
		val, err = client.Rpush("key", "value3", "value3")
		assert.Nil(t, err)
		assert.Equal(t, 7, val)
		_, err = NewRedis(client.Addr, "").Lrem("key", 2, "value1")
		assert.NotNil(t, err)
		n, err := client.Lrem("key", 2, "value1")
		assert.Nil(t, err)
		assert.Equal(t, 2, n)
		_, err = NewRedis(client.Addr, "").Lrange("key", 0, 10)
		assert.NotNil(t, err)
		vals, err = client.Lrange("key", 0, 10)
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"value2", "value3", "value4", "value3", "value3"}, vals)
		n, err = client.Lrem("key", -2, "value3")
		assert.Nil(t, err)
		assert.Equal(t, 2, n)
		vals, err = client.Lrange("key", 0, 10)
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"value2", "value3", "value4"}, vals)
	})
}

func TestRedis_Mget(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.Set("key1", "value1")
		assert.Nil(t, err)
		err = client.Set("key2", "value2")
		assert.Nil(t, err)
		_, err = NewRedis(client.Addr, "").Mget("key1", "key0", "key2", "key3")
		assert.NotNil(t, err)
		vals, err := client.Mget("key1", "key0", "key2", "key3")
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"value1", "", "value2", ""}, vals)
	})
}

func TestRedis_SetBit(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := NewRedis(client.Addr, "").SetBit("key", 1, 1)
		assert.NotNil(t, err)
		err = client.SetBit("key", 1, 1)
		assert.Nil(t, err)
	})
}

func TestRedis_GetBit(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.SetBit("key", 2, 1)
		assert.Nil(t, err)
		_, err = NewRedis(client.Addr, "").GetBit("key", 2)
		assert.NotNil(t, err)
		val, err := client.GetBit("key", 2)
		assert.Nil(t, err)
		assert.Equal(t, 1, val)
	})
}

func TestRedis_Persist(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := NewRedis(client.Addr, "").Persist("key")
		assert.NotNil(t, err)
		ok, err := client.Persist("key")
		assert.Nil(t, err)
		assert.False(t, ok)
		err = client.Set("key", "value")
		assert.Nil(t, err)
		ok, err = client.Persist("key")
		assert.Nil(t, err)
		assert.False(t, ok)
		err = NewRedis(client.Addr, "").Expire("key", 5)
		assert.NotNil(t, err)
		err = client.Expire("key", 5)
		assert.Nil(t, err)
		ok, err = client.Persist("key")
		assert.Nil(t, err)
		assert.True(t, ok)
		err = NewRedis(client.Addr, "").Expireat("key", time.Now().Unix()+5)
		assert.NotNil(t, err)
		err = client.Expireat("key", time.Now().Unix()+5)
		assert.Nil(t, err)
		ok, err = client.Persist("key")
		assert.Nil(t, err)
		assert.True(t, ok)
	})
}

func TestRedis_Ping(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		ok := client.Ping()
		assert.True(t, ok)
	})
}

func TestRedis_Scan(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.Set("key1", "value1")
		assert.Nil(t, err)
		err = client.Set("key2", "value2")
		assert.Nil(t, err)
		_, _, err = NewRedis(client.Addr, "").Scan(0, "*", 100)
		assert.NotNil(t, err)
		keys, _, err := client.Scan(0, "*", 100)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"key1", "key2"}, keys)
	})
}

func TestRedis_Sscan(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		key := "list"
		var list []string
		for i := 0; i < 1550; i++ {
			list = append(list, randomStr(i))
		}
		lens, err := client.Sadd(key, list)
		assert.Nil(t, err)
		assert.Equal(t, lens, 1550)

		var cursor uint64 = 0
		sum := 0
		for {
			_, _, err := NewRedis(client.Addr, "").Sscan(key, cursor, "", 100)
			assert.NotNil(t, err)
			keys, next, err := client.Sscan(key, cursor, "", 100)
			assert.Nil(t, err)
			sum += len(keys)
			if next == 0 {
				break
			}
			cursor = next
		}

		assert.Equal(t, sum, 1550)
		_, err = NewRedis(client.Addr, "").Del(key)
		assert.NotNil(t, err)
		_, err = client.Del(key)
		assert.Nil(t, err)
	})
}

func TestRedis_Set(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := NewRedis(client.Addr, "").Sadd("key", 1, 2, 3, 4)
		assert.NotNil(t, err)
		num, err := client.Sadd("key", 1, 2, 3, 4)
		assert.Nil(t, err)
		assert.Equal(t, 4, num)
		_, err = NewRedis(client.Addr, "").Scard("key")
		assert.NotNil(t, err)
		val, err := client.Scard("key")
		assert.Nil(t, err)
		assert.Equal(t, int64(4), val)
		_, err = NewRedis(client.Addr, "").Sismember("key", 2)
		assert.NotNil(t, err)
		ok, err := client.Sismember("key", 2)
		assert.Nil(t, err)
		assert.True(t, ok)
		_, err = NewRedis(client.Addr, "").Srem("key", 3, 4)
		assert.NotNil(t, err)
		num, err = client.Srem("key", 3, 4)
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		_, err = NewRedis(client.Addr, "").Smembers("key")
		assert.NotNil(t, err)
		vals, err := client.Smembers("key")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"1", "2"}, vals)
		_, err = NewRedis(client.Addr, "").Srandmember("key", 1)
		assert.NotNil(t, err)
		members, err := client.Srandmember("key", 1)
		assert.Nil(t, err)
		assert.Len(t, members, 1)
		assert.Contains(t, []string{"1", "2"}, members[0])
		_, err = NewRedis(client.Addr, "").Spop("key")
		assert.NotNil(t, err)
		member, err := client.Spop("key")
		assert.Nil(t, err)
		assert.Contains(t, []string{"1", "2"}, member)
		_, err = NewRedis(client.Addr, "").Smembers("key")
		assert.NotNil(t, err)
		vals, err = client.Smembers("key")
		assert.Nil(t, err)
		assert.NotContains(t, vals, member)
		_, err = NewRedis(client.Addr, "").Sadd("key1", 1, 2, 3, 4)
		assert.NotNil(t, err)
		num, err = client.Sadd("key1", 1, 2, 3, 4)
		assert.Nil(t, err)
		assert.Equal(t, 4, num)
		num, err = client.Sadd("key2", 2, 3, 4, 5)
		assert.Nil(t, err)
		assert.Equal(t, 4, num)
		_, err = NewRedis(client.Addr, "").Sunion("key1", "key2")
		assert.NotNil(t, err)
		vals, err = client.Sunion("key1", "key2")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"1", "2", "3", "4", "5"}, vals)
		_, err = NewRedis(client.Addr, "").Sunionstore("key3", "key1", "key2")
		assert.NotNil(t, err)
		num, err = client.Sunionstore("key3", "key1", "key2")
		assert.Nil(t, err)
		assert.Equal(t, 5, num)
		_, err = NewRedis(client.Addr, "").Sdiff("key1", "key2")
		assert.NotNil(t, err)
		vals, err = client.Sdiff("key1", "key2")
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"1"}, vals)
		_, err = NewRedis(client.Addr, "").Sdiffstore("key4", "key1", "key2")
		assert.NotNil(t, err)
		num, err = client.Sdiffstore("key4", "key1", "key2")
		assert.Nil(t, err)
		assert.Equal(t, 1, num)
	})
}

func TestRedis_SetGetDel(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := NewRedis(client.Addr, "").Set("hello", "world")
		assert.NotNil(t, err)
		err = client.Set("hello", "world")
		assert.Nil(t, err)
		_, err = NewRedis(client.Addr, "").Get("hello")
		assert.NotNil(t, err)
		val, err := client.Get("hello")
		assert.Nil(t, err)
		assert.Equal(t, "world", val)
		ret, err := client.Del("hello")
		assert.Nil(t, err)
		assert.Equal(t, 1, ret)
	})
}

func TestRedis_SetExNx(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := NewRedis(client.Addr, "").Setex("hello", "world", 5)
		assert.NotNil(t, err)
		err = client.Setex("hello", "world", 5)
		assert.Nil(t, err)
		_, err = NewRedis(client.Addr, "").Setnx("hello", "newworld")
		assert.NotNil(t, err)
		ok, err := client.Setnx("hello", "newworld")
		assert.Nil(t, err)
		assert.False(t, ok)
		ok, err = client.Setnx("newhello", "newworld")
		assert.Nil(t, err)
		assert.True(t, ok)
		val, err := client.Get("hello")
		assert.Nil(t, err)
		assert.Equal(t, "world", val)
		val, err = client.Get("newhello")
		assert.Nil(t, err)
		assert.Equal(t, "newworld", val)
		ttl, err := client.Ttl("hello")
		assert.Nil(t, err)
		assert.True(t, ttl > 0)
		_, err = NewRedis(client.Addr, "").SetnxEx("newhello", "newworld", 5)
		assert.NotNil(t, err)
		ok, err = client.SetnxEx("newhello", "newworld", 5)
		assert.Nil(t, err)
		assert.False(t, ok)
		num, err := client.Del("newhello")
		assert.Nil(t, err)
		assert.Equal(t, 1, num)
		ok, err = client.SetnxEx("newhello", "newworld", 5)
		assert.Nil(t, err)
		assert.True(t, ok)
		val, err = client.Get("newhello")
		assert.Nil(t, err)
		assert.Equal(t, "newworld", val)
	})
}

func TestRedis_SetGetDelHashField(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.Hset("key", "field", "value")
		assert.Nil(t, err)
		_, err = NewRedis(client.Addr, "").Hget("key", "field")
		assert.NotNil(t, err)
		val, err := client.Hget("key", "field")
		assert.Nil(t, err)
		assert.Equal(t, "value", val)
		_, err = NewRedis(client.Addr, "").Hexists("key", "field")
		assert.NotNil(t, err)
		ok, err := client.Hexists("key", "field")
		assert.Nil(t, err)
		assert.True(t, ok)
		_, err = NewRedis(client.Addr, "").Hdel("key", "field")
		assert.NotNil(t, err)
		ret, err := client.Hdel("key", "field")
		assert.Nil(t, err)
		assert.True(t, ret)
		ok, err = client.Hexists("key", "field")
		assert.Nil(t, err)
		assert.False(t, ok)
	})
}

func TestRedis_SortedSet(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		ok, err := client.Zadd("key", 1, "value1")
		assert.Nil(t, err)
		assert.True(t, ok)
		ok, err = client.Zadd("key", 2, "value1")
		assert.Nil(t, err)
		assert.False(t, ok)
		val, err := client.Zscore("key", "value1")
		assert.Nil(t, err)
		assert.Equal(t, int64(2), val)
		_, err = NewRedis(client.Addr, "").Zincrby("key", 3, "value1")
		assert.NotNil(t, err)
		val, err = client.Zincrby("key", 3, "value1")
		assert.Nil(t, err)
		assert.Equal(t, int64(5), val)
		_, err = NewRedis(client.Addr, "").Zscore("key", "value1")
		assert.NotNil(t, err)
		val, err = client.Zscore("key", "value1")
		assert.Nil(t, err)
		assert.Equal(t, int64(5), val)
		val, err = NewRedis(client.Addr, "").Zadds("key")
		assert.NotNil(t, err)
		val, err = client.Zadds("key", Pair{
			Key:   "value2",
			Score: 6,
		}, Pair{
			Key:   "value3",
			Score: 7,
		})
		assert.Nil(t, err)
		assert.Equal(t, int64(2), val)
		pairs, err := NewRedis(client.Addr, "").ZRevRangeWithScores("key", 1, 3)
		assert.NotNil(t, err)
		pairs, err = client.ZRevRangeWithScores("key", 1, 3)
		assert.Nil(t, err)
		assert.EqualValues(t, []Pair{
			{
				Key:   "value2",
				Score: 6,
			},
			{
				Key:   "value1",
				Score: 5,
			},
		}, pairs)
		rank, err := client.Zrank("key", "value2")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), rank)
		_, err = NewRedis(client.Addr, "").Zrank("key", "value4")
		assert.NotNil(t, err)
		_, err = client.Zrank("key", "value4")
		assert.Equal(t, Nil, err)
		_, err = NewRedis(client.Addr, "").Zrem("key", "value2", "value3")
		assert.NotNil(t, err)
		num, err := client.Zrem("key", "value2", "value3")
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		ok, err = client.Zadd("key", 6, "value2")
		assert.Nil(t, err)
		assert.True(t, ok)
		ok, err = client.Zadd("key", 7, "value3")
		assert.Nil(t, err)
		assert.True(t, ok)
		ok, err = client.Zadd("key", 8, "value4")
		assert.Nil(t, err)
		assert.True(t, ok)
		_, err = NewRedis(client.Addr, "").Zremrangebyscore("key", 6, 7)
		assert.NotNil(t, err)
		num, err = client.Zremrangebyscore("key", 6, 7)
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		ok, err = client.Zadd("key", 6, "value2")
		assert.Nil(t, err)
		assert.True(t, ok)
		_, err = NewRedis(client.Addr, "").Zadd("key", 7, "value3")
		assert.NotNil(t, err)
		ok, err = client.Zadd("key", 7, "value3")
		assert.Nil(t, err)
		assert.True(t, ok)
		_, err = NewRedis(client.Addr, "").Zcount("key", 6, 7)
		assert.NotNil(t, err)
		num, err = client.Zcount("key", 6, 7)
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		_, err = NewRedis(client.Addr, "").Zremrangebyrank("key", 1, 2)
		assert.NotNil(t, err)
		num, err = client.Zremrangebyrank("key", 1, 2)
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		_, err = NewRedis(client.Addr, "").Zcard("key")
		assert.NotNil(t, err)
		card, err := client.Zcard("key")
		assert.Nil(t, err)
		assert.Equal(t, 2, card)
		_, err = NewRedis(client.Addr, "").Zrange("key", 0, -1)
		assert.NotNil(t, err)
		vals, err := client.Zrange("key", 0, -1)
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"value1", "value4"}, vals)
		_, err = NewRedis(client.Addr, "").Zrevrange("key", 0, -1)
		assert.NotNil(t, err)
		vals, err = client.Zrevrange("key", 0, -1)
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"value4", "value1"}, vals)
		_, err = NewRedis(client.Addr, "").ZrangeWithScores("key", 0, -1)
		assert.NotNil(t, err)
		pairs, err = client.ZrangeWithScores("key", 0, -1)
		assert.Nil(t, err)
		assert.EqualValues(t, []Pair{
			{
				Key:   "value1",
				Score: 5,
			},
			{
				Key:   "value4",
				Score: 8,
			},
		}, pairs)
		_, err = NewRedis(client.Addr, "").ZrangebyscoreWithScores("key", 5, 8)
		assert.NotNil(t, err)
		pairs, err = client.ZrangebyscoreWithScores("key", 5, 8)
		assert.Nil(t, err)
		assert.EqualValues(t, []Pair{
			{
				Key:   "value1",
				Score: 5,
			},
			{
				Key:   "value4",
				Score: 8,
			},
		}, pairs)
		_, err = NewRedis(client.Addr, "").ZrangebyscoreWithScoresAndLimit(
			"key", 5, 8, 1, 1)
		assert.NotNil(t, err)
		pairs, err = client.ZrangebyscoreWithScoresAndLimit("key", 5, 8, 1, 1)
		assert.Nil(t, err)
		assert.EqualValues(t, []Pair{
			{
				Key:   "value4",
				Score: 8,
			},
		}, pairs)
		pairs, err = client.ZrangebyscoreWithScoresAndLimit("key", 5, 8, 1, 0)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(pairs))
		_, err = NewRedis(client.Addr, "").ZrevrangebyscoreWithScores("key", 5, 8)
		assert.NotNil(t, err)
		pairs, err = client.ZrevrangebyscoreWithScores("key", 5, 8)
		assert.Nil(t, err)
		assert.EqualValues(t, []Pair{
			{
				Key:   "value4",
				Score: 8,
			},
			{
				Key:   "value1",
				Score: 5,
			},
		}, pairs)
		_, err = NewRedis(client.Addr, "").ZrevrangebyscoreWithScoresAndLimit(
			"key", 5, 8, 1, 1)
		assert.NotNil(t, err)
		pairs, err = client.ZrevrangebyscoreWithScoresAndLimit("key", 5, 8, 1, 1)
		assert.Nil(t, err)
		assert.EqualValues(t, []Pair{
			{
				Key:   "value1",
				Score: 5,
			},
		}, pairs)
		pairs, err = client.ZrevrangebyscoreWithScoresAndLimit("key", 5, 8, 1, 0)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(pairs))
	})
}

func TestRedis_Pipelined(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.NotNil(t, NewRedis(client.Addr, "").Pipelined(func(pipeliner Pipeliner) error {
			return nil
		}))
		err := client.Pipelined(
			func(pipe Pipeliner) error {
				pipe.Incr("pipelined_counter")
				pipe.Expire("pipelined_counter", time.Hour)
				pipe.ZAdd("zadd", Z{Score: 12, Member: "zadd"})
				return nil
			},
		)
		assert.Nil(t, err)
		_, err = NewRedis(client.Addr, "").Ttl("pipelined_counter")
		assert.NotNil(t, err)
		ttl, err := client.Ttl("pipelined_counter")
		assert.Nil(t, err)
		assert.Equal(t, 3600, ttl)
		value, err := client.Get("pipelined_counter")
		assert.Nil(t, err)
		assert.Equal(t, "1", value)
		score, err := client.Zscore("zadd", "zadd")
		assert.Nil(t, err)
		assert.Equal(t, int64(12), score)
	})
}

func TestRedisString(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		client.Ping()
		_, err := getRedis(NewRedis(client.Addr, ClusterType))
		assert.Nil(t, err)
		assert.Equal(t, client.Addr, client.String())
		assert.NotNil(t, NewRedis(client.Addr, "").Ping())
	})
}

func TestRedisScriptLoad(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		client.Ping()
		_, err := NewRedis(client.Addr, "").scriptLoad("foo")
		assert.NotNil(t, err)
		_, err = client.scriptLoad("foo")
		assert.NotNil(t, err)
	})
}

func TestRedisToPairs(t *testing.T) {
	pairs := toPairs([]red.Z{
		{
			Member: 1,
			Score:  1,
		},
		{
			Member: 2,
			Score:  2,
		},
	})
	assert.EqualValues(t, []Pair{
		{
			Key:   "1",
			Score: 1,
		},
		{
			Key:   "2",
			Score: 2,
		},
	}, pairs)
}

func TestRedisToStrings(t *testing.T) {
	vals := toStrings([]interface{}{1, 2})
	assert.EqualValues(t, []string{"1", "2"}, vals)
}

func TestRedisBlpop(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		client.Ping()
		var node mockedNode
		_, err := client.Blpop(nil, "foo")
		assert.NotNil(t, err)
		_, err = client.Blpop(node, "foo")
		assert.NotNil(t, err)
	})
}

func TestRedisBlpopEx(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		client.Ping()
		var node mockedNode
		_, _, err := client.BlpopEx(nil, "foo")
		assert.NotNil(t, err)
		_, _, err = client.BlpopEx(node, "foo")
		assert.NotNil(t, err)
	})
}

func runOnRedis(t *testing.T, fn func(client *Redis)) {
	s, err := miniredis.Run()
	assert.Nil(t, err)
	defer func() {
		client, err := clientManager.GetResource(s.Addr(), func() (io.Closer, error) {
			return nil, errors.New("should already exist")
		})
		if err != nil {
			t.Error(err)
		}

		if client != nil {
			client.Close()
		}
	}()

	fn(NewRedis(s.Addr(), NodeType))
}

type mockedNode struct {
	RedisNode
}

func (n mockedNode) BLPop(timeout time.Duration, keys ...string) *red.StringSliceCmd {
	return red.NewStringSliceCmd("foo", "bar")
}
