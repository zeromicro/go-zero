// copy from core/stores/sqlx/sqlconn.go

package mocksql

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	// MockConn defines a mock connection instance for mysql
	MockConn struct {
		db *sql.DB
	}

	statement struct {
		stmt *sql.Stmt
	}
)

// NewMockConn creates an instance for MockConn
func NewMockConn(db *sql.DB) *MockConn {
	return &MockConn{db: db}
}

// Exec executes sql and returns the result
func (conn *MockConn) Exec(query string, args ...any) (sql.Result, error) {
	return exec(conn.db, query, args...)
}

// ExecCtx executes sql and returns the result
func (conn *MockConn) ExecCtx(_ context.Context, query string, args ...any) (sql.Result, error) {
	return exec(conn.db, query, args...)
}

// Prepare executes sql by sql.DB
func (conn *MockConn) Prepare(query string) (sqlx.StmtSession, error) {
	st, err := conn.db.Prepare(query)
	return statement{stmt: st}, err
}

// PrepareCtx executes sql by sql.DB
func (conn *MockConn) PrepareCtx(_ context.Context, query string) (sqlx.StmtSession, error) {
	return conn.Prepare(query)
}

// QueryRow executes sql and returns a query row
func (conn *MockConn) QueryRow(v any, q string, args ...any) error {
	return query(conn.db, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, q, args...)
}

// QueryRowCtx executes sql and returns a query row
func (conn *MockConn) QueryRowCtx(_ context.Context, v any, query string, args ...any) error {
	return conn.QueryRow(v, query, args...)
}

// QueryRowPartial executes sql and returns a partial query row
func (conn *MockConn) QueryRowPartial(v any, q string, args ...any) error {
	return query(conn.db, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, q, args...)
}

// QueryRowPartialCtx executes sql and returns a partial query row
func (conn *MockConn) QueryRowPartialCtx(_ context.Context, v any, query string, args ...any) error {
	return conn.QueryRowPartial(v, query, args...)
}

// QueryRows executes sql and returns  query rows
func (conn *MockConn) QueryRows(v any, q string, args ...any) error {
	return query(conn.db, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, q, args...)
}

// QueryRowsCtx executes sql and returns  query rows
func (conn *MockConn) QueryRowsCtx(_ context.Context, v any, query string, args ...any) error {
	return conn.QueryRows(v, query, args...)
}

// QueryRowsPartial executes sql and returns partial query rows
func (conn *MockConn) QueryRowsPartial(v any, q string, args ...any) error {
	return query(conn.db, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

// QueryRowsPartialCtx executes sql and returns partial query rows
func (conn *MockConn) QueryRowsPartialCtx(_ context.Context, v any, query string, args ...any) error {
	return conn.QueryRowsPartial(v, query, args...)
}

// RawDB returns the underlying sql.DB.
func (conn *MockConn) RawDB() (*sql.DB, error) {
	return conn.db, nil
}

// Transact is the implemention of sqlx.SqlConn, nothing to do
func (conn *MockConn) Transact(func(session sqlx.Session) error) error {
	return nil
}

// TransactCtx is the implemention of sqlx.SqlConn, nothing to do
func (conn *MockConn) TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return nil
}

func (s statement) Close() error {
	return s.stmt.Close()
}

func (s statement) Exec(args ...any) (sql.Result, error) {
	return execStmt(s.stmt, args...)
}

func (s statement) ExecCtx(_ context.Context, args ...any) (sql.Result, error) {
	return s.Exec(args...)
}

func (s statement) QueryRow(v any, args ...any) error {
	return queryStmt(s.stmt, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, args...)
}

func (s statement) QueryRowCtx(_ context.Context, v any, args ...any) error {
	return s.QueryRow(v, args...)
}

func (s statement) QueryRowPartial(v any, args ...any) error {
	return queryStmt(s.stmt, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, args...)
}

func (s statement) QueryRowPartialCtx(_ context.Context, v any, args ...any) error {
	return s.QueryRowPartial(v, args...)
}

func (s statement) QueryRows(v any, args ...any) error {
	return queryStmt(s.stmt, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, args...)
}

func (s statement) QueryRowsCtx(_ context.Context, v any, args ...any) error {
	return s.QueryRows(v, args...)
}

func (s statement) QueryRowsPartial(v any, args ...any) error {
	return queryStmt(s.stmt, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, args...)
}

func (s statement) QueryRowsPartialCtx(_ context.Context, v any, args ...any) error {
	return s.QueryRowsPartial(v, args...)
}
