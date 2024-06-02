package sqlx

import (
	"context"
	"database/sql"
	"errors"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"
)

// spanName is used to identify the span name for the SQL execution.
const spanName = "sql"

type (
	// Session stands for raw connections or transaction sessions
	Session interface {
		Exec(query string, args ...any) (sql.Result, error)
		ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error)
		Prepare(query string) (StmtSession, error)
		PrepareCtx(ctx context.Context, query string) (StmtSession, error)
		QueryRow(v any, query string, args ...any) error
		QueryRowCtx(ctx context.Context, v any, query string, args ...any) error
		QueryRowPartial(v any, query string, args ...any) error
		QueryRowPartialCtx(ctx context.Context, v any, query string, args ...any) error
		QueryRows(v any, query string, args ...any) error
		QueryRowsCtx(ctx context.Context, v any, query string, args ...any) error
		QueryRowsPartial(v any, query string, args ...any) error
		QueryRowsPartialCtx(ctx context.Context, v any, query string, args ...any) error
	}

	// SqlConn only stands for raw connections, so Transact method can be called.
	SqlConn interface {
		Session
		// RawDB is for other ORM to operate with, use it with caution.
		// Notice: don't close it.
		RawDB() (*sql.DB, error)
		Transact(fn func(Session) error) error
		TransactCtx(ctx context.Context, fn func(context.Context, Session) error) error
	}

	// SqlOption defines the method to customize a sql connection.
	SqlOption func(*commonSqlConn)

	// thread-safe
	// Because CORBA doesn't support PREPARE, so we need to combine the
	// query arguments into one string and do underlying query without arguments
	commonSqlConn struct {
		connProv connProvider
		onError  func(context.Context, error)
		beginTx  beginnable
		brk      breaker.Breaker
		accept   breaker.Acceptable
	}

	connProvider func() (*sql.DB, error)

	sessionConn interface {
		Exec(query string, args ...any) (sql.Result, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		Query(query string, args ...any) (*sql.Rows, error)
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	}
)

// NewSqlConn returns a SqlConn with given driver name and datasource.
func NewSqlConn(driverName, datasource string, opts ...SqlOption) SqlConn {
	conn := &commonSqlConn{
		connProv: func() (*sql.DB, error) {
			return getSqlConn(driverName, datasource)
		},
		onError: func(ctx context.Context, err error) {
			logInstanceError(ctx, datasource, err)
		},
		beginTx: begin,
		brk:     breaker.NewBreaker(),
	}
	for _, opt := range opts {
		opt(conn)
	}

	return conn
}

// NewSqlConnFromDB returns a SqlConn with the given sql.DB.
// Use it with caution; it's provided for other ORM to interact with.
func NewSqlConnFromDB(db *sql.DB, opts ...SqlOption) SqlConn {
	conn := &commonSqlConn{
		connProv: func() (*sql.DB, error) {
			return db, nil
		},
		onError: func(ctx context.Context, err error) {
			logx.WithContext(ctx).Errorf("Error on getting sql instance: %v", err)
		},
		beginTx: begin,
		brk:     breaker.NewBreaker(),
	}
	for _, opt := range opts {
		opt(conn)
	}

	return conn
}

// NewSqlConnFromSession returns a SqlConn with the given session.
func NewSqlConnFromSession(session Session) SqlConn {
	return txConn{
		Session: session,
	}
}

func (db *commonSqlConn) Exec(q string, args ...any) (result sql.Result, err error) {
	return db.ExecCtx(context.Background(), q, args...)
}

func (db *commonSqlConn) ExecCtx(ctx context.Context, q string, args ...any) (
	result sql.Result, err error) {
	ctx, span := startSpan(ctx, "Exec")
	defer func() {
		endSpan(span, err)
	}()

	err = db.brk.DoWithAcceptableCtx(ctx, func() error {
		var conn *sql.DB
		conn, err = db.connProv()
		if err != nil {
			db.onError(ctx, err)
			return err
		}

		result, err = exec(ctx, conn, q, args...)
		return err
	}, db.acceptable)
	if errors.Is(err, breaker.ErrServiceUnavailable) {
		metricReqErr.Inc("Exec", "breaker")
	}

	return
}

func (db *commonSqlConn) Prepare(query string) (stmt StmtSession, err error) {
	return db.PrepareCtx(context.Background(), query)
}

func (db *commonSqlConn) PrepareCtx(ctx context.Context, query string) (stmt StmtSession, err error) {
	ctx, span := startSpan(ctx, "Prepare")
	defer func() {
		endSpan(span, err)
	}()

	err = db.brk.DoWithAcceptableCtx(ctx, func() error {
		var conn *sql.DB
		conn, err = db.connProv()
		if err != nil {
			db.onError(ctx, err)
			return err
		}

		st, err := conn.PrepareContext(ctx, query)
		if err != nil {
			return err
		}

		stmt = statement{
			query:  query,
			stmt:   st,
			brk:    db.brk,
			accept: db.acceptable,
		}
		return nil
	}, db.acceptable)
	if errors.Is(err, breaker.ErrServiceUnavailable) {
		metricReqErr.Inc("Prepare", "breaker")
	}

	return
}

func (db *commonSqlConn) QueryRow(v any, q string, args ...any) error {
	return db.QueryRowCtx(context.Background(), v, q, args...)
}

func (db *commonSqlConn) QueryRowCtx(ctx context.Context, v any, q string,
	args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRow")
	defer func() {
		endSpan(span, err)
	}()

	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, q, args...)
}

func (db *commonSqlConn) QueryRowPartial(v any, q string, args ...any) error {
	return db.QueryRowPartialCtx(context.Background(), v, q, args...)
}

func (db *commonSqlConn) QueryRowPartialCtx(ctx context.Context, v any,
	q string, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowPartial")
	defer func() {
		endSpan(span, err)
	}()

	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, q, args...)
}

func (db *commonSqlConn) QueryRows(v any, q string, args ...any) error {
	return db.QueryRowsCtx(context.Background(), v, q, args...)
}

func (db *commonSqlConn) QueryRowsCtx(ctx context.Context, v any, q string,
	args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRows")
	defer func() {
		endSpan(span, err)
	}()

	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, q, args...)
}

func (db *commonSqlConn) QueryRowsPartial(v any, q string, args ...any) error {
	return db.QueryRowsPartialCtx(context.Background(), v, q, args...)
}

func (db *commonSqlConn) QueryRowsPartialCtx(ctx context.Context, v any,
	q string, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowsPartial")
	defer func() {
		endSpan(span, err)
	}()

	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

func (db *commonSqlConn) RawDB() (*sql.DB, error) {
	return db.connProv()
}

func (db *commonSqlConn) Transact(fn func(Session) error) error {
	return db.TransactCtx(context.Background(), func(_ context.Context, session Session) error {
		return fn(session)
	})
}

func (db *commonSqlConn) TransactCtx(ctx context.Context, fn func(context.Context, Session) error) (err error) {
	ctx, span := startSpan(ctx, "Transact")
	defer func() {
		endSpan(span, err)
	}()

	err = db.brk.DoWithAcceptableCtx(ctx, func() error {
		return transact(ctx, db, db.beginTx, fn)
	}, db.acceptable)
	if errors.Is(err, breaker.ErrServiceUnavailable) {
		metricReqErr.Inc("Transact", "breaker")
	}

	return
}

func (db *commonSqlConn) acceptable(err error) bool {
	if err == nil || errorx.In(err, sql.ErrNoRows, sql.ErrTxDone, context.Canceled) {
		return true
	}

	var e acceptableError
	if errors.As(err, &e) {
		return true
	}

	if db.accept == nil {
		return false
	}

	return db.accept(err)
}

func (db *commonSqlConn) queryRows(ctx context.Context, scanner func(*sql.Rows) error,
	q string, args ...any) (err error) {
	var scanFailed bool
	err = db.brk.DoWithAcceptableCtx(ctx, func() error {
		conn, err := db.connProv()
		if err != nil {
			db.onError(ctx, err)
			return err
		}

		return query(ctx, conn, func(rows *sql.Rows) error {
			e := scanner(rows)
			if isScanFailed(e) {
				scanFailed = true
			}
			return e
		}, q, args...)
	}, func(err error) bool {
		return scanFailed || db.acceptable(err)
	})
	if errors.Is(err, breaker.ErrServiceUnavailable) {
		metricReqErr.Inc("queryRows", "breaker")
	}

	return
}

// WithAcceptable returns a SqlOption that setting the acceptable function.
// acceptable is the func to check if the error can be accepted.
func WithAcceptable(acceptable func(err error) bool) SqlOption {
	return func(conn *commonSqlConn) {
		if conn.accept == nil {
			conn.accept = acceptable
		} else {
			pre := conn.accept
			conn.accept = func(err error) bool {
				return pre(err) || acceptable(err)
			}
		}
	}
}
