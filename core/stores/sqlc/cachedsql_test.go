package sqlc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func init() {
	logx.Disable()
	stat.SetReporter(nil)
}

func TestCachedConn_GetCache(t *testing.T) {
	resetStats()
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))
	var value string
	err = c.GetCache("any", &value)
	assert.Equal(t, ErrNotFound, err)
	r.Set("any", `"value"`)
	err = c.GetCache("any", &value)
	assert.Nil(t, err)
	assert.Equal(t, "value", value)
}

func TestStat(t *testing.T) {
	resetStats()
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))

	for i := 0; i < 10; i++ {
		var str string
		err = c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v interface{}) error {
			*v.(*string) = "zero"
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	}

	assert.Equal(t, uint64(10), atomic.LoadUint64(&stats.Total))
	assert.Equal(t, uint64(9), atomic.LoadUint64(&stats.Hit))
}

func TestCachedConn_QueryRowIndex_NoCache(t *testing.T) {
	resetStats()
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	c := NewConn(dummySqlConn{}, cache.CacheConf{
		{
			RedisConf: redis.RedisConf{
				Host: r.Addr,
				Type: redis.NodeType,
			},
			Weight: 100,
		},
	}, cache.WithExpiry(time.Second*10))

	var str string
	err = c.QueryRowIndex(&str, "index", func(s interface{}) string {
		return fmt.Sprintf("%s/1234", s)
	}, func(conn sqlx.SqlConn, v interface{}) (interface{}, error) {
		*v.(*string) = "zero"
		return "primary", errors.New("foo")
	}, func(conn sqlx.SqlConn, v, pri interface{}) error {
		assert.Equal(t, "primary", pri)
		*v.(*string) = "xin"
		return nil
	})
	assert.NotNil(t, err)

	err = c.QueryRowIndex(&str, "index", func(s interface{}) string {
		return fmt.Sprintf("%s/1234", s)
	}, func(conn sqlx.SqlConn, v interface{}) (interface{}, error) {
		*v.(*string) = "zero"
		return "primary", nil
	}, func(conn sqlx.SqlConn, v, pri interface{}) error {
		assert.Equal(t, "primary", pri)
		*v.(*string) = "xin"
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "zero", str)
	val, err := r.Get("index")
	assert.Nil(t, err)
	assert.Equal(t, `"primary"`, val)
	val, err = r.Get("primary/1234")
	assert.Nil(t, err)
	assert.Equal(t, `"zero"`, val)
}

func TestCachedConn_QueryRowIndex_HasCache(t *testing.T) {
	resetStats()
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10),
		cache.WithNotFoundExpiry(time.Second))

	var str string
	r.Set("index", `"primary"`)
	err = c.QueryRowIndex(&str, "index", func(s interface{}) string {
		return fmt.Sprintf("%s/1234", s)
	}, func(conn sqlx.SqlConn, v interface{}) (interface{}, error) {
		assert.Fail(t, "should not go here")
		return "primary", nil
	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
		*v.(*string) = "xin"
		assert.Equal(t, "primary", primary)
		return nil
	})
	assert.Nil(t, err)
	assert.Equal(t, "xin", str)
	val, err := r.Get("index")
	assert.Nil(t, err)
	assert.Equal(t, `"primary"`, val)
	val, err = r.Get("primary/1234")
	assert.Nil(t, err)
	assert.Equal(t, `"xin"`, val)
}

func TestCachedConn_QueryRowIndex_HasCache_IntPrimary(t *testing.T) {
	const (
		primaryInt8   int8   = 100
		primaryInt16  int16  = 10000
		primaryInt32  int32  = 10000000
		primaryInt64  int64  = 10000000
		primaryUint8  uint8  = 100
		primaryUint16 uint16 = 10000
		primaryUint32 uint32 = 10000000
		primaryUint64 uint64 = 10000000
	)
	tests := []struct {
		name         string
		primary      interface{}
		primaryCache string
	}{
		{
			name:         "int8 primary",
			primary:      primaryInt8,
			primaryCache: fmt.Sprint(primaryInt8),
		},
		{
			name:         "int16 primary",
			primary:      primaryInt16,
			primaryCache: fmt.Sprint(primaryInt16),
		},
		{
			name:         "int32 primary",
			primary:      primaryInt32,
			primaryCache: fmt.Sprint(primaryInt32),
		},
		{
			name:         "int64 primary",
			primary:      primaryInt64,
			primaryCache: fmt.Sprint(primaryInt64),
		},
		{
			name:         "uint8 primary",
			primary:      primaryUint8,
			primaryCache: fmt.Sprint(primaryUint8),
		},
		{
			name:         "uint16 primary",
			primary:      primaryUint16,
			primaryCache: fmt.Sprint(primaryUint16),
		},
		{
			name:         "uint32 primary",
			primary:      primaryUint32,
			primaryCache: fmt.Sprint(primaryUint32),
		},
		{
			name:         "uint64 primary",
			primary:      primaryUint64,
			primaryCache: fmt.Sprint(primaryUint64),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resetStats()
			r, clean, err := redistest.CreateRedis()
			assert.Nil(t, err)
			defer clean()

			c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10),
				cache.WithNotFoundExpiry(time.Second))

			var str string
			r.Set("index", test.primaryCache)
			err = c.QueryRowIndex(&str, "index", func(s interface{}) string {
				return fmt.Sprintf("%v/1234", s)
			}, func(conn sqlx.SqlConn, v interface{}) (interface{}, error) {
				assert.Fail(t, "should not go here")
				return test.primary, nil
			}, func(conn sqlx.SqlConn, v, primary interface{}) error {
				*v.(*string) = "xin"
				assert.Equal(t, primary, primary)
				return nil
			})
			assert.Nil(t, err)
			assert.Equal(t, "xin", str)
			val, err := r.Get("index")
			assert.Nil(t, err)
			assert.Equal(t, test.primaryCache, val)
			val, err = r.Get(test.primaryCache + "/1234")
			assert.Nil(t, err)
			assert.Equal(t, `"xin"`, val)
		})
	}
}

func TestCachedConn_QueryRowIndex_HasWrongCache(t *testing.T) {
	caches := map[string]string{
		"index":        "primary",
		"primary/1234": "xin",
	}

	for k, v := range caches {
		t.Run(k+"/"+v, func(t *testing.T) {
			resetStats()
			r, clean, err := redistest.CreateRedis()
			assert.Nil(t, err)
			defer clean()

			c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10),
				cache.WithNotFoundExpiry(time.Second))

			var str string
			r.Set(k, v)
			err = c.QueryRowIndex(&str, "index", func(s interface{}) string {
				return fmt.Sprintf("%s/1234", s)
			}, func(conn sqlx.SqlConn, v interface{}) (interface{}, error) {
				*v.(*string) = "xin"
				return "primary", nil
			}, func(conn sqlx.SqlConn, v, primary interface{}) error {
				*v.(*string) = "xin"
				assert.Equal(t, "primary", primary)
				return nil
			})
			assert.Nil(t, err)
			assert.Equal(t, "xin", str)
			val, err := r.Get("index")
			assert.Nil(t, err)
			assert.Equal(t, `"primary"`, val)
			val, err = r.Get("primary/1234")
			assert.Nil(t, err)
			assert.Equal(t, `"xin"`, val)
		})
	}
}

func TestStatCacheFails(t *testing.T) {
	resetStats()
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stdout)

	r := redis.New("localhost:59999")
	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))

	for i := 0; i < 20; i++ {
		var str string
		err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v interface{}) error {
			return errors.New("db failed")
		})
		assert.NotNil(t, err)
	}

	assert.Equal(t, uint64(20), atomic.LoadUint64(&stats.Total))
	assert.Equal(t, uint64(0), atomic.LoadUint64(&stats.Hit))
	assert.Equal(t, uint64(20), atomic.LoadUint64(&stats.Miss))
	assert.Equal(t, uint64(0), atomic.LoadUint64(&stats.DbFails))
}

func TestStatDbFails(t *testing.T) {
	resetStats()
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))

	for i := 0; i < 20; i++ {
		var str string
		err = c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v interface{}) error {
			return errors.New("db failed")
		})
		assert.NotNil(t, err)
	}

	assert.Equal(t, uint64(20), atomic.LoadUint64(&stats.Total))
	assert.Equal(t, uint64(0), atomic.LoadUint64(&stats.Hit))
	assert.Equal(t, uint64(20), atomic.LoadUint64(&stats.DbFails))
}

func TestStatFromMemory(t *testing.T) {
	resetStats()
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))

	var all sync.WaitGroup
	var wait sync.WaitGroup
	all.Add(10)
	wait.Add(4)
	go func() {
		var str string
		err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v interface{}) error {
			*v.(*string) = "zero"
			return nil
		})
		if err != nil {
			t.Error(err)
		}
		wait.Wait()
		runtime.Gosched()
		all.Done()
	}()

	for i := 0; i < 4; i++ {
		go func() {
			var str string
			wait.Done()
			err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v interface{}) error {
				*v.(*string) = "zero"
				return nil
			})
			if err != nil {
				t.Error(err)
			}
			all.Done()
		}()
	}
	for i := 0; i < 5; i++ {
		go func() {
			var str string
			err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v interface{}) error {
				*v.(*string) = "zero"
				return nil
			})
			if err != nil {
				t.Error(err)
			}
			all.Done()
		}()
	}
	all.Wait()

	assert.Equal(t, uint64(10), atomic.LoadUint64(&stats.Total))
	assert.Equal(t, uint64(9), atomic.LoadUint64(&stats.Hit))
}

func TestCachedConnQueryRow(t *testing.T) {
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	const (
		key   = "user"
		value = "any"
	)
	var conn trackedConn
	var user string
	var ran bool
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	err = c.QueryRow(&user, key, func(conn sqlx.SqlConn, v interface{}) error {
		ran = true
		user = value
		return nil
	})
	assert.Nil(t, err)
	actualValue, err := r.Get(key)
	assert.Nil(t, err)
	var actual string
	assert.Nil(t, json.Unmarshal([]byte(actualValue), &actual))
	assert.Equal(t, value, actual)
	assert.Equal(t, value, user)
	assert.True(t, ran)
}

func TestCachedConnQueryRowFromCache(t *testing.T) {
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	const (
		key   = "user"
		value = "any"
	)
	var conn trackedConn
	var user string
	var ran bool
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	assert.Nil(t, c.SetCache(key, value))
	err = c.QueryRow(&user, key, func(conn sqlx.SqlConn, v interface{}) error {
		ran = true
		user = value
		return nil
	})
	assert.Nil(t, err)
	actualValue, err := r.Get(key)
	assert.Nil(t, err)
	var actual string
	assert.Nil(t, json.Unmarshal([]byte(actualValue), &actual))
	assert.Equal(t, value, actual)
	assert.Equal(t, value, user)
	assert.False(t, ran)
}

func TestQueryRowNotFound(t *testing.T) {
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	const key = "user"
	var conn trackedConn
	var user string
	var ran int
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	for i := 0; i < 20; i++ {
		err = c.QueryRow(&user, key, func(conn sqlx.SqlConn, v interface{}) error {
			ran++
			return sql.ErrNoRows
		})
		assert.Exactly(t, sqlx.ErrNotFound, err)
	}
	assert.Equal(t, 1, ran)
}

func TestCachedConnExec(t *testing.T) {
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	var conn trackedConn
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	_, err = c.ExecNoCache("delete from user_table where id='kevin'")
	assert.Nil(t, err)
	assert.True(t, conn.execValue)
}

func TestCachedConnExecDropCache(t *testing.T) {
	r, err := miniredis.Run()
	assert.Nil(t, err)
	defer fx.DoWithTimeout(func() error {
		r.Close()
		return nil
	}, time.Second)

	const (
		key   = "user"
		value = "any"
	)
	var conn trackedConn
	c := NewNodeConn(&conn, redis.New(r.Addr()), cache.WithExpiry(time.Second*30))
	assert.Nil(t, c.SetCache(key, value))
	_, err = c.Exec(func(conn sqlx.SqlConn) (result sql.Result, e error) {
		return conn.Exec("delete from user_table where id='kevin'")
	}, key)
	assert.Nil(t, err)
	assert.True(t, conn.execValue)
	_, err = r.Get(key)
	assert.Exactly(t, miniredis.ErrKeyNotFound, err)
	_, err = c.Exec(func(conn sqlx.SqlConn) (result sql.Result, e error) {
		return nil, errors.New("foo")
	}, key)
	assert.NotNil(t, err)
}

func TestCachedConnExecDropCacheFailed(t *testing.T) {
	const key = "user"
	var conn trackedConn
	r := redis.New("anyredis:8888")
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	_, err := c.Exec(func(conn sqlx.SqlConn) (result sql.Result, e error) {
		return conn.Exec("delete from user_table where id='kevin'")
	}, key)
	// async background clean, retry logic
	assert.Nil(t, err)
}

func TestCachedConnQueryRows(t *testing.T) {
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	var conn trackedConn
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	var users []string
	err = c.QueryRowsNoCache(&users, "select user from user_table where id='kevin'")
	assert.Nil(t, err)
	assert.True(t, conn.queryRowsValue)
}

func TestCachedConnTransact(t *testing.T) {
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	var conn trackedConn
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	err = c.Transact(func(session sqlx.Session) error {
		return nil
	})
	assert.Nil(t, err)
	assert.True(t, conn.transactValue)
}

func TestQueryRowNoCache(t *testing.T) {
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	const (
		key   = "user"
		value = "any"
	)
	var user string
	var ran bool
	conn := dummySqlConn{queryRow: func(v interface{}, q string, args ...interface{}) error {
		user = value
		ran = true
		return nil
	}}
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	err = c.QueryRowNoCache(&user, key)
	assert.Nil(t, err)
	assert.Equal(t, value, user)
	assert.True(t, ran)
}

func TestNewConnWithCache(t *testing.T) {
	r, clean, err := redistest.CreateRedis()
	assert.Nil(t, err)
	defer clean()

	var conn trackedConn
	c := NewConnWithCache(&conn, cache.NewNode(r, singleFlights, stats, sql.ErrNoRows))
	_, err = c.ExecNoCache("delete from user_table where id='kevin'")
	assert.Nil(t, err)
	assert.True(t, conn.execValue)
}

func resetStats() {
	atomic.StoreUint64(&stats.Total, 0)
	atomic.StoreUint64(&stats.Hit, 0)
	atomic.StoreUint64(&stats.Miss, 0)
	atomic.StoreUint64(&stats.DbFails, 0)
}

type dummySqlConn struct {
	queryRow func(interface{}, string, ...interface{}) error
}

func (d dummySqlConn) ExecCtx(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (d dummySqlConn) PrepareCtx(ctx context.Context, query string) (sqlx.StmtSession, error) {
	return nil, nil
}

func (d dummySqlConn) QueryRowPartialCtx(ctx context.Context, v interface{}, query string, args ...interface{}) error {
	return nil
}

func (d dummySqlConn) QueryRowsCtx(ctx context.Context, v interface{}, query string, args ...interface{}) error {
	return nil
}

func (d dummySqlConn) QueryRowsPartialCtx(ctx context.Context, v interface{}, query string, args ...interface{}) error {
	return nil
}

func (d dummySqlConn) TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return nil
}

func (d dummySqlConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (d dummySqlConn) Prepare(query string) (sqlx.StmtSession, error) {
	return nil, nil
}

func (d dummySqlConn) QueryRow(v interface{}, query string, args ...interface{}) error {
	return d.QueryRowCtx(context.Background(), v, query, args...)
}

func (d dummySqlConn) QueryRowCtx(_ context.Context, v interface{}, query string, args ...interface{}) error {
	if d.queryRow != nil {
		return d.queryRow(v, query, args...)
	}
	return nil
}

func (d dummySqlConn) QueryRowPartial(v interface{}, query string, args ...interface{}) error {
	return nil
}

func (d dummySqlConn) QueryRows(v interface{}, query string, args ...interface{}) error {
	return nil
}

func (d dummySqlConn) QueryRowsPartial(v interface{}, query string, args ...interface{}) error {
	return nil
}

func (d dummySqlConn) RawDB() (*sql.DB, error) {
	return nil, nil
}

func (d dummySqlConn) Transact(func(session sqlx.Session) error) error {
	return nil
}

type trackedConn struct {
	dummySqlConn
	execValue      bool
	queryRowsValue bool
	transactValue  bool
}

func (c *trackedConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.ExecCtx(context.Background(), query, args...)
}

func (c *trackedConn) ExecCtx(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	c.execValue = true
	return c.dummySqlConn.ExecCtx(ctx, query, args...)
}

func (c *trackedConn) QueryRows(v interface{}, query string, args ...interface{}) error {
	return c.QueryRowsCtx(context.Background(), v, query, args...)
}

func (c *trackedConn) QueryRowsCtx(ctx context.Context, v interface{}, query string, args ...interface{}) error {
	c.queryRowsValue = true
	return c.dummySqlConn.QueryRowsCtx(ctx, v, query, args...)
}

func (c *trackedConn) RawDB() (*sql.DB, error) {
	return nil, nil
}

func (c *trackedConn) Transact(fn func(session sqlx.Session) error) error {
	return c.TransactCtx(context.Background(), func(_ context.Context, session sqlx.Session) error {
		return fn(session)
	})
}

func (c *trackedConn) TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	c.transactValue = true
	return c.dummySqlConn.TransactCtx(ctx, fn)
}
