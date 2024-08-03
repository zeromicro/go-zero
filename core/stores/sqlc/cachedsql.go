package sqlc

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/syncx"
)

// see doc/sql-cache.md
const cacheSafeGapBetweenIndexAndPrimary = time.Second * 5

var (
	// ErrNotFound is an alias of sqlx.ErrNotFound.
	ErrNotFound = sqlx.ErrNotFound

	// can't use one SingleFlight per conn, because multiple conns may share the same cache key.
	singleFlights = syncx.NewSingleFlight()
	stats         = cache.NewStat("sqlc")
)

type (
	// ExecFn defines the sql exec method.
	ExecFn func(conn sqlx.SqlConn) (sql.Result, error)
	// ExecCtxFn defines the sql exec method.
	ExecCtxFn func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error)
	// IndexQueryFn defines the query method that based on unique indexes.
	IndexQueryFn func(conn sqlx.SqlConn, v any) (any, error)
	// IndexQueryCtxFn defines the query method that based on unique indexes.
	IndexQueryCtxFn func(ctx context.Context, conn sqlx.SqlConn, v any) (any, error)
	// PrimaryQueryFn defines the query method that based on primary keys.
	PrimaryQueryFn func(conn sqlx.SqlConn, v, primary any) error
	// PrimaryQueryCtxFn defines the query method that based on primary keys.
	PrimaryQueryCtxFn func(ctx context.Context, conn sqlx.SqlConn, v, primary any) error
	// QueryFn defines the query method.
	QueryFn func(conn sqlx.SqlConn, v any) error
	// QueryCtxFn defines the query method.
	QueryCtxFn func(ctx context.Context, conn sqlx.SqlConn, v any) error

	// A CachedConn is a DB connection with cache capability.
	CachedConn struct {
		db    sqlx.SqlConn
		cache cache.Cache
	}
)

// NewConn returns a CachedConn with a redis cluster cache.
func NewConn(db sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CachedConn {
	cc := cache.New(c, singleFlights, stats, sql.ErrNoRows, opts...)
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
	c := cache.NewNode(rds, singleFlights, stats, sql.ErrNoRows, opts...)
	return NewConnWithCache(db, c)
}

// DelCache deletes cache with keys.
func (cc CachedConn) DelCache(keys ...string) error {
	return cc.DelCacheCtx(context.Background(), keys...)
}

// DelCacheCtx deletes cache with keys.
func (cc CachedConn) DelCacheCtx(ctx context.Context, keys ...string) error {
	return cc.cache.DelCtx(ctx, keys...)
}

// GetCache unmarshals cache with given key into v.
func (cc CachedConn) GetCache(key string, v any) error {
	return cc.GetCacheCtx(context.Background(), key, v)
}

// GetCacheCtx unmarshals cache with given key into v.
func (cc CachedConn) GetCacheCtx(ctx context.Context, key string, v any) error {
	return cc.cache.GetCtx(ctx, key, v)
}

// Exec runs given exec on given keys, and returns execution result.
func (cc CachedConn) Exec(exec ExecFn, keys ...string) (sql.Result, error) {
	execCtx := func(_ context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return exec(conn)
	}
	return cc.ExecCtx(context.Background(), execCtx, keys...)
}

// ExecCtx runs given exec on given keys, and returns execution result.
// If DB operation succeeds, it will delete cache with given keys,
// if DB operation fails, it will return nil result and non-nil error,
// if DB operation succeeds but cache deletion fails, it will return result and non-nil error.
func (cc CachedConn) ExecCtx(ctx context.Context, exec ExecCtxFn, keys ...string) (
	sql.Result, error) {
	res, err := exec(ctx, cc.db)
	if err != nil {
		return nil, err
	}

	return res, cc.DelCacheCtx(ctx, keys...)
}

// ExecNoCache runs exec with given sql statement, without affecting cache.
func (cc CachedConn) ExecNoCache(q string, args ...any) (sql.Result, error) {
	return cc.ExecNoCacheCtx(context.Background(), q, args...)
}

// ExecNoCacheCtx runs exec with given sql statement, without affecting cache.
func (cc CachedConn) ExecNoCacheCtx(ctx context.Context, q string, args ...any) (
	sql.Result, error) {
	return cc.db.ExecCtx(ctx, q, args...)
}

// QueryRow unmarshals into v with given key and query func.
func (cc CachedConn) QueryRow(v any, key string, query QueryFn) error {
	queryCtx := func(_ context.Context, conn sqlx.SqlConn, v any) error {
		return query(conn, v)
	}
	return cc.QueryRowCtx(context.Background(), v, key, queryCtx)
}

// QueryRowCtx unmarshals into v with given key and query func.
func (cc CachedConn) QueryRowCtx(ctx context.Context, v any, key string, query QueryCtxFn) error {
	return cc.cache.TakeCtx(ctx, v, key, func(v any) error {
		return query(ctx, cc.db, v)
	})
}

// QueryRowIndex unmarshals into v with given key.
func (cc CachedConn) QueryRowIndex(v any, key string, keyer func(primary any) string,
	indexQuery IndexQueryFn, primaryQuery PrimaryQueryFn) error {
	indexQueryCtx := func(_ context.Context, conn sqlx.SqlConn, v any) (any, error) {
		return indexQuery(conn, v)
	}
	primaryQueryCtx := func(_ context.Context, conn sqlx.SqlConn, v, primary any) error {
		return primaryQuery(conn, v, primary)
	}

	return cc.QueryRowIndexCtx(context.Background(), v, key, keyer, indexQueryCtx, primaryQueryCtx)
}

// QueryRowIndexCtx unmarshals into v with given key.
func (cc CachedConn) QueryRowIndexCtx(ctx context.Context, v any, key string,
	keyer func(primary any) string, indexQuery IndexQueryCtxFn,
	primaryQuery PrimaryQueryCtxFn) error {
	var primaryKey any
	var found bool

	if err := cc.cache.TakeWithExpireCtx(ctx, &primaryKey, key,
		func(val any, expire time.Duration) (err error) {
			primaryKey, err = indexQuery(ctx, cc.db, v)
			if err != nil {
				return
			}

			found = true
			return cc.cache.SetWithExpireCtx(ctx, keyer(primaryKey), v,
				expire+cacheSafeGapBetweenIndexAndPrimary)
		}); err != nil {
		return err
	}

	if found {
		return nil
	}

	return cc.cache.TakeCtx(ctx, v, keyer(primaryKey), func(v any) error {
		return primaryQuery(ctx, cc.db, v, primaryKey)
	})
}

// QueryRowNoCache unmarshals into v with given statement.
func (cc CachedConn) QueryRowNoCache(v any, q string, args ...any) error {
	return cc.QueryRowNoCacheCtx(context.Background(), v, q, args...)
}

// QueryRowNoCacheCtx unmarshals into v with given statement.
func (cc CachedConn) QueryRowNoCacheCtx(ctx context.Context, v any, q string,
	args ...any) error {
	return cc.db.QueryRowCtx(ctx, v, q, args...)
}

// QueryRowPartialNoCache unmarshals into v with given statement.
func (cc CachedConn) QueryRowPartialNoCache(v any, q string, args ...any) error {
	return cc.QueryRowPartialNoCacheCtx(context.Background(), v, q, args...)
}

// QueryRowPartialNoCacheCtx unmarshals into v with given statement.
func (cc CachedConn) QueryRowPartialNoCacheCtx(ctx context.Context, v any, q string,
	args ...any) error {
	return cc.db.QueryRowPartialCtx(ctx, v, q, args...)
}

// QueryRowsNoCache unmarshals into v with given statement.
// It doesn't use cache, because it might cause consistency problem.
func (cc CachedConn) QueryRowsNoCache(v any, q string, args ...any) error {
	return cc.QueryRowsNoCacheCtx(context.Background(), v, q, args...)
}

// QueryRowsNoCacheCtx unmarshals into v with given statement.
// It doesn't use cache, because it might cause consistency problem.
func (cc CachedConn) QueryRowsNoCacheCtx(ctx context.Context, v any, q string,
	args ...any) error {
	return cc.db.QueryRowsCtx(ctx, v, q, args...)
}

// QueryRowsPartialNoCache unmarshals into v with given statement.
// It doesn't use cache, because it might cause consistency problem.
func (cc CachedConn) QueryRowsPartialNoCache(v any, q string, args ...any) error {
	return cc.QueryRowsPartialNoCacheCtx(context.Background(), v, q, args...)
}

// QueryRowsPartialNoCacheCtx unmarshals into v with given statement.
// It doesn't use cache, because it might cause consistency problem.
func (cc CachedConn) QueryRowsPartialNoCacheCtx(ctx context.Context, v any, q string,
	args ...any) error {
	return cc.db.QueryRowsPartialCtx(ctx, v, q, args...)
}

// SetCache sets v into cache with given key.
func (cc CachedConn) SetCache(key string, val any) error {
	return cc.SetCacheCtx(context.Background(), key, val)
}

// SetCacheCtx sets v into cache with given key.
func (cc CachedConn) SetCacheCtx(ctx context.Context, key string, val any) error {
	return cc.cache.SetCtx(ctx, key, val)
}

// SetCacheWithExpire sets v into cache with given key with given expire.
func (cc CachedConn) SetCacheWithExpire(key string, val any, expire time.Duration) error {
	return cc.SetCacheWithExpireCtx(context.Background(), key, val, expire)
}

// SetCacheWithExpireCtx sets v into cache with given key with given expire.
func (cc CachedConn) SetCacheWithExpireCtx(ctx context.Context, key string, val any,
	expire time.Duration) error {
	return cc.cache.SetWithExpireCtx(ctx, key, val, expire)
}

// Transact runs given fn in transaction mode.
func (cc CachedConn) Transact(fn func(sqlx.Session) error) error {
	fnCtx := func(_ context.Context, session sqlx.Session) error {
		return fn(session)
	}
	return cc.TransactCtx(context.Background(), fnCtx)
}

// TransactCtx runs given fn in transaction mode.
func (cc CachedConn) TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return cc.db.TransactCtx(ctx, fn)
}

// WithSession returns a new CachedConn with given session.
// If query from session, the uncommitted data might be returned.
// Don't query for the uncommitted data, you should just use it,
// and don't use the cache for the uncommitted data.
// Not recommend to use cache within transactions due to consistency problem.
func (cc CachedConn) WithSession(session sqlx.Session) CachedConn {
	return CachedConn{
		db:    sqlx.NewSqlConnFromSession(session),
		cache: cc.cache,
	}
}
