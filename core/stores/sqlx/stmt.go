package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
)

const defaultSlowThreshold = time.Millisecond * 500

var (
	slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)
	logSql        = syncx.ForAtomicBool(true)
	logSlowSql    = syncx.ForAtomicBool(true)
)

type (
	// StmtSession interface represents a session that can be used to execute statements.
	StmtSession interface {
		Close() error
		Exec(args ...any) (sql.Result, error)
		ExecCtx(ctx context.Context, args ...any) (sql.Result, error)
		QueryRow(v any, args ...any) error
		QueryRowCtx(ctx context.Context, v any, args ...any) error
		QueryRowPartial(v any, args ...any) error
		QueryRowPartialCtx(ctx context.Context, v any, args ...any) error
		QueryRows(v any, args ...any) error
		QueryRowsCtx(ctx context.Context, v any, args ...any) error
		QueryRowsPartial(v any, args ...any) error
		QueryRowsPartialCtx(ctx context.Context, v any, args ...any) error
	}

	statement struct {
		query  string
		stmt   *sql.Stmt
		brk    breaker.Breaker
		accept breaker.Acceptable
	}

	stmtConn interface {
		Exec(args ...any) (sql.Result, error)
		ExecContext(ctx context.Context, args ...any) (sql.Result, error)
		Query(args ...any) (*sql.Rows, error)
		QueryContext(ctx context.Context, args ...any) (*sql.Rows, error)
	}
)

func (s statement) Close() error {
	return s.stmt.Close()
}

func (s statement) Exec(args ...any) (sql.Result, error) {
	return s.ExecCtx(context.Background(), args...)
}

func (s statement) ExecCtx(ctx context.Context, args ...any) (result sql.Result, err error) {
	ctx, span := startSpan(ctx, "Exec")
	defer func() {
		endSpan(span, err)
	}()

	err = s.brk.DoWithAcceptableCtx(ctx, func() error {
		result, err = execStmt(ctx, s.stmt, s.query, args...)
		return err
	}, func(err error) bool {
		return s.accept(err)
	})
	if errors.Is(err, breaker.ErrServiceUnavailable) {
		metricReqErr.Inc("stmt_exec", "breaker")
	}

	return
}

func (s statement) QueryRow(v any, args ...any) error {
	return s.QueryRowCtx(context.Background(), v, args...)
}

func (s statement) QueryRowCtx(ctx context.Context, v any, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRow")
	defer func() {
		endSpan(span, err)
	}()

	return s.queryRows(ctx, func(v any, scanner rowsScanner) error {
		return unmarshalRow(v, scanner, true)
	}, v, args...)
}

func (s statement) QueryRowPartial(v any, args ...any) error {
	return s.QueryRowPartialCtx(context.Background(), v, args...)
}

func (s statement) QueryRowPartialCtx(ctx context.Context, v any, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowPartial")
	defer func() {
		endSpan(span, err)
	}()

	return s.queryRows(ctx, func(v any, scanner rowsScanner) error {
		return unmarshalRow(v, scanner, false)
	}, v, args...)
}

func (s statement) QueryRows(v any, args ...any) error {
	return s.QueryRowsCtx(context.Background(), v, args...)
}

func (s statement) QueryRowsCtx(ctx context.Context, v any, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRows")
	defer func() {
		endSpan(span, err)
	}()

	return s.queryRows(ctx, func(v any, scanner rowsScanner) error {
		return unmarshalRows(v, scanner, true)
	}, v, args...)
}

func (s statement) QueryRowsPartial(v any, args ...any) error {
	return s.QueryRowsPartialCtx(context.Background(), v, args...)
}

func (s statement) QueryRowsPartialCtx(ctx context.Context, v any, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowsPartial")
	defer func() {
		endSpan(span, err)
	}()

	return s.queryRows(ctx, func(v any, scanner rowsScanner) error {
		return unmarshalRows(v, scanner, false)
	}, v, args...)
}

func (s statement) queryRows(ctx context.Context, scanFn func(any, rowsScanner) error,
	v any, args ...any) error {
	var scanFailed bool
	err := s.brk.DoWithAcceptableCtx(ctx, func() error {
		return queryStmt(ctx, s.stmt, func(rows *sql.Rows) error {
			err := scanFn(v, rows)
			if isScanFailed(err) {
				scanFailed = true
			}
			return err
		}, s.query, args...)
	}, func(err error) bool {
		return scanFailed || s.accept(err)
	})
	if errors.Is(err, breaker.ErrServiceUnavailable) {
		metricReqErr.Inc("stmt_queryRows", "breaker")
	}

	return err
}

// DisableLog disables logging of sql statements, includes info and slow logs.
func DisableLog() {
	logSql.Set(false)
	logSlowSql.Set(false)
}

// DisableStmtLog disables info logging of sql statements, but keeps slow logs.
func DisableStmtLog() {
	logSql.Set(false)
}

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

func exec(ctx context.Context, conn sessionConn, q string, args ...any) (sql.Result, error) {
	guard := newGuard("exec")
	if err := guard.start(q, args...); err != nil {
		return nil, err
	}

	result, err := conn.ExecContext(ctx, q, args...)
	guard.finish(ctx, err)

	return result, err
}

func execStmt(ctx context.Context, conn stmtConn, q string, args ...any) (sql.Result, error) {
	guard := newGuard("execStmt")
	if err := guard.start(q, args...); err != nil {
		return nil, err
	}

	result, err := conn.ExecContext(ctx, args...)
	guard.finish(ctx, err)

	return result, err
}

func query(ctx context.Context, conn sessionConn, scanner func(*sql.Rows) error,
	q string, args ...any) error {
	guard := newGuard("query")
	if err := guard.start(q, args...); err != nil {
		return err
	}

	rows, err := conn.QueryContext(ctx, q, args...)
	guard.finish(ctx, err)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanner(rows)
}

func queryStmt(ctx context.Context, conn stmtConn, scanner func(*sql.Rows) error,
	q string, args ...any) error {
	guard := newGuard("queryStmt")
	if err := guard.start(q, args...); err != nil {
		return err
	}

	rows, err := conn.QueryContext(ctx, args...)
	guard.finish(ctx, err)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanner(rows)
}

type (
	sqlGuard interface {
		start(q string, args ...any) error
		finish(ctx context.Context, err error)
	}

	nilGuard struct{}

	realSqlGuard struct {
		command   string
		stmt      string
		startTime time.Duration
	}
)

func newGuard(command string) sqlGuard {
	if logSql.True() || logSlowSql.True() {
		return &realSqlGuard{
			command: command,
		}
	}

	return nilGuard{}
}

func (n nilGuard) start(_ string, _ ...any) error {
	return nil
}

func (n nilGuard) finish(_ context.Context, _ error) {
}

func (e *realSqlGuard) finish(ctx context.Context, err error) {
	duration := timex.Since(e.startTime)
	if duration > slowThreshold.Load() {
		logx.WithContext(ctx).WithDuration(duration).Slowf("[SQL] %s: slowcall - %s", e.command, e.stmt)
		metricSlowCount.Inc(e.command)
	} else if logSql.True() {
		logx.WithContext(ctx).WithDuration(duration).Infof("sql %s: %s", e.command, e.stmt)
	}

	if err != nil {
		logSqlError(ctx, e.stmt, err)
	}

	metricReqDur.ObserveFloat(float64(duration)/float64(time.Millisecond), e.command)
}

func (e *realSqlGuard) start(q string, args ...any) error {
	stmt, err := format(q, args...)
	if err != nil {
		return err
	}

	e.stmt = stmt
	e.startTime = timex.Now()

	return nil
}
