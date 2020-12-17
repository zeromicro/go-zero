// copy from core/stores/sqlx/sqlconn.go
package mocksql

import (
	"database/sql"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type (
	MockConn struct {
		db *sql.DB
	}
	statement struct {
		stmt *sql.Stmt
	}
)

func NewMockConn(db *sql.DB) *MockConn {
	return &MockConn{db: db}
}

func (conn *MockConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return exec(conn.db, query, args...)
}

func (conn *MockConn) Prepare(query string) (sqlx.StmtSession, error) {
	st, err := conn.db.Prepare(query)
	return statement{stmt: st}, err
}

func (conn *MockConn) QueryRow(v interface{}, q string, args ...interface{}) error {
	return query(conn.db, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, q, args...)
}

func (conn *MockConn) QueryRowPartial(v interface{}, q string, args ...interface{}) error {
	return query(conn.db, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, q, args...)
}

func (conn *MockConn) QueryRows(v interface{}, q string, args ...interface{}) error {
	return query(conn.db, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, q, args...)
}

func (conn *MockConn) QueryRowsPartial(v interface{}, q string, args ...interface{}) error {
	return query(conn.db, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

func (conn *MockConn) Transact(func(session sqlx.Session) error) error {
	return nil
}

func (s statement) Close() error {
	return s.stmt.Close()
}

func (s statement) Exec(args ...interface{}) (sql.Result, error) {
	return execStmt(s.stmt, args...)
}

func (s statement) QueryRow(v interface{}, args ...interface{}) error {
	return queryStmt(s.stmt, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, args...)
}

func (s statement) QueryRowPartial(v interface{}, args ...interface{}) error {
	return queryStmt(s.stmt, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, args...)
}

func (s statement) QueryRows(v interface{}, args ...interface{}) error {
	return queryStmt(s.stmt, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, args...)
}

func (s statement) QueryRowsPartial(v interface{}, args ...interface{}) error {
	return queryStmt(s.stmt, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, args...)
}
