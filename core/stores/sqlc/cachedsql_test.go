package sqlc

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/stat"
	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

func init() {
	logx.Disable()
	stat.SetReporter(nil)
}

func TestCachedConn_GetCache(t *testing.T) {
	resetStats()
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	r := redis.NewRedis(s.Addr(), redis.NodeType)
	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))
	var value string
	err = c.GetCache("any", &value)
	assert.Equal(t, ErrNotFound, err)
	s.Set("any", `"value"`)
	err = c.GetCache("any", &value)
	assert.Nil(t, err)
	assert.Equal(t, "value", value)
}

func TestStat(t *testing.T) {
	resetStats()
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	r := redis.NewRedis(s.Addr(), redis.NodeType)
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
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	r := redis.NewRedis(s.Addr(), redis.NodeType)
	c := NewNodeConn(dummySqlConn{}, r, cache.WithExpiry(time.Second*10))

	var str string
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
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	r := redis.NewRedis(s.Addr(), redis.NodeType)
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

func TestCachedConn_QueryRowIndex_HasWrongCache(t *testing.T) {
	caches := map[string]string{
		"index":        "primary",
		"primary/1234": "xin",
	}

	for k, v := range caches {
		t.Run(k+"/"+v, func(t *testing.T) {
			resetStats()
			s, err := miniredis.Run()
			if err != nil {
				t.Error(err)
			}

			r := redis.NewRedis(s.Addr(), redis.NodeType)
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
	log.SetOutput(ioutil.Discard)
	defer log.SetOutput(os.Stdout)

	r := redis.NewRedis("localhost:59999", redis.NodeType)
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
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	r := redis.NewRedis(s.Addr(), redis.NodeType)
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
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	r := redis.NewRedis(s.Addr(), redis.NodeType)
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
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	const (
		key   = "user"
		value = "any"
	)
	var conn trackedConn
	var user string
	var ran bool
	r := redis.NewRedis(s.Addr(), redis.NodeType)
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	err = c.QueryRow(&user, key, func(conn sqlx.SqlConn, v interface{}) error {
		ran = true
		user = value
		return nil
	})
	assert.Nil(t, err)
	actualValue, err := s.Get(key)
	assert.Nil(t, err)
	var actual string
	assert.Nil(t, json.Unmarshal([]byte(actualValue), &actual))
	assert.Equal(t, value, actual)
	assert.Equal(t, value, user)
	assert.True(t, ran)
}

func TestCachedConnQueryRowFromCache(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	const (
		key   = "user"
		value = "any"
	)
	var conn trackedConn
	var user string
	var ran bool
	r := redis.NewRedis(s.Addr(), redis.NodeType)
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	assert.Nil(t, c.SetCache(key, value))
	err = c.QueryRow(&user, key, func(conn sqlx.SqlConn, v interface{}) error {
		ran = true
		user = value
		return nil
	})
	assert.Nil(t, err)
	actualValue, err := s.Get(key)
	assert.Nil(t, err)
	var actual string
	assert.Nil(t, json.Unmarshal([]byte(actualValue), &actual))
	assert.Equal(t, value, actual)
	assert.Equal(t, value, user)
	assert.False(t, ran)
}

func TestQueryRowNotFound(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	const key = "user"
	var conn trackedConn
	var user string
	var ran int
	r := redis.NewRedis(s.Addr(), redis.NodeType)
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
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	var conn trackedConn
	r := redis.NewRedis(s.Addr(), redis.NodeType)
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	_, err = c.ExecNoCache("delete from user_table where id='kevin'")
	assert.Nil(t, err)
	assert.True(t, conn.execValue)
}

func TestCachedConnExecDropCache(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	const (
		key   = "user"
		value = "any"
	)
	var conn trackedConn
	r := redis.NewRedis(s.Addr(), redis.NodeType)
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*30))
	assert.Nil(t, c.SetCache(key, value))
	_, err = c.Exec(func(conn sqlx.SqlConn) (result sql.Result, e error) {
		return conn.Exec("delete from user_table where id='kevin'")
	}, key)
	assert.Nil(t, err)
	assert.True(t, conn.execValue)
	_, err = s.Get(key)
	assert.Exactly(t, miniredis.ErrKeyNotFound, err)
}

func TestCachedConnExecDropCacheFailed(t *testing.T) {
	const key = "user"
	var conn trackedConn
	r := redis.NewRedis("anyredis:8888", redis.NodeType)
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	_, err := c.Exec(func(conn sqlx.SqlConn) (result sql.Result, e error) {
		return conn.Exec("delete from user_table where id='kevin'")
	}, key)
	// async background clean, retry logic
	assert.Nil(t, err)
}

func TestCachedConnQueryRows(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	var conn trackedConn
	r := redis.NewRedis(s.Addr(), redis.NodeType)
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	var users []string
	err = c.QueryRowsNoCache(&users, "select user from user_table where id='kevin'")
	assert.Nil(t, err)
	assert.True(t, conn.queryRowsValue)
}

func TestCachedConnTransact(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error(err)
	}

	var conn trackedConn
	r := redis.NewRedis(s.Addr(), redis.NodeType)
	c := NewNodeConn(&conn, r, cache.WithExpiry(time.Second*10))
	err = c.Transact(func(session sqlx.Session) error {
		return nil
	})
	assert.Nil(t, err)
	assert.True(t, conn.transactValue)
}

func resetStats() {
	atomic.StoreUint64(&stats.Total, 0)
	atomic.StoreUint64(&stats.Hit, 0)
	atomic.StoreUint64(&stats.Miss, 0)
	atomic.StoreUint64(&stats.DbFails, 0)
}

type dummySqlConn struct {
}

func (d dummySqlConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (d dummySqlConn) Prepare(query string) (sqlx.StmtSession, error) {
	return nil, nil
}

func (d dummySqlConn) QueryRow(v interface{}, query string, args ...interface{}) error {
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
	c.execValue = true
	return c.dummySqlConn.Exec(query, args...)
}

func (c *trackedConn) QueryRows(v interface{}, query string, args ...interface{}) error {
	c.queryRowsValue = true
	return c.dummySqlConn.QueryRows(v, query, args...)
}

func (c *trackedConn) Transact(fn func(session sqlx.Session) error) error {
	c.transactValue = true
	return c.dummySqlConn.Transact(fn)
}
