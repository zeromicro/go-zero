package sqlx

import (
	"context"
	"database/sql"
	"reflect"
	"strings"
	"time"

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

// FormatSql expand slice or array arguments
func FormatSql(query string, args ...interface{}) (string, []interface{}) {
	var n0, n1, j int
	for _, arg := range args {
		switch av := reflect.ValueOf(arg); av.Kind() {
		case reflect.Slice, reflect.Array:
			if j = av.Len(); j > 0 {
				n0 += j
				n1++
			} else {
				n0++
			}
		default:
			n0++
		}
	}

	if n1 == 0 || n0 == 0 {
		return query, args
	}

	type argMap struct{ index, len int }
	var (
		resp    = make([]interface{}, 0, n0)
		argMaps = make([]argMap, 0, n1)
	)
	n0 = 0
	for i, arg := range args {
		switch av := reflect.ValueOf(arg); av.Kind() {
		case reflect.Slice, reflect.Array:
			if n1 = av.Len() - 1; n1 >= 0 {
				for j = 0; j <= n1; j++ {
					resp = append(resp, av.Index(j).Interface())
				}
				argMaps = append(argMaps, argMap{index: i, len: n1})
				n0 += n1
			} else {
				resp = append(resp, "NULL")
			}
		default:
			resp = append(resp, arg)
		}
	}

	var b strings.Builder
	b.Grow(len(query) + 2*n0)
	n0, n1 = 0, 0
	for _, v := range query {
		b.WriteRune(v)
		if v == '?' && n0 < len(argMaps) {
			if argMaps[n0].index == n1 {
				for j = argMaps[n0].len; j > 0; j-- {
					b.WriteString(",?")
				}
				n0++
			}
			n1++
		}
	}
	return b.String(), resp
}

func exec(ctx context.Context, conn sessionConn, q string, args ...interface{}) (sql.Result, error) {
	q, args = FormatSql(q, args...)
	guard := newGuard("exec")
	if err := guard.start(q, args...); err != nil {
		return nil, err
	}

	result, err := conn.ExecContext(ctx, q, args...)
	guard.finish(ctx, err)

	return result, err
}

func execStmt(ctx context.Context, conn stmtConn, q string, args ...interface{}) (sql.Result, error) {
	q, args = FormatSql(q, args...)
	guard := newGuard("execStmt")
	if err := guard.start(q, args...); err != nil {
		return nil, err
	}

	result, err := conn.ExecContext(ctx, args...)
	guard.finish(ctx, err)

	return result, err
}

func query(ctx context.Context, conn sessionConn, scanner func(*sql.Rows) error,
	q string, args ...interface{}) error {
	q, args = FormatSql(q, args...)
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
	q string, args ...interface{}) error {
	q, args = FormatSql(q, args...)
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
		start(q string, args ...interface{}) error
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

func (n nilGuard) start(_ string, _ ...interface{}) error {
	return nil
}

func (n nilGuard) finish(_ context.Context, _ error) {
}

func (e *realSqlGuard) finish(ctx context.Context, err error) {
	duration := timex.Since(e.startTime)
	if duration > slowThreshold.Load() {
		logx.WithContext(ctx).WithDuration(duration).Slowf("[SQL] %s: slowcall - %s", e.command, e.stmt)
	} else if logSql.True() {
		logx.WithContext(ctx).WithDuration(duration).Infof("sql %s: %s", e.command, e.stmt)
	}

	if err != nil {
		logSqlError(ctx, e.stmt, err)
	}

	metricReqDur.Observe(int64(duration/time.Millisecond), e.command)
}

func (e *realSqlGuard) start(q string, args ...interface{}) error {
	stmt, err := format(q, args...)
	if err != nil {
		return err
	}

	e.stmt = stmt
	e.startTime = timex.Now()

	return nil
}
