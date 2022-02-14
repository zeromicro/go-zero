package sqlx

import (
	"context"
	"database/sql"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
)

const defaultSlowThreshold = time.Millisecond * 500

var slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

func exec(ctx context.Context, conn sessionConn, q string, args ...interface{}) (sql.Result, error) {
	stmt, err := format(q, args...)
	if err != nil {
		return nil, err
	}

	startTime := timex.Now()
	result, err := conn.ExecContext(ctx, q, args...)
	duration := timex.Since(startTime)
	if duration > slowThreshold.Load() {
		logx.WithContext(ctx).WithDuration(duration).Slowf("[SQL] exec: slowcall - %s", stmt)
	} else {
		logx.WithContext(ctx).WithDuration(duration).Infof("sql exec: %s", stmt)
	}
	if err != nil {
		logSqlError(ctx, stmt, err)
	}

	return result, err
}

func execStmt(ctx context.Context, conn stmtConn, q string, args ...interface{}) (sql.Result, error) {
	stmt, err := format(q, args...)
	if err != nil {
		return nil, err
	}

	startTime := timex.Now()
	result, err := conn.ExecContext(ctx, args...)
	duration := timex.Since(startTime)
	if duration > slowThreshold.Load() {
		logx.WithContext(ctx).WithDuration(duration).Slowf("[SQL] execStmt: slowcall - %s", stmt)
	} else {
		logx.WithContext(ctx).WithDuration(duration).Infof("sql execStmt: %s", stmt)
	}
	if err != nil {
		logSqlError(ctx, stmt, err)
	}

	return result, err
}

func query(ctx context.Context, conn sessionConn, scanner func(*sql.Rows) error,
	q string, args ...interface{}) error {
	stmt, err := format(q, args...)
	if err != nil {
		return err
	}

	startTime := timex.Now()
	rows, err := conn.QueryContext(ctx, q, args...)
	duration := timex.Since(startTime)
	if duration > slowThreshold.Load() {
		logx.WithContext(ctx).WithDuration(duration).Slowf("[SQL] query: slowcall - %s", stmt)
	} else {
		logx.WithContext(ctx).WithDuration(duration).Infof("sql query: %s", stmt)
	}
	if err != nil {
		logSqlError(ctx, stmt, err)
		return err
	}
	defer rows.Close()

	return scanner(rows)
}

func queryStmt(ctx context.Context, conn stmtConn, scanner func(*sql.Rows) error,
	q string, args ...interface{}) error {
	stmt, err := format(q, args...)
	if err != nil {
		return err
	}

	startTime := timex.Now()
	rows, err := conn.QueryContext(ctx, args...)
	duration := timex.Since(startTime)
	if duration > slowThreshold.Load() {
		logx.WithContext(ctx).WithDuration(duration).Slowf("[SQL] queryStmt: slowcall - %s", stmt)
	} else {
		logx.WithContext(ctx).WithDuration(duration).Infof("sql queryStmt: %s", stmt)
	}
	if err != nil {
		logSqlError(ctx, stmt, err)
		return err
	}
	defer rows.Close()

	return scanner(rows)
}
