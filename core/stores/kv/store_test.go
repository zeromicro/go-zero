package kv

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/stores/internal"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/stringx"
)

var s1, _ = miniredis.Run()
var s2, _ = miniredis.Run()

func TestRedis_Exists(t *testing.T) {
	runOnCluster(t, func(client Store) {
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
	runOnCluster(t, func(client Store) {
		_, err := client.Eval(`redis.call("EXISTS", KEYS[1])`, "notexist")
		assert.Equal(t, redis.Nil, err)
		err = client.Set("key1", "value1")
		assert.Nil(t, err)
		_, err = client.Eval(`redis.call("EXISTS", KEYS[1])`, "key1")
		assert.Equal(t, redis.Nil, err)
		val, err := client.Eval(`return redis.call("EXISTS", KEYS[1])`, "key1")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
	})
}

func TestRedis_Hgetall(t *testing.T) {
	runOnCluster(t, func(client Store) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		vals, err := client.Hgetall("a")
		assert.Nil(t, err)
		assert.EqualValues(t, map[string]string{
			"aa": "aaa",
			"bb": "bbb",
		}, vals)
	})
}

func TestRedis_Hvals(t *testing.T) {
	runOnCluster(t, func(client Store) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		vals, err := client.Hvals("a")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"aaa", "bbb"}, vals)
	})
}

func TestRedis_Hsetnx(t *testing.T) {
	runOnCluster(t, func(client Store) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
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
	runOnCluster(t, func(client Store) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
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
	runOnCluster(t, func(client Store) {
		val, err := client.Hincrby("key", "field", 2)
		assert.Nil(t, err)
		assert.Equal(t, 2, val)
		val, err = client.Hincrby("key", "field", 3)
		assert.Nil(t, err)
		assert.Equal(t, 5, val)
	})
}

func TestRedis_Hkeys(t *testing.T) {
	runOnCluster(t, func(client Store) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		vals, err := client.Hkeys("a")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"aa", "bb"}, vals)
	})
}

func TestRedis_Hmget(t *testing.T) {
	runOnCluster(t, func(client Store) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		vals, err := client.Hmget("a", "aa", "bb")
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"aaa", "bbb"}, vals)
		vals, err = client.Hmget("a", "aa", "no", "bb")
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"aaa", "", "bbb"}, vals)
	})
}

func TestRedis_Hmset(t *testing.T) {
	runOnCluster(t, func(client Store) {
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
	runOnCluster(t, func(client Store) {
		val, err := client.Incr("a")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
		val, err = client.Incr("a")
		assert.Nil(t, err)
		assert.Equal(t, int64(2), val)
	})
}

func TestRedis_IncrBy(t *testing.T) {
	runOnCluster(t, func(client Store) {
		val, err := client.Incrby("a", 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(2), val)
		val, err = client.Incrby("a", 3)
		assert.Nil(t, err)
		assert.Equal(t, int64(5), val)
	})
}

func TestRedis_List(t *testing.T) {
	runOnCluster(t, func(client Store) {
		val, err := client.Lpush("key", "value1", "value2")
		assert.Nil(t, err)
		assert.Equal(t, 2, val)
		val, err = client.Rpush("key", "value3", "value4")
		assert.Nil(t, err)
		assert.Equal(t, 4, val)
		val, err = client.Llen("key")
		assert.Nil(t, err)
		assert.Equal(t, 4, val)
		vals, err := client.Lrange("key", 0, 10)
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"value2", "value1", "value3", "value4"}, vals)
		v, err := client.Lpop("key")
		assert.Nil(t, err)
		assert.Equal(t, "value2", v)
		val, err = client.Lpush("key", "value1", "value2")
		assert.Nil(t, err)
		assert.Equal(t, 5, val)
		val, err = client.Rpush("key", "value3", "value3")
		assert.Nil(t, err)
		assert.Equal(t, 7, val)
		n, err := client.Lrem("key", 2, "value1")
		assert.Nil(t, err)
		assert.Equal(t, 2, n)
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

func TestRedis_Persist(t *testing.T) {
	runOnCluster(t, func(client Store) {
		ok, err := client.Persist("key")
		assert.Nil(t, err)
		assert.False(t, ok)
		err = client.Set("key", "value")
		assert.Nil(t, err)
		ok, err = client.Persist("key")
		assert.Nil(t, err)
		assert.False(t, ok)
		err = client.Expire("key", 5)
		assert.Nil(t, err)
		ok, err = client.Persist("key")
		assert.Nil(t, err)
		assert.True(t, ok)
		err = client.Expireat("key", time.Now().Unix()+5)
		assert.Nil(t, err)
		ok, err = client.Persist("key")
		assert.Nil(t, err)
		assert.True(t, ok)
	})
}

func TestRedis_Sscan(t *testing.T) {
	runOnCluster(t, func(client Store) {
		key := "list"
		var list []string
		for i := 0; i < 1550; i++ {
			list = append(list, stringx.Randn(i))
		}
		lens, err := client.Sadd(key, list)
		assert.Nil(t, err)
		assert.Equal(t, lens, 1550)

		var cursor uint64 = 0
		sum := 0
		for {
			keys, next, err := client.Sscan(key, cursor, "", 100)
			assert.Nil(t, err)
			sum += len(keys)
			if next == 0 {
				break
			}
			cursor = next
		}

		assert.Equal(t, sum, 1550)
		_, err = client.Del(key)
		assert.Nil(t, err)
	})
}

func TestRedis_Set(t *testing.T) {
	runOnCluster(t, func(client Store) {
		num, err := client.Sadd("key", 1, 2, 3, 4)
		assert.Nil(t, err)
		assert.Equal(t, 4, num)
		val, err := client.Scard("key")
		assert.Nil(t, err)
		assert.Equal(t, int64(4), val)
		ok, err := client.Sismember("key", 2)
		assert.Nil(t, err)
		assert.True(t, ok)
		num, err = client.Srem("key", 3, 4)
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		vals, err := client.Smembers("key")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"1", "2"}, vals)
		members, err := client.Srandmember("key", 1)
		assert.Nil(t, err)
		assert.Len(t, members, 1)
		assert.Contains(t, []string{"1", "2"}, members[0])
		member, err := client.Spop("key")
		assert.Nil(t, err)
		assert.Contains(t, []string{"1", "2"}, member)
		vals, err = client.Smembers("key")
		assert.Nil(t, err)
		assert.NotContains(t, vals, member)
		num, err = client.Sadd("key1", 1, 2, 3, 4)
		assert.Nil(t, err)
		assert.Equal(t, 4, num)
		num, err = client.Sadd("key2", 2, 3, 4, 5)
		assert.Nil(t, err)
		assert.Equal(t, 4, num)
	})
}

func TestRedis_SetGetDel(t *testing.T) {
	runOnCluster(t, func(client Store) {
		err := client.Set("hello", "world")
		assert.Nil(t, err)
		val, err := client.Get("hello")
		assert.Nil(t, err)
		assert.Equal(t, "world", val)
		ret, err := client.Del("hello")
		assert.Nil(t, err)
		assert.Equal(t, 1, ret)
	})
}

func TestRedis_SetExNx(t *testing.T) {
	runOnCluster(t, func(client Store) {
		err := client.Setex("hello", "world", 5)
		assert.Nil(t, err)
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
	runOnCluster(t, func(client Store) {
		err := client.Hset("key", "field", "value")
		assert.Nil(t, err)
		val, err := client.Hget("key", "field")
		assert.Nil(t, err)
		assert.Equal(t, "value", val)
		ok, err := client.Hexists("key", "field")
		assert.Nil(t, err)
		assert.True(t, ok)
		ret, err := client.Hdel("key", "field")
		assert.Nil(t, err)
		assert.True(t, ret)
		ok, err = client.Hexists("key", "field")
		assert.Nil(t, err)
		assert.False(t, ok)
	})
}

func TestRedis_SortedSet(t *testing.T) {
	runOnCluster(t, func(client Store) {
		ok, err := client.Zadd("key", 1, "value1")
		assert.Nil(t, err)
		assert.True(t, ok)
		ok, err = client.Zadd("key", 2, "value1")
		assert.Nil(t, err)
		assert.False(t, ok)
		val, err := client.Zscore("key", "value1")
		assert.Nil(t, err)
		assert.Equal(t, int64(2), val)
		val, err = client.Zincrby("key", 3, "value1")
		assert.Nil(t, err)
		assert.Equal(t, int64(5), val)
		val, err = client.Zscore("key", "value1")
		assert.Nil(t, err)
		assert.Equal(t, int64(5), val)
		ok, err = client.Zadd("key", 6, "value2")
		assert.Nil(t, err)
		assert.True(t, ok)
		ok, err = client.Zadd("key", 7, "value3")
		assert.Nil(t, err)
		assert.True(t, ok)
		rank, err := client.Zrank("key", "value2")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), rank)
		_, err = client.Zrank("key", "value4")
		assert.Equal(t, redis.Nil, err)
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
		num, err = client.Zremrangebyscore("key", 6, 7)
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		ok, err = client.Zadd("key", 6, "value2")
		assert.Nil(t, err)
		assert.True(t, ok)
		ok, err = client.Zadd("key", 7, "value3")
		assert.Nil(t, err)
		assert.True(t, ok)
		num, err = client.Zcount("key", 6, 7)
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		num, err = client.Zremrangebyrank("key", 1, 2)
		assert.Nil(t, err)
		assert.Equal(t, 2, num)
		card, err := client.Zcard("key")
		assert.Nil(t, err)
		assert.Equal(t, 2, card)
		vals, err := client.Zrange("key", 0, -1)
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"value1", "value4"}, vals)
		vals, err = client.Zrevrange("key", 0, -1)
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"value4", "value1"}, vals)
		pairs, err := client.ZrangeWithScores("key", 0, -1)
		assert.Nil(t, err)
		assert.EqualValues(t, []redis.Pair{
			{
				Key:   "value1",
				Score: 5,
			},
			{
				Key:   "value4",
				Score: 8,
			},
		}, pairs)
		pairs, err = client.ZrangebyscoreWithScores("key", 5, 8)
		assert.Nil(t, err)
		assert.EqualValues(t, []redis.Pair{
			{
				Key:   "value1",
				Score: 5,
			},
			{
				Key:   "value4",
				Score: 8,
			},
		}, pairs)
		pairs, err = client.ZrangebyscoreWithScoresAndLimit("key", 5, 8, 1, 1)
		assert.Nil(t, err)
		assert.EqualValues(t, []redis.Pair{
			{
				Key:   "value4",
				Score: 8,
			},
		}, pairs)
		pairs, err = client.ZrevrangebyscoreWithScores("key", 5, 8)
		assert.Nil(t, err)
		assert.EqualValues(t, []redis.Pair{
			{
				Key:   "value4",
				Score: 8,
			},
			{
				Key:   "value1",
				Score: 5,
			},
		}, pairs)
		pairs, err = client.ZrevrangebyscoreWithScoresAndLimit("key", 5, 8, 1, 1)
		assert.Nil(t, err)
		assert.EqualValues(t, []redis.Pair{
			{
				Key:   "value1",
				Score: 5,
			},
		}, pairs)
	})
}

func runOnCluster(t *testing.T, fn func(cluster Store)) {
	s1.FlushAll()
	s2.FlushAll()

	store := NewStore([]internal.NodeConf{
		{
			RedisConf: redis.RedisConf{
				Host: s1.Addr(),
				Type: redis.NodeType,
			},
			Weight: 100,
		},
		{
			RedisConf: redis.RedisConf{
				Host: s2.Addr(),
				Type: redis.NodeType,
			},
			Weight: 100,
		},
	})

	fn(store)
}
