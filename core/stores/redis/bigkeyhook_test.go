package redis

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"

	"github.com/zeromicro/go-zero/core/logx"
)

func TestBigKeyHook_AfterProcess_Get(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	ctx := context.Background()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	err := r.Set("foo", "123456")
	assert.NoError(t, err)

	_, _ = r.Get("foo")
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	err = r.Set("foo2", "1234")
	assert.NoError(t, err)

	_, _ = r.Get("foo2")
	assert.NotContains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.GetCtx(ctx, "foo")
	assert.Contains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_AfterProcess_Set(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	_ = r.Set("foo", "123456")
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_ = r.Set("foo2", "1234")
	assert.NotContains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.Setnx("foo3", "123456")
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_ = r.Setex("foo4", "123456", 10)
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.SetnxExCtx(context.Background(), "foo5", "123456", 10)
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.SetBit("foo6", 1, 1)
	assert.NotContains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_AfterProcess_GetSet(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	_, _ = r.GetSet("foo", "123456")
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.GetSet("foo2", "1234")

	_, _ = r.GetSet("foo2", "123456")
	assert.NotContains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_AfterProcess_Hgetall(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	_, _ = r.Hgetall("foo")
	assert.NotContains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_ = r.Hset("foo", "bar", "123456")
	_ = r.Hset("foo", "bar2", "123456")
	_, _ = r.Hgetall("foo")
	assert.Contains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_AfterProcess_Hget(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	_, _ = r.Hget("foo", "bar")
	assert.NotContains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_ = r.Hset("foo", "bar", "123456")
	_, _ = r.Hget("foo", "bar")
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_ = r.Hset("foo", "bar2", "123456")
	_, _ = r.Hmget("foo", "bar1", "bar2")
	assert.Contains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_AfterProcess_Hset(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	_ = r.Hset("foo", "bar", "123456")
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.Hsetnx("foo2", "bar", "123456")
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_ = r.Hmset("foo3", map[string]string{"bar": "123456", "bar2": "123456"})
	assert.Contains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_AfterProcess_Sadd(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	_, _ = r.Sadd("foo", "123456", "123456")
	assert.Contains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_AfterProcess_Smembers(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	_, _ = r.Sadd("foo", "123456", "123456")
	_, _ = r.Smembers("foo")
	assert.Contains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_AfterProcess_Zadd(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	_, _ = r.Zadd("foo", 1, "123456")
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.ZaddFloat("foo", 1, "123456")
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	r.Zadds("foo2", Pair{"111111", 1}, Pair{"2222222", 2})
	assert.Contains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_AfterProcess_Zrange(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:     true,
			LimitSize:  5,
			LimitCount: 1,
			BufferLen:  0,
		},
	})

	_, _ = r.Zadd("foo", 1, "123456")

	buf.Reset()
	_, _ = r.Zrange("foo", 0, -1)
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.ZrangebyscoreWithScoresCtx(context.Background(), "foo", 0, 100)
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.ZrangeWithScores("foo", 0, 100)
	assert.Contains(t, buf.String(), "BigKey limit")

	buf.Reset()
	_, _ = r.ZrangebyscoreWithScoresAndLimit("foo", 0, 100, 0, 10)
	assert.Contains(t, buf.String(), "BigKey limit")
}

func TestBigKeyHook_stat(t *testing.T) {
	var buf bytes.Buffer
	logx.Reset()
	logx.SetLevel(logx.InfoLevel)
	logx.SetWriter(logx.NewWriter(&buf))
	defer logx.Reset()

	r := MustNewRedis(RedisConf{
		Host: miniredis.RunT(t).Addr(),
		Type: "node",
		VerifyBigKey: BigKeyHookConfig{
			Enable:       true,
			LimitSize:    5,
			LimitCount:   1,
			BufferLen:    100,
			StatInterval: time.Millisecond * 100,
		},
	})

	err := r.Set("foo", "123456")
	assert.NoError(t, err)

	for i := 0; i < 99; i++ {
		_, _ = r.Get("foo")
	}

	time.Sleep(time.Second)

	assert.Contains(t, buf.String(), "[REDIS] BigKey limit, key: foo, size: 6, count: 100")

}
