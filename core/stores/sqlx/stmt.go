package sqlx

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/core/timex"
)

const slowThreshold = time.Millisecond * 500

func exec(conn sessionConn, q string, args ...interface{}) (sql.Result, error) {
	stmt, err := format(q, args...)
	if err != nil {
		return nil, err
	}

	startTime := timex.Now()
	result, err := conn.Exec(q, args...)
	duration := timex.Since(startTime)
	if duration > slowThreshold {
		logx.WithDuration(duration).Slowf("[SQL] exec: slowcall - %s", stmt)
	} else {
		logx.WithDuration(duration).Infof("sql exec: %s", stmt)
	}
	if err != nil {
		logSqlError(stmt, err)
	}

	return result, err
}

func execStmt(conn stmtConn, args ...interface{}) (sql.Result, error) {
	stmt := fmt.Sprint(args...)
	startTime := timex.Now()
	result, err := conn.Exec(args...)
	duration := timex.Since(startTime)
	if duration > slowThreshold {
		logx.WithDuration(duration).Slowf("[SQL] execStmt: slowcall - %s", stmt)
	} else {
		logx.WithDuration(duration).Infof("sql execStmt: %s", stmt)
	}
	if err != nil {
		logSqlError(stmt, err)
	}

	return result, err
}

func query(conn sessionConn, scanner func(*sql.Rows) error, q string, args ...interface{}) error {
	stmt, err := format(q, args...)
	if err != nil {
		return err
	}

	startTime := timex.Now()
	rows, err := conn.Query(q, args...)
	duration := timex.Since(startTime)
	if duration > slowThreshold {
		logx.WithDuration(duration).Slowf("[SQL] query: slowcall - %s", stmt)
	} else {
		logx.WithDuration(duration).Infof("sql query: %s", stmt)
	}
	if err != nil {
		logSqlError(stmt, err)
		return err
	}
	defer rows.Close()

	return scanner(rows)
}

func queryStmt(conn stmtConn, scanner func(*sql.Rows) error, args ...interface{}) error {
	stmt := fmt.Sprint(args...)
	startTime := timex.Now()
	rows, err := conn.Query(args...)
	duration := timex.Since(startTime)
	if duration > slowThreshold {
		logx.WithDuration(duration).Slowf("[SQL] queryStmt: slowcall - %s", stmt)
	} else {
		logx.WithDuration(duration).Infof("sql queryStmt: %s", stmt)
	}
	if err != nil {
		logSqlError(stmt, err)
		return err
	}
	defer rows.Close()

	return scanner(rows)
}
