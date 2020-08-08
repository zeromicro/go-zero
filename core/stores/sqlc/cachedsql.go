package sqlc

import (
	"database/sql"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/internal"
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/syncx"
)

// see doc/sql-cache.md
const cacheSafeGapBetweenIndexAndPrimary = time.Second * 5

var (
	ErrNotFound = sqlx.ErrNotFound

	// can't use one SharedCalls per conn, because multiple conns may share the same cache key.
	exclusiveCalls = syncx.NewSharedCalls()
	stats          = internal.NewCacheStat("sqlc")
)

type (
	ExecFn         func(conn sqlx.SqlConn) (sql.Result, error)
	IndexQueryFn   func(conn sqlx.SqlConn, v interface{}) (interface{}, error)
	PrimaryQueryFn func(conn sqlx.SqlConn, v, primary interface{}) error
	QueryFn        func(conn sqlx.SqlConn, v interface{}) error

	CachedConn struct {
		db    sqlx.SqlConn
		cache internal.Cache
	}
)

func NewNodeConn(db sqlx.SqlConn, rds *redis.Redis, opts ...cache.Option) CachedConn {
	return CachedConn{
		db:    db,
		cache: internal.NewCacheNode(rds, exclusiveCalls, stats, sql.ErrNoRows, opts...),
	}
}

func NewConn(db sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CachedConn {
	return CachedConn{
		db:    db,
		cache: internal.NewCache(c, exclusiveCalls, stats, sql.ErrNoRows, opts...),
	}
}

func (cc CachedConn) DelCache(keys ...string) error {
	return cc.cache.DelCache(keys...)
}

func (cc CachedConn) GetCache(key string, v interface{}) error {
	return cc.cache.GetCache(key, v)
}

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

func (cc CachedConn) ExecNoCache(q string, args ...interface{}) (sql.Result, error) {
	return cc.db.Exec(q, args...)
}

func (cc CachedConn) QueryRow(v interface{}, key string, query QueryFn) error {
	return cc.cache.Take(v, key, func(v interface{}) error {
		return query(cc.db, v)
	})
}

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
		return cc.cache.SetCacheWithExpire(keyer(primaryKey), v, expire+cacheSafeGapBetweenIndexAndPrimary)
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

func (cc CachedConn) QueryRowNoCache(v interface{}, q string, args ...interface{}) error {
	return cc.db.QueryRow(v, q, args...)
}

// QueryRowsNoCache doesn't use cache, because it might cause consistency problem.
func (cc CachedConn) QueryRowsNoCache(v interface{}, q string, args ...interface{}) error {
	return cc.db.QueryRows(v, q, args...)
}

func (cc CachedConn) SetCache(key string, v interface{}) error {
	return cc.cache.SetCache(key, v)
}

func (cc CachedConn) Transact(fn func(sqlx.Session) error) error {
	return cc.db.Transact(fn)
}
