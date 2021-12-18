package sqlc

import (
	"database/sql"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/syncx"
)

// see doc/sql-cache.md
const cacheSafeGapBetweenIndexAndPrimary = time.Second * 5

var (
	// ErrNotFound is an alias of sqlx.ErrNotFound.
	ErrNotFound = sqlx.ErrNotFound

	// can't use one SingleFlight per conn, because multiple conns may share the same cache key.
	exclusiveCalls = syncx.NewSingleFlight()
	stats          = cache.NewStat("sqlc")
)

type (
	// ExecFn defines the sql exec method.
	ExecFn func(conn sqlx.SqlConn) (sql.Result, error)
	// IndexQueryFn defines the query method that based on unique indexes.
	IndexQueryFn func(conn sqlx.SqlConn, v interface{}) (interface{}, error)
	// PrimaryQueryFn defines the query method that based on primary keys.
	PrimaryQueryFn func(conn sqlx.SqlConn, v, primary interface{}) error
	// QueryFn defines the query method.
	QueryFn func(conn sqlx.SqlConn, v interface{}) error

	// A CachedConn is a DB connection with cache capability.
	CachedConn struct {
		db    sqlx.SqlConn
		cache cache.Cache
	}
)

// NewConn returns a CachedConn with a redis cluster cache.
func NewConn(db sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CachedConn {
	cc := cache.New(c, exclusiveCalls, stats, sql.ErrNoRows, opts...)
	return NewConnWithCache(db, cc)
}

// NewConnWithCache returns a CachedConn with a custom cache.
func NewConnWithCache(db sqlx.SqlConn, c cache.Cache) CachedConn {
	return CachedConn{
		db:    db,
		cache: c,
	}
}

// NewNodeConn returns a CachedConn with a redis node cache.
func NewNodeConn(db sqlx.SqlConn, rds *redis.Redis, opts ...cache.Option) CachedConn {
	c := cache.NewNode(rds, exclusiveCalls, stats, sql.ErrNoRows, opts...)
	return NewConnWithCache(db, c)
}

// DelCache deletes cache with keys.
func (cc CachedConn) DelCache(keys ...string) error {
	return cc.cache.Del(keys...)
}

// GetCache unmarshals cache with given key into v.
func (cc CachedConn) GetCache(key string, v interface{}) error {
	return cc.cache.Get(key, v)
}

// Exec runs given exec on given keys, and returns execution result.
func (cc CachedConn) Exec(exec ExecFn, keys ...string) (sql.Result, error) {
	res, err := exec(cc.db)
	if err != nil {
		return nil, err
	}

	if err := cc.DelCache(keys...); err != nil {
		return nil, err
	}

	return res, nil
}

// ExecNoCache runs exec with given sql statement, without affecting cache.
func (cc CachedConn) ExecNoCache(q string, args ...interface{}) (sql.Result, error) {
	return cc.db.Exec(q, args...)
}

// QueryRow unmarshals into v with given key and query func.
func (cc CachedConn) QueryRow(v interface{}, key string, query QueryFn) error {
	return cc.cache.Take(v, key, func(v interface{}) error {
		return query(cc.db, v)
	})
}

// QueryRowIndex unmarshals into v with given key.
func (cc CachedConn) QueryRowIndex(v interface{}, key string, keyer func(primary interface{}) string,
	indexQuery IndexQueryFn, primaryQuery PrimaryQueryFn) error {
	var primaryKey interface{}
	var found bool

	if err := cc.cache.TakeWithExpire(&primaryKey, key, func(val interface{}, expire time.Duration) (err error) {
		primaryKey, err = indexQuery(cc.db, v)
		if err != nil {
			return
		}

		found = true
		return cc.cache.SetWithExpire(keyer(primaryKey), v, expire+cacheSafeGapBetweenIndexAndPrimary)
	}); err != nil {
		return err
	}

	if found {
		return nil
	}

	return cc.cache.Take(v, keyer(primaryKey), func(v interface{}) error {
		return primaryQuery(cc.db, v, primaryKey)
	})
}

// QueryRowNoCache unmarshals into v with given statement.
func (cc CachedConn) QueryRowNoCache(v interface{}, q string, args ...interface{}) error {
	return cc.db.QueryRow(v, q, args...)
}

// QueryRowsNoCache unmarshals into v with given statement.
// It doesn't use cache, because it might cause consistency problem.
func (cc CachedConn) QueryRowsNoCache(v interface{}, q string, args ...interface{}) error {
	return cc.db.QueryRows(v, q, args...)
}

// SetCache sets v into cache with given key.
func (cc CachedConn) SetCache(key string, v interface{}) error {
	return cc.cache.Set(key, v)
}

// Transact runs given fn in transaction mode.
func (cc CachedConn) Transact(fn func(sqlx.Session) error) error {
	return cc.db.Transact(fn)
}
