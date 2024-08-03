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

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/fx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stat"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/dbtest"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/redis/redistest"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
)

func init() {
	logx.Disable()
	stat.SetReporter(nil)
}

func TestCachedConn_GetCache(t *testing.T) {
	resetStats()
	r := redistest.CreateRedis(t)

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))
	var value string
	err := c.GetCache("any", &value)
	assert.Equal(t, ErrNotFound, err)
	_ = r.Set("any", `"value"`)
	err = c.GetCache("any", &value)
	assert.Nil(t, err)
	assert.Equal(t, "value", value)
}

func TestStat(t *testing.T) {
	resetStats()
	r := redistest.CreateRedis(t)

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))

	for i := 0; i < 10; i++ {
		var str string
		err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v any) error {
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
	r := redistest.CreateRedis(t)

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
	err := c.QueryRowIndex(&str, "index", func(s any) string {
		return fmt.Sprintf("%s/1234", s)
	}, func(conn sqlx.SqlConn, v any) (any, error) {
		*v.(*string) = "zero"
		return "primary", errors.New("foo")
	}, func(conn sqlx.SqlConn, v, pri any) error {
		assert.Equal(t, "primary", pri)
		*v.(*string) = "xin"
		return nil
	})
	assert.NotNil(t, err)

	err = c.QueryRowIndex(&str, "index", func(s any) string {
		return fmt.Sprintf("%s/1234", s)
	}, func(conn sqlx.SqlConn, v any) (any, error) {
		*v.(*string) = "zero"
		return "primary", nil
	}, func(conn sqlx.SqlConn, v, pri any) error {
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
	r := redistest.CreateRedis(t)

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10),
		cache.WithNotFoundExpiry(time.Second))

	var str string
	r.Set("index", `"primary"`)
	err := c.QueryRowIndex(&str, "index", func(s any) string {
		return fmt.Sprintf("%s/1234", s)
	}, func(conn sqlx.SqlConn, v any) (any, error) {
		assert.Fail(t, "should not go here")
		return "primary", nil
	}, func(conn sqlx.SqlConn, v, primary any) error {
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
		primary      any
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
			r := redistest.CreateRedis(t)

			c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10),
				cache.WithNotFoundExpiry(time.Second))

			var str string
			r.Set("index", test.primaryCache)
			err := c.QueryRowIndex(&str, "index", func(s any) string {
				return fmt.Sprintf("%v/1234", s)
			}, func(conn sqlx.SqlConn, v any) (any, error) {
				assert.Fail(t, "should not go here")
				return test.primary, nil
			}, func(conn sqlx.SqlConn, v, primary any) error {
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
			r := redistest.CreateRedis(t)

			c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10),
				cache.WithNotFoundExpiry(time.Second))

			var str string
			r.Set(k, v)
			err := c.QueryRowIndex(&str, "index", func(s any) string {
				return fmt.Sprintf("%s/1234", s)
			}, func(conn sqlx.SqlConn, v any) (any, error) {
				*v.(*string) = "xin"
				return "primary", nil
			}, func(conn sqlx.SqlConn, v, primary any) error {
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
		err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v any) error {
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
	r := redistest.CreateRedis(t)

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))

	for i := 0; i < 20; i++ {
		var str string
		err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v any) error {
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
	r := redistest.CreateRedis(t)

	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))

	var all sync.WaitGroup
	var wait sync.WaitGroup
	all.Add(10)
	wait.Add(4)
	go func() {
		var str string
		err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v any) error {
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
			err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v any) error {
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
			err := c.QueryRow(&str, "name", func(conn sqlx.SqlConn, v any) error {
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

func TestCachedConn_DelCache(t *testing.T) {
	r := redistest.CreateRedis(t)

	const (
		key   = "user"
		value = "any"
	)
	assert.NoError(t, r.Set(key, value))

	c := NewNodeConn(&trackedConn{}, r, cache.WithExpiry(time.Second*30))
	err := c.DelCache(key)
	assert.Nil(t, err)

	val, err := r.Get(key)
	assert.Nil(t, err)
	assert.Empty(t, val)
}

func TestCachedConnQueryRow(t *testing.T) {
	r := redistest.CreateRedis(t)

	const (
		key   = "user"
		value = "any"
	)
	var conn trackedConn
	var user string
	var ran bool
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	err := c.QueryRow(&user, key, func(conn sqlx.SqlConn, v any) error {
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
	r := redistest.CreateRedis(t)

	const (
		key   = "user"
		value = "any"
	)
	var conn trackedConn
	var user string
	var ran bool
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	assert.Nil(t, c.SetCache(key, value))
	err := c.QueryRow(&user, key, func(conn sqlx.SqlConn, v any) error {
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
	r := redistest.CreateRedis(t)

	const key = "user"
	var conn trackedConn
	var user string
	var ran int
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	for i := 0; i < 20; i++ {
		err := c.QueryRow(&user, key, func(conn sqlx.SqlConn, v any) error {
			ran++
			return sql.ErrNoRows
		})
		assert.Exactly(t, sqlx.ErrNotFound, err)
	}
	assert.Equal(t, 1, ran)
}

func TestCachedConnExec(t *testing.T) {
	r := redistest.CreateRedis(t)

	var conn trackedConn
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	_, err := c.ExecNoCache("delete from user_table where id='kevin'")
	assert.Nil(t, err)
	assert.True(t, conn.execValue)
}

func TestCachedConnExecDropCache(t *testing.T) {
	t.Run("drop cache", func(t *testing.T) {
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
	})
}

func TestCachedConn_SetCacheWithExpire(t *testing.T) {
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
	assert.Nil(t, c.SetCacheWithExpire(key, value, time.Minute))
	val, err := r.Get(key)
	if assert.NoError(t, err) {
		ttl := r.TTL(key)
		assert.True(t, ttl > 0 && ttl <= time.Minute)
		assert.Equal(t, fmt.Sprintf("%q", value), val)
	}
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
	r := redistest.CreateRedis(t)

	var conn trackedConn
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	var users []string
	err := c.QueryRowsNoCache(&users, "select user from user_table where id='kevin'")
	assert.Nil(t, err)
	assert.True(t, conn.queryRowsValue)
}

func TestCachedConnTransact(t *testing.T) {
	r := redistest.CreateRedis(t)

	var conn trackedConn
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	err := c.Transact(func(session sqlx.Session) error {
		return nil
	})
	assert.Nil(t, err)
	assert.True(t, conn.transactValue)
}

func TestQueryRowNoCache(t *testing.T) {
	r := redistest.CreateRedis(t)

	const (
		key   = "user"
		value = "any"
	)
	var user string
	var ran bool
	conn := dummySqlConn{queryRow: func(v any, q string, args ...any) error {
		user = value
		ran = true
		return nil
	}}
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	err := c.QueryRowNoCache(&user, key)
	assert.Nil(t, err)
	assert.Equal(t, value, user)
	assert.True(t, ran)
}

func TestQueryRowPartialNoCache(t *testing.T) {
	r := redistest.CreateRedis(t)

	const (
		key   = "user"
		value = "any"
	)
	var user string
	var ran bool
	conn := dummySqlConn{queryRow: func(v any, q string, args ...any) error {
		user = value
		ran = true
		return nil
	}}
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	err := c.QueryRowPartialNoCache(&user, key)
	assert.Nil(t, err)
	assert.Equal(t, value, user)
	assert.True(t, ran)
}

func TestQueryRowsPartialNoCache(t *testing.T) {
	r := redistest.CreateRedis(t)

	var (
		key    = "user"
		values = []string{"any", "any"}
	)
	var users []string
	var ran bool
	conn := dummySqlConn{queryRows: func(v any, q string, args ...any) error {
		users = values
		ran = true
		return nil
	}}
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	err := c.QueryRowsPartialNoCache(&users, key)
	assert.Nil(t, err)
	assert.Equal(t, values, users)
	assert.True(t, ran)
}

func TestNewConnWithCache(t *testing.T) {
	r := redistest.CreateRedis(t)

	var conn trackedConn
	c := NewConnWithCache(&conn, cache.NewNode(r, singleFlights, stats, sql.ErrNoRows))
	_, err := c.ExecNoCache("delete from user_table where id='kevin'")
	assert.Nil(t, err)
	assert.True(t, conn.execValue)
}

func TestCachedConn_WithSession(t *testing.T) {
	dbtest.RunTxTest(t, func(tx *sql.Tx, mock sqlmock.Sqlmock) {
		mock.ExpectExec("any").WillReturnResult(sqlmock.NewResult(2, 3))

		r := redistest.CreateRedis(t)
		conn := CachedConn{
			cache: cache.NewNode(r, syncx.NewSingleFlight(), stats, sql.ErrNoRows),
		}
		conn = conn.WithSession(sqlx.NewSessionFromTx(tx))
		res, err := conn.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
			return conn.Exec("any")
		}, "foo")
		assert.NoError(t, err)
		last, err := res.LastInsertId()
		assert.NoError(t, err)
		assert.Equal(t, int64(2), last)
		affected, err := res.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(3), affected)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("any").WillReturnResult(sqlmock.NewResult(2, 3))
		mock.ExpectCommit()

		r := redistest.CreateRedis(t)
		conn := CachedConn{
			db:    sqlx.NewSqlConnFromDB(db),
			cache: cache.NewNode(r, syncx.NewSingleFlight(), stats, sql.ErrNoRows),
		}
		assert.NoError(t, conn.Transact(func(session sqlx.Session) error {
			conn = conn.WithSession(session)
			res, err := conn.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
				return conn.Exec("any")
			}, "foo")
			assert.NoError(t, err)
			last, err := res.LastInsertId()
			assert.NoError(t, err)
			assert.Equal(t, int64(2), last)
			affected, err := res.RowsAffected()
			assert.NoError(t, err)
			assert.Equal(t, int64(3), affected)
			return nil
		}))
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("any").WillReturnError(errors.New("foo"))
		mock.ExpectRollback()

		r := redistest.CreateRedis(t)
		conn := CachedConn{
			db:    sqlx.NewSqlConnFromDB(db),
			cache: cache.NewNode(r, syncx.NewSingleFlight(), stats, sql.ErrNoRows),
		}
		assert.Error(t, conn.Transact(func(session sqlx.Session) error {
			conn = conn.WithSession(session)
			_, err := conn.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
				return conn.Exec("any")
			}, "bar")
			return err
		}))
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectQuery("any").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
		mock.ExpectCommit()

		r := redistest.CreateRedis(t)
		conn := CachedConn{
			db:    sqlx.NewSqlConnFromDB(db),
			cache: cache.NewNode(r, syncx.NewSingleFlight(), stats, sql.ErrNoRows),
		}
		assert.NoError(t, conn.Transact(func(session sqlx.Session) error {
			var val string
			conn = conn.WithSession(session)
			err := conn.QueryRow(&val, "foo", func(conn sqlx.SqlConn, v interface{}) error {
				return conn.QueryRow(v, "any")
			})
			assert.Equal(t, "2", val)
			return err
		}))
		val, err := r.Get("foo")
		assert.NoError(t, err)
		assert.Equal(t, `"2"`, val)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectQuery("any").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
		mock.ExpectExec("any").WillReturnResult(sqlmock.NewResult(2, 3))
		mock.ExpectCommit()

		r := redistest.CreateRedis(t)
		conn := CachedConn{
			db:    sqlx.NewSqlConnFromDB(db),
			cache: cache.NewNode(r, syncx.NewSingleFlight(), stats, sql.ErrNoRows),
		}
		assert.NoError(t, conn.Transact(func(session sqlx.Session) error {
			var val string
			conn = conn.WithSession(session)
			assert.NoError(t, conn.QueryRow(&val, "foo", func(conn sqlx.SqlConn, v interface{}) error {
				return conn.QueryRow(v, "any")
			}))
			assert.Equal(t, "2", val)
			_, err := conn.Exec(func(conn sqlx.SqlConn) (sql.Result, error) {
				return conn.Exec("any")
			}, "foo")
			return err
		}))
		val, err := r.Get("foo")
		assert.NoError(t, err)
		assert.Empty(t, val)
	})
}

func resetStats() {
	atomic.StoreUint64(&stats.Total, 0)
	atomic.StoreUint64(&stats.Hit, 0)
	atomic.StoreUint64(&stats.Miss, 0)
	atomic.StoreUint64(&stats.DbFails, 0)
}

type dummySqlConn struct {
	queryRow  func(any, string, ...any) error
	queryRows func(any, string, ...any) error
}

func (d dummySqlConn) ExecCtx(_ context.Context, _ string, _ ...any) (sql.Result, error) {
	return nil, nil
}

func (d dummySqlConn) PrepareCtx(_ context.Context, _ string) (sqlx.StmtSession, error) {
	return nil, nil
}

func (d dummySqlConn) QueryRowPartialCtx(_ context.Context, v any, query string, args ...any) error {
	if d.queryRow != nil {
		return d.queryRow(v, query, args...)
	}

	return nil
}

func (d dummySqlConn) QueryRowsCtx(_ context.Context, _ any, _ string, _ ...any) error {
	return nil
}

func (d dummySqlConn) QueryRowsPartialCtx(_ context.Context, v any, query string, args ...any) error {
	if d.queryRows != nil {
		return d.queryRows(v, query, args...)
	}

	return nil
}

func (d dummySqlConn) TransactCtx(_ context.Context, _ func(context.Context, sqlx.Session) error) error {
	return nil
}

func (d dummySqlConn) Exec(_ string, _ ...any) (sql.Result, error) {
	return nil, nil
}

func (d dummySqlConn) Prepare(_ string) (sqlx.StmtSession, error) {
	return nil, nil
}

func (d dummySqlConn) QueryRow(v any, query string, args ...any) error {
	return d.QueryRowCtx(context.Background(), v, query, args...)
}

func (d dummySqlConn) QueryRowCtx(_ context.Context, v any, query string, args ...any) error {
	if d.queryRow != nil {
		return d.queryRow(v, query, args...)
	}
	return nil
}

func (d dummySqlConn) QueryRowPartial(_ any, _ string, _ ...any) error {
	return nil
}

func (d dummySqlConn) QueryRows(_ any, _ string, _ ...any) error {
	return nil
}

func (d dummySqlConn) QueryRowsPartial(_ any, _ string, _ ...any) error {
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

func (c *trackedConn) Exec(query string, args ...any) (sql.Result, error) {
	return c.ExecCtx(context.Background(), query, args...)
}

func (c *trackedConn) ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error) {
	c.execValue = true
	return c.dummySqlConn.ExecCtx(ctx, query, args...)
}

func (c *trackedConn) QueryRows(v any, query string, args ...any) error {
	return c.QueryRowsCtx(context.Background(), v, query, args...)
}

func (c *trackedConn) QueryRowsCtx(ctx context.Context, v any, query string, args ...any) error {
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
