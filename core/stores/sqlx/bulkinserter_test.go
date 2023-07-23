package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/internal/dbtest"
)

type mockedConn struct {
	query   string
	args    []any
	execErr error
}

func (c *mockedConn) ExecCtx(_ context.Context, query string, args ...any) (sql.Result, error) {
	c.query = query
	c.args = args
	return nil, c.execErr
}

func (c *mockedConn) PrepareCtx(ctx context.Context, query string) (StmtSession, error) {
	panic("implement me")
}

func (c *mockedConn) QueryRowCtx(ctx context.Context, v any, query string, args ...any) error {
	panic("implement me")
}

func (c *mockedConn) QueryRowPartialCtx(ctx context.Context, v any, query string, args ...any) error {
	panic("implement me")
}

func (c *mockedConn) QueryRowsCtx(ctx context.Context, v any, query string, args ...any) error {
	panic("implement me")
}

func (c *mockedConn) QueryRowsPartialCtx(ctx context.Context, v any, query string, args ...any) error {
	panic("implement me")
}

func (c *mockedConn) TransactCtx(ctx context.Context, fn func(context.Context, Session) error) error {
	panic("should not called")
}

func (c *mockedConn) Exec(query string, args ...any) (sql.Result, error) {
	return c.ExecCtx(context.Background(), query, args...)
}

func (c *mockedConn) Prepare(query string) (StmtSession, error) {
	panic("should not called")
}

func (c *mockedConn) QueryRow(v any, query string, args ...any) error {
	panic("should not called")
}

func (c *mockedConn) QueryRowPartial(v any, query string, args ...any) error {
	panic("should not called")
}

func (c *mockedConn) QueryRows(v any, query string, args ...any) error {
	panic("should not called")
}

func (c *mockedConn) QueryRowsPartial(v any, query string, args ...any) error {
	panic("should not called")
}

func (c *mockedConn) RawDB() (*sql.DB, error) {
	panic("should not called")
}

func (c *mockedConn) Transact(func(session Session) error) error {
	panic("should not called")
}

func TestBulkInserter(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var conn mockedConn
		inserter, err := NewBulkInserter(&conn, `INSERT INTO classroom_dau(classroom, user, count) VALUES(?, ?, ?)`)
		assert.Nil(t, err)
		for i := 0; i < 5; i++ {
			assert.Nil(t, inserter.Insert("class_"+strconv.Itoa(i), "user_"+strconv.Itoa(i), i))
		}
		inserter.Flush()
		assert.Equal(t, `INSERT INTO classroom_dau(classroom, user, count) VALUES `+
			`('class_0', 'user_0', 0), ('class_1', 'user_1', 1), ('class_2', 'user_2', 2), `+
			`('class_3', 'user_3', 3), ('class_4', 'user_4', 4)`,
			conn.query)
		assert.Nil(t, conn.args)
	})
}

func TestBulkInserterSuffix(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var conn mockedConn
		inserter, err := NewBulkInserter(&conn, `INSERT INTO classroom_dau(classroom, user, count) VALUES`+
			`(?, ?, ?) ON DUPLICATE KEY UPDATE is_overtime=VALUES(is_overtime)`)
		assert.Nil(t, err)
		assert.Nil(t, inserter.UpdateStmt(`INSERT INTO classroom_dau(classroom, user, count) VALUES`+
			`(?, ?, ?) ON DUPLICATE KEY UPDATE is_overtime=VALUES(is_overtime)`))
		for i := 0; i < 5; i++ {
			assert.Nil(t, inserter.Insert("class_"+strconv.Itoa(i), "user_"+strconv.Itoa(i), i))
		}
		inserter.SetResultHandler(func(result sql.Result, err error) {})
		inserter.Flush()
		assert.Equal(t, `INSERT INTO classroom_dau(classroom, user, count) VALUES `+
			`('class_0', 'user_0', 0), ('class_1', 'user_1', 1), ('class_2', 'user_2', 2), `+
			`('class_3', 'user_3', 3), ('class_4', 'user_4', 4) ON DUPLICATE KEY UPDATE is_overtime=VALUES(is_overtime)`,
			conn.query)
		assert.Nil(t, conn.args)
	})
}

func TestBulkInserterBadStatement(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		var conn mockedConn
		_, err := NewBulkInserter(&conn, "foo")
		assert.NotNil(t, err)
	})
}

func TestBulkInserter_Update(t *testing.T) {
	conn := mockedConn{
		execErr: errors.New("foo"),
	}
	_, err := NewBulkInserter(&conn, `INSERT INTO classroom_dau(classroom, user, count) VALUES()`)
	assert.NotNil(t, err)
	_, err = NewBulkInserter(&conn, `INSERT INTO classroom_dau(classroom, user, count) VALUES(?)`)
	assert.NotNil(t, err)
	inserter, err := NewBulkInserter(&conn, `INSERT INTO classroom_dau(classroom, user, count) VALUES(?, ?, ?)`)
	assert.Nil(t, err)
	inserter.inserter.Execute([]string{"bar"})
	inserter.SetResultHandler(func(result sql.Result, err error) {
	})
	inserter.UpdateOrDelete(func() {})
	inserter.inserter.Execute([]string(nil))
	assert.NotNil(t, inserter.UpdateStmt("foo"))
	assert.NotNil(t, inserter.Insert("foo", "bar"))
}
