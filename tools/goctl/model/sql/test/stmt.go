// copy from core/stores/sqlx/stmt.go

package mocksql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/timex"
)

const slowThreshold = time.Millisecond * 500

func exec(db *sql.DB, q string, args ...any) (sql.Result, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	stmt, err := format(q, args...)
	if err != nil {
		return nil, err
	}

	startTime := timex.Now()
	result, err := tx.Exec(q, args...)
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

func execStmt(conn *sql.Stmt, args ...any) (sql.Result, error) {
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

func query(db *sql.DB, scanner func(*sql.Rows) error, q string, args ...any) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	stmt, err := format(q, args...)
	if err != nil {
		return err
	}

	startTime := timex.Now()
	rows, err := tx.Query(q, args...)
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

func queryStmt(conn *sql.Stmt, scanner func(*sql.Rows) error, args ...any) error {
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
