package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"strconv"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	red "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stringx"
)

type myHook struct {
	includePing bool
}

var _ red.Hook = myHook{}

func (m myHook) DialHook(next red.DialHook) red.DialHook {
	return next
}

func (m myHook) ProcessPipelineHook(next red.ProcessPipelineHook) red.ProcessPipelineHook {
	return next
}

func (m myHook) ProcessHook(next red.ProcessHook) red.ProcessHook {
	return func(ctx context.Context, cmd red.Cmder) error {
		// skip ping cmd
		if cmd.Name() == "ping" && !m.includePing {
			return next(ctx, cmd)
		}
		return errors.New("durationHook error")
	}
}

func TestNewRedis(t *testing.T) {
	r1, err := miniredis.Run()
	assert.NoError(t, err)
	defer r1.Close()

	r2, err := miniredis.Run()
	assert.NoError(t, err)
	defer r2.Close()
	r2.SetError("mock")

	tests := []struct {
		name string
		RedisConf
		ok       bool
		redisErr bool
	}{
		{
			name: "missing host",
			RedisConf: RedisConf{
				Host: "",
				Type: NodeType,
				Pass: "",
			},
			ok: false,
		},
		{
			name: "missing type",
			RedisConf: RedisConf{
				Host: "localhost:6379",
				Type: "",
				Pass: "",
			},
			ok: false,
		},
		{
			name: "ok",
			RedisConf: RedisConf{
				Host: r1.Addr(),
				Type: NodeType,
				Pass: "",
			},
			ok: true,
		},
		{
			name: "ok",
			RedisConf: RedisConf{
				Host: r1.Addr(),
				Type: ClusterType,
				Pass: "",
			},
			ok: true,
		},
		{
			name: "password",
			RedisConf: RedisConf{
				Host: r1.Addr(),
				Type: NodeType,
				Pass: "pw",
			},
			ok: true,
		},
		{
			name: "tls",
			RedisConf: RedisConf{
				Host: r1.Addr(),
				Type: NodeType,
				Tls:  true,
			},
			ok: true,
		},
		{
			name: "node error",
			RedisConf: RedisConf{
				Host: r2.Addr(),
				Type: NodeType,
				Pass: "",
			},
			ok:       true,
			redisErr: true,
		},
		{
			name: "cluster error",
			RedisConf: RedisConf{
				Host: r2.Addr(),
				Type: ClusterType,
				Pass: "",
			},
			ok:       true,
			redisErr: true,
		},
	}

	for _, test := range tests {
		t.Run(stringx.RandId(), func(t *testing.T) {
			rds, err := NewRedis(test.RedisConf)
			if test.ok {
				if test.redisErr {
					assert.Error(t, err)
					assert.Nil(t, rds)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, rds)
				}
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestRedis_NonBlock(t *testing.T) {
	logx.Disable()

	t.Run("nonBlock true", func(t *testing.T) {
		s := miniredis.RunT(t)
		// use durationHook to simulate redis ping error
		_, err := NewRedis(RedisConf{
			Host:     s.Addr(),
			NonBlock: true,
			Type:     NodeType,
		}, WithHook(myHook{includePing: true}))
		assert.NoError(t, err)
	})

	t.Run("nonBlock false", func(t *testing.T) {
		s := miniredis.RunT(t)
		_, err := NewRedis(RedisConf{
			Host:     s.Addr(),
			NonBlock: false,
			Type:     NodeType,
		}, WithHook(myHook{includePing: true}))
		assert.ErrorContains(t, err, "redis connect error")
	})
}

func TestRedis_Decr(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Decr("a")
		assert.NotNil(t, err)
		val, err := client.Decr("a")
		assert.Nil(t, err)
		assert.Equal(t, int64(-1), val)
		val, err = client.Decr("a")
		assert.Nil(t, err)
		assert.Equal(t, int64(-2), val)
	})
}

func TestRedis_DecrBy(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Decrby("a", 2)
		assert.NotNil(t, err)
		val, err := client.Decrby("a", 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(-2), val)
		val, err = client.Decrby("a", 3)
		assert.Nil(t, err)
		assert.Equal(t, int64(-5), val)
	})
}

func TestRedis_Exists(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Exists("a")
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

func TestRedisTLS_Exists(t *testing.T) {
	runOnRedisTLS(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Exists("a")
		assert.NotNil(t, err)
		ok, err := client.Exists("a")
		assert.NotNil(t, err)
		assert.False(t, ok)
		assert.NotNil(t, client.Set("a", "b"))
		ok, err = client.Exists("a")
		assert.NotNil(t, err)
		assert.False(t, ok)
	})
}

func TestRedis_ExistsMany(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		// Attempt to create a new Redis instance with an incorrect type and call ExistsMany
		_, err := newRedis(client.Addr, badType()).ExistsMany("key1", "key2")
		assert.NotNil(t, err)

		// Check if key1 and key2 exist, expecting that they do not
		val, err := client.ExistsMany("key1", "key2")
		assert.Nil(t, err)
		assert.Equal(t, int64(0), val)

		// Set the value for key1 and check if key1 exists
		assert.Nil(t, client.Set("key1", "value1"))
		val, err = client.ExistsMany("key1")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)

		// Set the value for key2 and check if key1 and key2 exist
		assert.Nil(t, client.Set("key2", "value2"))
		val, err = client.ExistsMany("key1", "key2")
		assert.Nil(t, err)
		assert.Equal(t, int64(2), val)

		// Check if key1, key2, and a non-existent key3 exist, expecting that key1 and key2 do
		val, err = client.ExistsMany("key1", "key2", "key3")
		assert.Nil(t, err)
		assert.Equal(t, int64(2), val)
	})
}

func TestRedis_Eval(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Eval(`redis.call("EXISTS", KEYS[1])`, []string{"notexist"})
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

func TestRedis_ScriptRun(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		sc := NewScript(`redis.call("EXISTS", KEYS[1])`)
		sc2 := NewScript(`return redis.call("EXISTS", KEYS[1])`)
		_, err := newRedis(client.Addr, badType()).ScriptRun(sc, []string{"notexist"})
		assert.NotNil(t, err)
		_, err = client.ScriptRun(sc, []string{"notexist"})
		assert.Equal(t, Nil, err)
		err = client.Set("key1", "value1")
		assert.Nil(t, err)
		_, err = client.ScriptRun(sc, []string{"key1"})
		assert.Equal(t, Nil, err)
		val, err := client.ScriptRun(sc2, []string{"key1"})
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
	})
}

func TestRedis_GeoHash(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := client.GeoHash("parent", "child1", "child2")
		assert.Error(t, err)
		_, err = newRedis(client.Addr, badType()).GeoHash("parent", "child1", "child2")
		assert.Error(t, err)
	})
}

func TestRedis_Hgetall(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		_, err := newRedis(client.Addr, badType()).Hgetall("a")
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
		assert.NotNil(t, newRedis(client.Addr, badType()).Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "aa", "aaa"))
		assert.Nil(t, client.Hset("a", "bb", "bbb"))
		_, err := newRedis(client.Addr, badType()).Hvals("a")
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
		_, err := newRedis(client.Addr, badType()).Hsetnx("a", "bb", "ccc")
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
		_, err := newRedis(client.Addr, badType()).Hlen("a")
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
		_, err := newRedis(client.Addr, badType()).Hincrby("key", "field", 2)
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
		_, err := newRedis(client.Addr, badType()).Hkeys("a")
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
		_, err := newRedis(client.Addr, badType()).Hmget("a", "aa", "bb")
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
		assert.NotNil(t, newRedis(client.Addr, badType()).Hmset("a", nil))
		assert.Nil(t, client.Hmset("a", map[string]string{
			"aa": "aaa",
			"bb": "bbb",
		}))
		vals, err := client.Hmget("a", "aa", "bb")
		assert.Nil(t, err)
		assert.EqualValues(t, []string{"aaa", "bbb"}, vals)
	})
}

func TestRedis_Hscan(t *testing.T) {
	t.Run("scan", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			key := "hash:test"
			fieldsAndValues := make(map[string]string)
			for i := 0; i < 1550; i++ {
				fieldsAndValues["filed_"+strconv.Itoa(i)] = stringx.Randn(i)
			}
			err := client.Hmset(key, fieldsAndValues)
			assert.Nil(t, err)

			var cursor uint64 = 0
			sum := 0
			for {
				_, _, err := newRedis(client.Addr, badType()).Hscan(key, cursor, "*", 100)
				assert.NotNil(t, err)
				reMap, next, err := client.Hscan(key, cursor, "*", 100)
				assert.Nil(t, err)
				sum += len(reMap)
				if next == 0 {
					break
				}
				cursor = next
			}

			assert.Equal(t, sum, 3100)
			_, err = newRedis(client.Addr, badType()).Del(key)
			assert.Error(t, err)
			_, err = client.Del(key)
			assert.NoError(t, err)
		})
	})

	t.Run("scan with error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Del("hash:test")
			assert.Error(t, err)
		})
	})
}

func TestRedis_Incr(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Incr("a")
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
		_, err := newRedis(client.Addr, badType()).Incrby("a", 2)
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
		_, err = newRedis(client.Addr, badType()).Keys("*")
		assert.NotNil(t, err)
		keys, err := client.Keys("*")
		assert.Nil(t, err)
		assert.ElementsMatch(t, []string{"key1", "key2"}, keys)
	})
}

func TestRedis_HyperLogLog(t *testing.T) {
	t.Run("hyperloglog", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			client.Ping()
			r := newRedis(client.Addr)
			_, err := newRedis(client.Addr, badType()).Pfadd("key1", "val1")
			assert.Error(t, err)
			ok, err := r.Pfadd("key1", "val1")
			assert.Nil(t, err)
			assert.True(t, ok)
			_, err = newRedis(client.Addr, badType()).Pfcount("key1")
			assert.Error(t, err)
			val, err := r.Pfcount("key1")
			assert.Nil(t, err)
			assert.Equal(t, int64(1), val)
			ok, err = r.Pfadd("key2", "val2")
			assert.Nil(t, err)
			assert.True(t, ok)
			val, err = r.Pfcount("key2")
			assert.Nil(t, err)
			assert.Equal(t, int64(1), val)
			err = newRedis(client.Addr, badType()).Pfmerge("key3", "key1", "key2")
			assert.Error(t, err)
			err = r.Pfmerge("key1", "key2")
			assert.Nil(t, err)
			val, err = r.Pfcount("key1")
			assert.Nil(t, err)
			assert.Equal(t, int64(2), val)
		})
	})

	t.Run("hyperloglog with error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Pfadd("key1", "val1")
			assert.Error(t, err)
		})
	})
}

func TestRedis_List(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			_, err := newRedis(client.Addr, badType()).Lpush("key", "value1", "value2")
			assert.NotNil(t, err)
			val, err := client.Lpush("key", "value1", "value2")
			assert.Nil(t, err)
			assert.Equal(t, 2, val)
			_, err = newRedis(client.Addr, badType()).Rpush("key", "value3", "value4")
			assert.NotNil(t, err)
			val, err = client.Rpush("key", "value3", "value4")
			assert.Nil(t, err)
			assert.Equal(t, 4, val)
			_, err = newRedis(client.Addr, badType()).Llen("key")
			assert.NotNil(t, err)
			val, err = client.Llen("key")
			assert.Nil(t, err)
			assert.Equal(t, 4, val)
			_, err = newRedis(client.Addr, badType()).Lindex("key", 1)
			assert.NotNil(t, err)
			value, err := client.Lindex("key", 0)
			assert.Nil(t, err)
			assert.Equal(t, "value2", value)
			vals, err := client.Lrange("key", 0, 10)
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"value2", "value1", "value3", "value4"}, vals)
			_, err = newRedis(client.Addr, badType()).Lpop("key")
			assert.NotNil(t, err)
			v, err := client.Lpop("key")
			assert.Nil(t, err)
			assert.Equal(t, "value2", v)
			val, err = client.Lpush("key", "value1", "value2")
			assert.Nil(t, err)
			assert.Equal(t, 5, val)
			_, err = newRedis(client.Addr, badType()).Rpop("key")
			assert.NotNil(t, err)
			v, err = client.Rpop("key")
			assert.Nil(t, err)
			assert.Equal(t, "value4", v)
			val, err = client.Rpush("key", "value4", "value3", "value3")
			assert.Nil(t, err)
			assert.Equal(t, 7, val)
			_, err = newRedis(client.Addr, badType()).Lrem("key", 2, "value1")
			assert.NotNil(t, err)
			n, err := client.Lrem("key", 2, "value1")
			assert.Nil(t, err)
			assert.Equal(t, 2, n)
			_, err = newRedis(client.Addr, badType()).Lrange("key", 0, 10)
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
			err = newRedis(client.Addr, badType()).Ltrim("key", 0, 1)
			assert.Error(t, err)
			err = client.Ltrim("key", 0, 1)
			assert.Nil(t, err)
			vals, err = client.Lrange("key", 0, 10)
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"value2", "value3"}, vals)
			vals, err = client.LpopCount("key", 2)
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"value2", "value3"}, vals)
			_, err = client.Lpush("key", "value1", "value2")
			assert.Nil(t, err)
			vals, err = client.RpopCount("key", 4)
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"value1", "value2"}, vals)
		})
	})

	t.Run("list error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Llen("key")
			assert.Error(t, err)

			_, err = client.Lpush("key", "value1", "value2")
			assert.Error(t, err)

			_, err = client.Lrem("key", 2, "value1")
			assert.Error(t, err)

			_, err = client.Rpush("key", "value3", "value4")
			assert.Error(t, err)

			_, err = client.LpopCount("key", 2)
			assert.Error(t, err)

			_, err = client.RpopCount("key", 2)
			assert.Error(t, err)
		})
	})
	t.Run("list redis type error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			client.Type = "nil"
			_, err := client.Llen("key")
			assert.Error(t, err)

			_, err = client.Lpush("key", "value1", "value2")
			assert.Error(t, err)

			_, err = client.Lrem("key", 2, "value1")
			assert.Error(t, err)

			_, err = client.Rpush("key", "value3", "value4")
			assert.Error(t, err)

			_, err = client.LpopCount("key", 2)
			assert.Error(t, err)

			_, err = client.RpopCount("key", 2)
			assert.Error(t, err)
		})
	})
}

func TestRedis_Mset(t *testing.T) {
	t.Run("mset", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			// Attempt to Mget with a bad client type, expecting an error.
			_, err := New(client.Addr, badType()).Mset("key1", "value1")
			assert.NotNil(t, err)

			// Set multiple key-value pairs using Mset and expect no error.
			_, err = client.Mset("key1", "value1", "key2", "value2")
			assert.Nil(t, err)

			// Retrieve the values for the keys set above using Mget and expect no error.
			vals, err := client.Mget("key1", "key2")
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"value1", "value2"}, vals)
		})
	})

	// Test case for Mset operation with an incorrect number of arguments, expecting an error.
	t.Run("mset error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Mset("key1", "value1", "key2")
			assert.Error(t, err)
		})
	})
}

func TestRedis_Mget(t *testing.T) {
	t.Run("mget", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			err := client.Set("key1", "value1")
			assert.Nil(t, err)
			err = client.Set("key2", "value2")
			assert.Nil(t, err)
			_, err = newRedis(client.Addr, badType()).Mget("key1", "key0", "key2", "key3")
			assert.NotNil(t, err)
			vals, err := client.Mget("key1", "key0", "key2", "key3")
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"value1", "", "value2", ""}, vals)
		})
	})

	t.Run("mget error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Mget("key1", "key0")
			assert.Error(t, err)
		})
	})
}

func TestRedis_SetBit(t *testing.T) {
	t.Run("setbit", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			_, err := newRedis(client.Addr, badType()).SetBit("key", 1, 1)
			assert.Error(t, err)
			val, err := client.SetBit("key", 1, 1)
			if assert.NoError(t, err) {
				assert.Equal(t, 0, val)
			}
		})
	})

	t.Run("setbit error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.SetBit("key", 1, 1)
			assert.Error(t, err)
		})
	})
}

func TestRedis_GetBit(t *testing.T) {
	t.Run("getbit", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			val, err := client.SetBit("key", 2, 1)
			assert.Nil(t, err)
			assert.Equal(t, 0, val)
			_, err = newRedis(client.Addr, badType()).GetBit("key", 2)
			assert.NotNil(t, err)
			v, err := client.GetBit("key", 2)
			assert.Nil(t, err)
			assert.Equal(t, 1, v)
		})
	})

	t.Run("getbit error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.GetBit("key", 2)
			assert.Error(t, err)
		})
	})
}

func TestRedis_BitCount(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		for i := 0; i < 11; i++ {
			val, err := client.SetBit("key", int64(i), 1)
			assert.Nil(t, err)
			assert.Equal(t, 0, val)
		}

		_, err := newRedis(client.Addr, badType()).BitCount("key", 0, -1)
		assert.NotNil(t, err)
		val, err := client.BitCount("key", 0, -1)
		assert.Nil(t, err)
		assert.Equal(t, int64(11), val)

		val, err = client.BitCount("key", 0, 0)
		assert.Nil(t, err)
		assert.Equal(t, int64(8), val)

		val, err = client.BitCount("key", 1, 1)
		assert.Nil(t, err)
		assert.Equal(t, int64(3), val)

		val, err = client.BitCount("key", 0, 1)
		assert.Nil(t, err)
		assert.Equal(t, int64(11), val)

		val, err = client.BitCount("key", 2, 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(0), val)
	})
}

func TestRedis_BitOpAnd(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.Set("key1", "0")
		assert.Nil(t, err)
		err = client.Set("key2", "1")
		assert.Nil(t, err)
		_, err = newRedis(client.Addr, badType()).BitOpAnd("destKey", "key1", "key2")
		assert.NotNil(t, err)
		val, err := client.BitOpAnd("destKey", "key1", "key2")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
		valStr, err := client.Get("destKey")
		assert.Nil(t, err)
		// destKey  binary 110000   ascii 0
		assert.Equal(t, "0", valStr)
	})
}

func TestRedis_BitOpNot(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.Set("key1", "\u0000")
		assert.Nil(t, err)
		_, err = newRedis(client.Addr, badType()).BitOpNot("destKey", "key1")
		assert.NotNil(t, err)
		val, err := client.BitOpNot("destKey", "key1")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
		valStr, err := client.Get("destKey")
		assert.Nil(t, err)
		assert.Equal(t, "\xff", valStr)
	})
}

func TestRedis_BitOpOr(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.Set("key1", "1")
		assert.Nil(t, err)
		err = client.Set("key2", "0")
		assert.Nil(t, err)
		_, err = newRedis(client.Addr, badType()).BitOpOr("destKey", "key1", "key2")
		assert.NotNil(t, err)
		val, err := client.BitOpOr("destKey", "key1", "key2")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
		valStr, err := client.Get("destKey")
		assert.Nil(t, err)
		assert.Equal(t, "1", valStr)
	})
}

func TestRedis_BitOpXor(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.Set("key1", "\xff")
		assert.Nil(t, err)
		err = client.Set("key2", "\x0f")
		assert.Nil(t, err)
		_, err = newRedis(client.Addr, badType()).BitOpXor("destKey", "key1", "key2")
		assert.NotNil(t, err)
		val, err := client.BitOpXor("destKey", "key1", "key2")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), val)
		valStr, err := client.Get("destKey")
		assert.Nil(t, err)
		assert.Equal(t, "\xf0", valStr)
	})
}

func TestRedis_BitPos(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		// 11111111 11110000 00000000
		err := client.Set("key", "\xff\xf0\x00")
		assert.Nil(t, err)

		_, err = newRedis(client.Addr, badType()).BitPos("key", 0, 0, -1)
		assert.NotNil(t, err)
		val, err := client.BitPos("key", 0, 0, 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(12), val)

		val, err = client.BitPos("key", 1, 0, 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(0), val)

		val, err = client.BitPos("key", 0, 1, 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(12), val)

		val, err = client.BitPos("key", 1, 1, 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(8), val)

		val, err = client.BitPos("key", 1, 2, 2)
		assert.Nil(t, err)
		assert.Equal(t, int64(-1), val)
	})
}

func TestRedis_Persist(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Persist("key")
		assert.NotNil(t, err)
		ok, err := client.Persist("key")
		assert.Nil(t, err)
		assert.False(t, ok)
		err = client.Set("key", "value")
		assert.Nil(t, err)
		ok, err = client.Persist("key")
		assert.Nil(t, err)
		assert.False(t, ok)
		err = newRedis(client.Addr, badType()).Expire("key", 5)
		assert.NotNil(t, err)
		err = client.Expire("key", 5)
		assert.Nil(t, err)
		ok, err = client.Persist("key")
		assert.Nil(t, err)
		assert.True(t, ok)
		err = newRedis(client.Addr, badType()).Expireat("key", time.Now().Unix()+5)
		assert.NotNil(t, err)
		err = client.Expireat("key", time.Now().Unix()+5)
		assert.Nil(t, err)
		ok, err = client.Persist("key")
		assert.Nil(t, err)
		assert.True(t, ok)
	})
}

func TestRedis_Ping(t *testing.T) {
	t.Run("ping", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			ok := client.Ping()
			assert.True(t, ok)
		})
	})
}

func TestRedis_Scan(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := client.Set("key1", "value1")
		assert.Nil(t, err)
		err = client.Set("key2", "value2")
		assert.Nil(t, err)
		_, _, err = newRedis(client.Addr, badType()).Scan(0, "*", 100)
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
			list = append(list, stringx.Randn(i))
		}
		lens, err := client.Sadd(key, list)
		assert.Nil(t, err)
		assert.Equal(t, lens, 1550)

		var cursor uint64 = 0
		sum := 0
		for {
			_, _, err := newRedis(client.Addr, badType()).Sscan(key, cursor, "", 100)
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
		_, err = newRedis(client.Addr, badType()).Del(key)
		assert.NotNil(t, err)
		_, err = client.Del(key)
		assert.Nil(t, err)
	})
}

func TestRedis_Set(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			_, err := newRedis(client.Addr, badType()).Sadd("key", 1, 2, 3, 4)
			assert.NotNil(t, err)
			num, err := client.Sadd("key", 1, 2, 3, 4)
			assert.Nil(t, err)
			assert.Equal(t, 4, num)
			_, err = newRedis(client.Addr, badType()).Scard("key")
			assert.NotNil(t, err)
			val, err := client.Scard("key")
			assert.Nil(t, err)
			assert.Equal(t, int64(4), val)
			_, err = newRedis(client.Addr, badType()).Sismember("key", 2)
			assert.NotNil(t, err)
			ok, err := client.Sismember("key", 2)
			assert.Nil(t, err)
			assert.True(t, ok)
			_, err = newRedis(client.Addr, badType()).Srem("key", 3, 4)
			assert.NotNil(t, err)
			num, err = client.Srem("key", 3, 4)
			assert.Nil(t, err)
			assert.Equal(t, 2, num)
			_, err = newRedis(client.Addr, badType()).Smembers("key")
			assert.NotNil(t, err)
			vals, err := client.Smembers("key")
			assert.Nil(t, err)
			assert.ElementsMatch(t, []string{"1", "2"}, vals)
			_, err = newRedis(client.Addr, badType()).Srandmember("key", 1)
			assert.NotNil(t, err)
			members, err := client.Srandmember("key", 1)
			assert.Nil(t, err)
			assert.Len(t, members, 1)
			assert.Contains(t, []string{"1", "2"}, members[0])
			_, err = newRedis(client.Addr, badType()).Spop("key")
			assert.NotNil(t, err)
			member, err := client.Spop("key")
			assert.Nil(t, err)
			assert.Contains(t, []string{"1", "2"}, member)
			_, err = newRedis(client.Addr, badType()).Smembers("key")
			assert.NotNil(t, err)
			vals, err = client.Smembers("key")
			assert.Nil(t, err)
			assert.NotContains(t, vals, member)
			_, err = newRedis(client.Addr, badType()).Sadd("key1", 1, 2, 3, 4)
			assert.NotNil(t, err)
			num, err = client.Sadd("key1", 1, 2, 3, 4)
			assert.Nil(t, err)
			assert.Equal(t, 4, num)
			num, err = client.Sadd("key2", 2, 3, 4, 5)
			assert.Nil(t, err)
			assert.Equal(t, 4, num)
			_, err = newRedis(client.Addr, badType()).Sunion("key1", "key2")
			assert.NotNil(t, err)
			vals, err = client.Sunion("key1", "key2")
			assert.Nil(t, err)
			assert.ElementsMatch(t, []string{"1", "2", "3", "4", "5"}, vals)
			_, err = newRedis(client.Addr, badType()).Sunionstore("key3", "key1", "key2")
			assert.NotNil(t, err)
			num, err = client.Sunionstore("key3", "key1", "key2")
			assert.Nil(t, err)
			assert.Equal(t, 5, num)
			_, err = newRedis(client.Addr, badType()).Sdiff("key1", "key2")
			assert.NotNil(t, err)
			vals, err = client.Sdiff("key1", "key2")
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"1"}, vals)
			_, err = newRedis(client.Addr, badType()).Sdiffstore("key4", "key1", "key2")
			assert.NotNil(t, err)
			num, err = client.Sdiffstore("key4", "key1", "key2")
			assert.Nil(t, err)
			assert.Equal(t, 1, num)
			_, err = newRedis(client.Addr, badType()).Sinter("key1", "key2")
			assert.NotNil(t, err)
			vals, err = client.Sinter("key1", "key2")
			assert.Nil(t, err)
			assert.ElementsMatch(t, []string{"2", "3", "4"}, vals)
			_, err = newRedis(client.Addr, badType()).Sinterstore("key4", "key1", "key2")
			assert.NotNil(t, err)
			num, err = client.Sinterstore("key4", "key1", "key2")
			assert.Nil(t, err)
			assert.Equal(t, 3, num)
		})
	})

	t.Run("set with error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Sadd("key", 1, 2, 3, 4)
			assert.Error(t, err)

			_, err = client.Srem("key", 3, 4)
			assert.Error(t, err)

			_, err = client.Sunionstore("key3", "key1", "key2")
			assert.Error(t, err)

			_, err = client.Sdiffstore("key4", "key1", "key2")
			assert.Error(t, err)

			_, err = client.Sinterstore("key4", "key1", "key2")
			assert.Error(t, err)
		})
	})
}

func TestRedis_GetSet(t *testing.T) {
	t.Run("set_get", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			_, err := newRedis(client.Addr, badType()).GetSet("hello", "world")
			assert.NotNil(t, err)
			val, err := client.GetSet("hello", "world")
			assert.Nil(t, err)
			assert.Equal(t, "", val)
			val, err = client.Get("hello")
			assert.Nil(t, err)
			assert.Equal(t, "world", val)
			val, err = client.GetSet("hello", "newworld")
			assert.Nil(t, err)
			assert.Equal(t, "world", val)
			val, err = client.Get("hello")
			assert.Nil(t, err)
			assert.Equal(t, "newworld", val)
			ret, err := client.Del("hello")
			assert.Nil(t, err)
			assert.Equal(t, 1, ret)
		})
	})

	t.Run("set_get with notexists", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			_, err := client.Get("hello")
			assert.NoError(t, err)
		})
	})

	t.Run("set_get with error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Get("hello")
			assert.Error(t, err)
		})
	})
}

func TestRedis_SetGetDel(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := newRedis(client.Addr, badType()).Set("hello", "world")
		assert.NotNil(t, err)
		err = client.Set("hello", "world")
		assert.Nil(t, err)
		_, err = newRedis(client.Addr, badType()).Get("hello")
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
		err := newRedis(client.Addr, badType()).Setex("hello", "world", 5)
		assert.NotNil(t, err)
		err = client.Setex("hello", "world", 5)
		assert.Nil(t, err)
		_, err = newRedis(client.Addr, badType()).Setnx("hello", "newworld")
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
		_, err = newRedis(client.Addr, badType()).SetnxEx("newhello", "newworld", 5)
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
	t.Run("hash", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			err := client.Hset("key", "field", "value")
			assert.Nil(t, err)
			_, err = newRedis(client.Addr, badType()).Hget("key", "field")
			assert.NotNil(t, err)
			val, err := client.Hget("key", "field")
			assert.Nil(t, err)
			assert.Equal(t, "value", val)
			_, err = newRedis(client.Addr, badType()).Hexists("key", "field")
			assert.NotNil(t, err)
			ok, err := client.Hexists("key", "field")
			assert.Nil(t, err)
			assert.True(t, ok)
			_, err = newRedis(client.Addr, badType()).Hdel("key", "field")
			assert.NotNil(t, err)
			ret, err := client.Hdel("key", "field")
			assert.Nil(t, err)
			assert.True(t, ret)
			ok, err = client.Hexists("key", "field")
			assert.Nil(t, err)
			assert.False(t, ok)
		})
	})

	t.Run("hash error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Hdel("key", "field")
			assert.Error(t, err)

			_, err = client.Hincrby("key", "field", 1)
			assert.Error(t, err)

			_, err = client.HincrbyFloat("key", "field", 1)
			assert.Error(t, err)

			_, err = client.Hlen("key")
			assert.Error(t, err)

			_, err = client.Hmget("key", "field")
			assert.Error(t, err)
		})
	})
}

func TestRedis_SortedSet(t *testing.T) {
	t.Run("sorted set", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			ok, err := client.ZaddFloat("key", 1, "value1")
			assert.Nil(t, err)
			assert.True(t, ok)
			ok, err = client.Zadd("key", 2, "value1")
			assert.Nil(t, err)
			assert.False(t, ok)
			val, err := client.Zscore("key", "value1")
			assert.Nil(t, err)
			assert.Equal(t, int64(2), val)
			_, err = newRedis(client.Addr, badType()).Zincrby("key", 3, "value1")
			assert.NotNil(t, err)
			val, err = client.Zincrby("key", 3, "value1")
			assert.Nil(t, err)
			assert.Equal(t, int64(5), val)
			_, err = newRedis(client.Addr, badType()).Zscore("key", "value1")
			assert.NotNil(t, err)
			val, err = client.Zscore("key", "value1")
			assert.Nil(t, err)
			assert.Equal(t, int64(5), val)
			_, err = newRedis(client.Addr, badType()).Zadds("key")
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
			_, err = newRedis(client.Addr, badType()).ZRevRangeWithScores("key", 1, 3)
			assert.NotNil(t, err)
			_, err = client.ZRevRangeWithScores("key", 1, 3)
			assert.Nil(t, err)
			_, err = client.ZRevRangeWithScoresCtx(context.Background(), "key", 1, 3)
			assert.Nil(t, err)
			pairs, err := client.ZrevrangeWithScores("key", 1, 3)
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
			rank, err = client.Zrevrank("key", "value1")
			assert.Nil(t, err)
			assert.Equal(t, int64(2), rank)
			_, err = newRedis(client.Addr, badType()).Zrank("key", "value4")
			assert.NotNil(t, err)
			_, err = client.Zrank("key", "value4")
			assert.Equal(t, Nil, err)
			_, err = newRedis(client.Addr, badType()).Zrem("key", "value2", "value3")
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
			_, err = newRedis(client.Addr, badType()).Zremrangebyscore("key", 6, 7)
			assert.NotNil(t, err)
			num, err = client.Zremrangebyscore("key", 6, 7)
			assert.Nil(t, err)
			assert.Equal(t, 2, num)
			ok, err = client.Zadd("key", 6, "value2")
			assert.Nil(t, err)
			assert.True(t, ok)
			_, err = newRedis(client.Addr, badType()).Zadd("key", 7, "value3")
			assert.NotNil(t, err)
			ok, err = client.Zadd("key", 7, "value3")
			assert.Nil(t, err)
			assert.True(t, ok)
			_, err = newRedis(client.Addr, badType()).Zcount("key", 6, 7)
			assert.NotNil(t, err)
			num, err = client.Zcount("key", 6, 7)
			assert.Nil(t, err)
			assert.Equal(t, 2, num)
			_, err = newRedis(client.Addr, badType()).Zremrangebyrank("key", 1, 2)
			assert.NotNil(t, err)
			num, err = client.Zremrangebyrank("key", 1, 2)
			assert.Nil(t, err)
			assert.Equal(t, 2, num)
			_, err = newRedis(client.Addr, badType()).Zcard("key")
			assert.NotNil(t, err)
			card, err := client.Zcard("key")
			assert.Nil(t, err)
			assert.Equal(t, 2, card)
			_, err = newRedis(client.Addr, badType()).Zrange("key", 0, -1)
			assert.NotNil(t, err)
			vals, err := client.Zrange("key", 0, -1)
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"value1", "value4"}, vals)
			_, err = newRedis(client.Addr, badType()).Zrevrange("key", 0, -1)
			assert.NotNil(t, err)
			vals, err = client.Zrevrange("key", 0, -1)
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"value4", "value1"}, vals)
			_, err = newRedis(client.Addr, badType()).ZrangeWithScores("key", 0, -1)
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
			_, err = newRedis(client.Addr, badType()).ZrangebyscoreWithScores("key", 5, 8)
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
			_, err = newRedis(client.Addr, badType()).ZrangebyscoreWithScoresAndLimit(
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
			_, err = newRedis(client.Addr, badType()).ZrevrangebyscoreWithScores("key", 5, 8)
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
			_, err = newRedis(client.Addr, badType()).ZrevrangebyscoreWithScoresAndLimit(
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
			_, err = newRedis(client.Addr, badType()).Zrevrank("key", "value")
			assert.NotNil(t, err)
			_, _ = client.Zadd("second", 2, "aa")
			_, _ = client.Zadd("third", 3, "bbb")
			val, err = client.Zunionstore("union", &ZStore{
				Keys:      []string{"second", "third"},
				Weights:   []float64{1, 2},
				Aggregate: "SUM",
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(2), val)
			_, err = newRedis(client.Addr, badType()).Zunionstore("union", &ZStore{})
			assert.NotNil(t, err)
			vals, err = client.Zrange("union", 0, 10000)
			assert.Nil(t, err)
			assert.EqualValues(t, []string{"aa", "bbb"}, vals)
			ival, err := client.Zcard("union")
			assert.Nil(t, err)
			assert.Equal(t, 2, ival)
		})
	})

	t.Run("sorted set error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZaddFloat("key", 1, "value")
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Zadds("key")
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Zcard("key")
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Zcount("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Zincrby("key", 1, "value")
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Zscore("key", "value")
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Zrem("key", "value")
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Zremrangebyscore("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Zremrangebyrank("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrangeWithScores("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrangeWithScoresByFloat("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZRevRangeWithScores("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZRevRangeWithScoresByFloat("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrangebyscoreWithScores("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrangebyscoreWithScoresByFloat("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrangebyscoreWithScoresAndLimit("key", 1, 2, 1, 1)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrangebyscoreWithScoresByFloatAndLimit("key", 1, 2, 1, 1)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrevrangebyscoreWithScores("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrevrangebyscoreWithScoresByFloat("key", 1, 2)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrevrangebyscoreWithScoresAndLimit("key", 1, 2, 1, 1)
			assert.Error(t, err)
		})

		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.ZrevrangebyscoreWithScoresByFloatAndLimit("key", 1, 2, 1, 1)
			assert.Error(t, err)
		})
	})
}

func TestRedis_SortedSetByFloat64(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		ok, err := client.ZaddFloat("key", 10.345, "value1")
		assert.Nil(t, err)
		assert.True(t, ok)
		val, err := client.ZscoreByFloat("key", "value1")
		assert.Nil(t, err)
		assert.Equal(t, 10.345, val)
		_, err = newRedis(client.Addr, badType()).ZscoreByFloat("key", "value1")
		assert.Error(t, err)
		_, _ = client.ZaddFloat("key", 10.346, "value2")
		_, err = newRedis(client.Addr, badType()).ZRevRangeWithScoresByFloat("key", 0, -1)
		assert.NotNil(t, err)
		_, err = client.ZRevRangeWithScoresByFloat("key", 0, -1)
		assert.Nil(t, err)
		_, err = client.ZRevRangeWithScoresByFloatCtx(context.Background(), "key", 0, -1)
		assert.Nil(t, err)
		pairs, err := client.ZrevrangeWithScoresByFloat("key", 0, -1)
		assert.Nil(t, err)
		assert.EqualValues(t, []FloatPair{
			{
				Key:   "value2",
				Score: 10.346,
			},
			{
				Key:   "value1",
				Score: 10.345,
			},
		}, pairs)

		_, err = newRedis(client.Addr, badType()).ZrangeWithScoresByFloat("key", 0, -1)
		assert.Error(t, err)

		pairs, err = client.ZrangeWithScoresByFloat("key", 0, -1)
		if assert.NoError(t, err) {
			assert.EqualValues(t, []FloatPair{
				{
					Key:   "value1",
					Score: 10.345,
				},
				{
					Key:   "value2",
					Score: 10.346,
				},
			}, pairs)
		}

		_, err = newRedis(client.Addr, badType()).ZrangebyscoreWithScoresByFloat("key", 0, 20)
		assert.NotNil(t, err)
		pairs, err = client.ZrangebyscoreWithScoresByFloat("key", 0, 20)
		assert.Nil(t, err)
		assert.EqualValues(t, []FloatPair{
			{
				Key:   "value1",
				Score: 10.345,
			},
			{
				Key:   "value2",
				Score: 10.346,
			},
		}, pairs)
		_, err = client.ZrangebyscoreWithScoresByFloatAndLimit("key", 10.1, 12.2, 1, 0)
		assert.NoError(t, err)
		_, err = newRedis(client.Addr, badType()).ZrangebyscoreWithScoresByFloatAndLimit(
			"key", 10.1, 12.2, 1, 1)
		assert.NotNil(t, err)
		pairs, err = client.ZrangebyscoreWithScoresByFloatAndLimit("key", 10.1, 12.2, 1, 1)
		assert.Nil(t, err)
		assert.EqualValues(t, []FloatPair{
			{
				Key:   "value2",
				Score: 10.346,
			},
		}, pairs)
		_, err = newRedis(client.Addr, badType()).ZrevrangebyscoreWithScoresByFloat("key", 10, 12)
		assert.NotNil(t, err)
		pairs, err = client.ZrevrangebyscoreWithScoresByFloat("key", 10, 12)
		assert.Nil(t, err)
		assert.EqualValues(t, []FloatPair{
			{
				Key:   "value2",
				Score: 10.346,
			},
			{
				Key:   "value1",
				Score: 10.345,
			},
		}, pairs)
		_, err = client.ZrevrangebyscoreWithScoresByFloatAndLimit("key", 10, 12, 1, 0)
		assert.NoError(t, err)
		_, err = newRedis(client.Addr, badType()).ZrevrangebyscoreWithScoresByFloatAndLimit(
			"key", 10, 12, 1, 1)
		assert.NotNil(t, err)
		pairs, err = client.ZrevrangebyscoreWithScoresByFloatAndLimit("key", 10, 12, 1, 1)
		assert.Nil(t, err)
		assert.EqualValues(t, []FloatPair{
			{
				Key:   "value1",
				Score: 10.345,
			},
		}, pairs)
	})
}

func TestRedis_Zaddnx(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Zaddnx("key", 1, "value1")
		assert.Error(t, err)

		ok, err := client.Zadd("key", 1, "value1")
		assert.NoError(t, err)
		assert.True(t, ok)
		ok, err = client.Zaddnx("key", 2, "value1")
		assert.NoError(t, err)
		assert.False(t, ok)

		ok, err = client.ZaddFloat("key", 1.1, "value2")
		assert.NoError(t, err)
		assert.True(t, ok)
		ok, err = client.ZaddnxFloat("key", 1.1, "value3")
		assert.NoError(t, err)
		assert.True(t, ok)

		assert.NoError(t, client.Set("newkey", "value"))
		ok, err = client.Zaddnx("newkey", 1, "value")
		assert.Error(t, err)
		assert.False(t, ok)
	})
}

func TestRedis_IncrbyFloat(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		incrVal, err := client.IncrbyFloat("key", 0.002)
		assert.Nil(t, err)
		assert.Equal(t, 0.002, incrVal)
		_, err = newRedis(client.Addr, badType()).IncrbyFloat("key", -0.001)
		assert.Error(t, err)
		incrVal2, err := client.IncrbyFloat("key", -0.001)
		assert.Nil(t, err)
		assert.Equal(t, 0.001, incrVal2)
		_, err = newRedis(client.Addr, badType()).HincrbyFloat("hkey", "i", 0.002)
		assert.Error(t, err)
		hincrVal, err := client.HincrbyFloat("hkey", "i", 0.002)
		assert.Nil(t, err)
		assert.Equal(t, 0.002, hincrVal)
		hincrVal2, err := client.HincrbyFloat("hkey", "i", -0.001)
		assert.Nil(t, err)
		assert.Equal(t, 0.001, hincrVal2)
	})
}

func TestRedis_Pipelined(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		assert.NotNil(t, newRedis(client.Addr, badType()).Pipelined(func(pipeliner Pipeliner) error {
			return nil
		}))
		err := client.Pipelined(
			func(pipe Pipeliner) error {
				pipe.Incr(context.Background(), "pipelined_counter")
				pipe.Expire(context.Background(), "pipelined_counter", time.Hour)
				pipe.ZAdd(context.Background(), "zadd", Z{Score: 12, Member: "zadd"})
				return nil
			},
		)
		assert.Nil(t, err)
		_, err = newRedis(client.Addr, badType()).Ttl("pipelined_counter")
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
		_, err := getRedis(newRedis(client.Addr, Cluster()))
		assert.Nil(t, err)
		assert.Equal(t, client.Addr, client.String())
		assert.NotNil(t, newRedis(client.Addr, badType()).Ping())
	})
}

func TestRedisScriptLoad(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		client.Ping()
		_, err := newRedis(client.Addr, badType()).ScriptLoad("foo")
		assert.NotNil(t, err)
		_, err = client.ScriptLoad("foo")
		assert.NotNil(t, err)
	})
}

func TestRedisEvalSha(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		client.Ping()
		scriptHash, err := client.ScriptLoad(`return redis.call("EXISTS", KEYS[1])`)
		assert.Nil(t, err)
		_, err = newRedis(client.Addr, badType()).EvalSha(scriptHash, []string{"key1"})
		assert.Error(t, err)
		result, err := client.EvalSha(scriptHash, []string{"key1"})
		assert.Nil(t, err)
		assert.Equal(t, int64(0), result)
	})
}

func TestRedis_Ttl(t *testing.T) {
	t.Run("TTL", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			if assert.NoError(t, client.Set("key", "value")) {
				_, err := client.Ttl("key")
				assert.NoError(t, err)
			}
		})
	})

	t.Run("TTL error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.Ttl("key")
			assert.Error(t, err)
		})
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

func TestRedisToFloatPairs(t *testing.T) {
	pairs := toFloatPairs([]red.Z{
		{
			Member: 1,
			Score:  1,
		},
		{
			Member: 2,
			Score:  2,
		},
	})
	assert.EqualValues(t, []FloatPair{
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

func TestRedis_Zscan(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		key := "key"
		for i := 0; i < 1550; i++ {
			ok, err := client.Zadd(key, int64(i), "value_"+strconv.Itoa(i))
			assert.Nil(t, err)
			assert.True(t, ok)
		}
		var cursor uint64 = 0
		sum := 0
		for {
			_, _, err := newRedis(client.Addr, badType()).Zscan(key, cursor, "value_*", 100)
			assert.NotNil(t, err)
			keys, next, err := client.Zscan(key, cursor, "value_*", 100)
			assert.Nil(t, err)
			sum += len(keys)
			if next == 0 {
				break
			}
			cursor = next
		}

		assert.Equal(t, sum, 3100)
		_, err := newRedis(client.Addr, badType()).Del(key)
		assert.NotNil(t, err)
		_, err = client.Del(key)
		assert.Nil(t, err)
	})
}

func TestRedisToStrings(t *testing.T) {
	vals := toStrings([]any{1, 2})
	assert.EqualValues(t, []string{"1", "2"}, vals)
}

func TestRedisBlpop(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		client.Ping()
		var node mockedNode
		_, err := client.Blpop(nil, "foo")
		assert.Error(t, err)
		_, err = client.Blpop(node, "foo")
		assert.NoError(t, err)
	})
}

func TestRedisBlpopEx(t *testing.T) {
	t.Run("blpopex", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			client.Ping()
			var node mockedNode
			_, _, err := client.BlpopEx(nil, "foo")
			assert.Error(t, err)
			_, _, err = client.BlpopEx(node, "foo")
			assert.NoError(t, err)
		})
	})

	t.Run("blpopex with expire", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			client.Ping()
			node, err := getRedis(client)
			assert.NoError(t, err)
			assert.NoError(t, client.Set("foo", "bar"))
			_, _, err = client.BlpopEx(node, "foo")
			assert.Error(t, err)
		})
	})

	t.Run("blpopex with bad return", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			client.Ping()
			node := &mockedNode{args: []string{"bar"}}
			_, _, err := client.BlpopEx(node, "foo")
			assert.Error(t, err)
		})
	})
}

func TestRedisBlpopWithTimeout(t *testing.T) {
	t.Run("blpop_withTimeout", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			client.Ping()
			var node mockedNode
			_, err := client.BlpopWithTimeout(nil, 10*time.Second, "foo")
			assert.Error(t, err)
			_, err = client.BlpopWithTimeout(node, 10*time.Second, "foo")
			assert.NoError(t, err)
		})
	})

	t.Run("blpop_withTimeout_error", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			client.Ping()
			node, err := getRedis(client)
			assert.NoError(t, err)
			assert.NoError(t, client.Set("foo", "bar"))
			_, err = client.BlpopWithTimeout(node, time.Millisecond, "foo")
			assert.Error(t, err)
		})
	})

	t.Run("blpop_with_bad_return", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			node := &mockedNode{args: []string{"foo"}}
			_, err := client.Blpop(node, "foo")
			assert.Error(t, err)
		})
	})
}

func TestRedisGeo(t *testing.T) {
	t.Run("geo", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			client.Ping()
			geoLocation := []*GeoLocation{{Longitude: 13.361389, Latitude: 38.115556, Name: "Palermo"},
				{Longitude: 15.087269, Latitude: 37.502669, Name: "Catania"}}
			v, err := client.GeoAdd("sicily", geoLocation...)
			assert.Nil(t, err)
			assert.Equal(t, int64(2), v)
			_, err = newRedis(client.Addr, badType()).GeoDist("sicily", "Palermo", "Catania", "m")
			assert.Error(t, err)
			v2, err := client.GeoDist("sicily", "Palermo", "Catania", "m")
			assert.Nil(t, err)
			assert.Equal(t, 166274, int(v2))
			// GeoHash not support
			_, err = newRedis(client.Addr, badType()).GeoPos("sicily", "Palermo", "Catania")
			assert.Error(t, err)
			v3, err := client.GeoPos("sicily", "Palermo", "Catania")
			assert.Nil(t, err)
			assert.Equal(t, int64(v3[0].Longitude), int64(13))
			assert.Equal(t, int64(v3[0].Latitude), int64(38))
			assert.Equal(t, int64(v3[1].Longitude), int64(15))
			assert.Equal(t, int64(v3[1].Latitude), int64(37))
			_, err = newRedis(client.Addr, badType()).GeoRadius("sicily", 15, 37,
				&red.GeoRadiusQuery{WithDist: true, Unit: "km", Radius: 200})
			assert.Error(t, err)
			v4, err := client.GeoRadius("sicily", 15, 37, &red.GeoRadiusQuery{
				WithDist: true,
				Unit:     "km", Radius: 200,
			})
			assert.Nil(t, err)
			assert.Equal(t, int64(v4[0].Dist), int64(190))
			assert.Equal(t, int64(v4[1].Dist), int64(56))
			geoLocation2 := []*GeoLocation{{Longitude: 13.583333, Latitude: 37.316667, Name: "Agrigento"}}
			_, err = newRedis(client.Addr, badType()).GeoAdd("sicily", geoLocation2...)
			assert.Error(t, err)
			v5, err := client.GeoAdd("sicily", geoLocation2...)
			assert.Nil(t, err)
			assert.Equal(t, int64(1), v5)
			_, err = newRedis(client.Addr, badType()).GeoRadiusByMember("sicily", "Agrigento",
				&red.GeoRadiusQuery{Unit: "km", Radius: 100})
			assert.Error(t, err)
			v6, err := client.GeoRadiusByMember("sicily", "Agrigento",
				&red.GeoRadiusQuery{Unit: "km", Radius: 100})
			assert.Nil(t, err)
			assert.Equal(t, v6[0].Name, "Agrigento")
			assert.Equal(t, v6[1].Name, "Palermo")
		})
	})

	t.Run("geo error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			_, err := client.GeoAdd("sicily", &GeoLocation{
				Longitude: 13.3,
				Latitude:  38.1,
				Name:      "Palermo",
			})
			assert.Error(t, err)

			_, err = client.GeoDist("sicily", "Palermo", "Catania", "m")
			assert.Error(t, err)

			_, err = client.GeoRadius("sicily", 15, 37, &red.GeoRadiusQuery{
				WithDist: true,
			})
			assert.Error(t, err)

			_, err = client.GeoRadiusByMember("sicily", "Agrigento", &red.GeoRadiusQuery{
				Unit: "km",
			})
			assert.Error(t, err)

			_, err = client.GeoPos("sicily", "Palermo", "Catania")
			assert.Error(t, err)
		})
	})
}

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold.Load())
	SetSlowThreshold(time.Second)
	assert.Equal(t, time.Second, slowThreshold.Load())
}

func TestRedis_WithUserPass(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		err := newRedis(client.Addr, WithUser("any"), WithPass("any")).Ping()
		assert.NotNil(t, err)
	})
}

func TestRedis_checkConnection(t *testing.T) {
	t.Run("checkConnection", func(t *testing.T) {
		runOnRedis(t, func(client *Redis) {
			client.Ping()
			assert.NoError(t, client.checkConnection(time.Millisecond))
		})
	})

	t.Run("checkConnection error", func(t *testing.T) {
		runOnRedisWithError(t, func(client *Redis) {
			assert.Error(t, newRedis(client.Addr, badType()).checkConnection(time.Millisecond))
			assert.Error(t, client.checkConnection(time.Millisecond))
		})
	})
}

func runOnRedis(t *testing.T, fn func(client *Redis)) {
	logx.Disable()

	s := miniredis.RunT(t)
	fn(MustNewRedis(RedisConf{
		Host: s.Addr(),
		Type: NodeType,
	}))
}

func runOnRedisWithError(t *testing.T, fn func(client *Redis)) {
	logx.Disable()

	s := miniredis.RunT(t)
	s.SetError("mock error")
	fn(newRedis(s.Addr()))
}

func runOnRedisTLS(t *testing.T, fn func(client *Redis)) {
	logx.Disable()

	s, err := miniredis.RunTLS(&tls.Config{
		Certificates:       make([]tls.Certificate, 1),
		InsecureSkipVerify: true,
	})
	assert.Nil(t, err)
	defer func() {
		client, err := clientManager.GetResource(s.Addr(), func() (io.Closer, error) {
			return nil, errors.New("should already exist")
		})
		if err != nil {
			t.Error(err)
		}
		if client != nil {
			_ = client.Close()
		}
	}()
	fn(newRedis(s.Addr(), WithTLS()))
}

func badType() Option {
	return func(r *Redis) {
		r.Type = "bad"
	}
}

type mockedNode struct {
	RedisNode
	args []string
}

func (n mockedNode) BLPop(_ context.Context, _ time.Duration, _ ...string) *red.StringSliceCmd {
	cmd := red.NewStringSliceCmd(context.Background())
	if len(n.args) == 0 {
		cmd.SetVal([]string{"foo", "bar"})
	} else {
		cmd.SetVal(n.args)
	}

	return cmd
}

func TestRedisPublish(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Publish("Test", "message")
		assert.NotNil(t, err)
		_, err = client.Publish("Test", "message")
		assert.Nil(t, err)
	})
}

func TestRedisRPopLPush(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).RPopLPush("Source", "Destination")
		assert.NotNil(t, err)
		_, err = client.Rpush("Source", "Destination")
		assert.Nil(t, err)
		_, err = client.RPopLPush("Source", "Destination")
		assert.Nil(t, err)
	})
}

func TestRedisUnlink(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		_, err := newRedis(client.Addr, badType()).Unlink("Key1", "Key2")
		assert.NotNil(t, err)
		err = client.Set("Key1", "Key2")
		assert.Nil(t, err)
		get, err := client.Get("Key1")
		assert.Nil(t, err)
		assert.Equal(t, "Key2", get)
		res, err := client.Unlink("Key1")
		assert.Nil(t, err)
		assert.Equal(t, int64(1), res)
	})
}

func TestRedisTxPipeline(t *testing.T) {
	runOnRedis(t, func(client *Redis) {
		ctx := context.Background()
		_, err := newRedis(client.Addr, badType()).TxPipeline()
		assert.NotNil(t, err)
		pipe, err := client.TxPipeline()
		assert.Nil(t, err)
		key := "key"
		hashKey := "field"
		hashValue := "value"

		// setting value
		pipe.HSet(ctx, key, hashKey, hashValue)

		existsCmd := pipe.Exists(ctx, key)
		getCmd := pipe.HGet(ctx, key, hashKey)

		// execution
		_, err = pipe.Exec(ctx)
		assert.Nil(t, err)

		// verification results
		exists, err := existsCmd.Result()
		assert.Nil(t, err)
		assert.Equal(t, int64(1), exists)

		value, err := getCmd.Result()
		assert.Nil(t, err)
		assert.Equal(t, hashValue, value)
	})
}
