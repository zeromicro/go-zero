package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/stores/dbtest"
)

const (
	mockCommit   = 1
	mockRollback = 2
)

type mockTx struct {
	status int
}

func (mt *mockTx) Commit() error {
	mt.status |= mockCommit
	return nil
}

func (mt *mockTx) Exec(_ string, _ ...any) (sql.Result, error) {
	return nil, nil
}

func (mt *mockTx) ExecCtx(_ context.Context, _ string, _ ...any) (sql.Result, error) {
	return nil, nil
}

func (mt *mockTx) Prepare(_ string) (StmtSession, error) {
	return nil, nil
}

func (mt *mockTx) PrepareCtx(_ context.Context, _ string) (StmtSession, error) {
	return nil, nil
}

func (mt *mockTx) QueryRow(_ any, _ string, _ ...any) error {
	return nil
}

func (mt *mockTx) QueryRowCtx(_ context.Context, _ any, _ string, _ ...any) error {
	return nil
}

func (mt *mockTx) QueryRowPartial(_ any, _ string, _ ...any) error {
	return nil
}

func (mt *mockTx) QueryRowPartialCtx(_ context.Context, _ any, _ string, _ ...any) error {
	return nil
}

func (mt *mockTx) QueryRows(_ any, _ string, _ ...any) error {
	return nil
}

func (mt *mockTx) QueryRowsCtx(_ context.Context, _ any, _ string, _ ...any) error {
	return nil
}

func (mt *mockTx) QueryRowsPartial(_ any, _ string, _ ...any) error {
	return nil
}

func (mt *mockTx) QueryRowsPartialCtx(_ context.Context, _ any, _ string, _ ...any) error {
	return nil
}

func (mt *mockTx) Rollback() error {
	mt.status |= mockRollback
	return nil
}

func beginMock(mock *mockTx) beginnable {
	return func(*sql.DB) (trans, error) {
		return mock, nil
	}
}

func TestTransactCommit(t *testing.T) {
	mock := &mockTx{}
	err := transactOnConn(context.Background(), nil, beginMock(mock),
		func(context.Context, Session) error {
			return nil
		})
	assert.Equal(t, mockCommit, mock.status)
	assert.Nil(t, err)
}

func TestTransactRollback(t *testing.T) {
	mock := &mockTx{}
	err := transactOnConn(context.Background(), nil, beginMock(mock),
		func(context.Context, Session) error {
			return errors.New("rollback")
		})
	assert.Equal(t, mockRollback, mock.status)
	assert.NotNil(t, err)
}

func TestTxExceptions(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectCommit()
		conn := NewSqlConnFromDB(db)
		assert.NoError(t, conn.Transact(func(session Session) error {
			return nil
		}))
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		conn := &commonSqlConn{
			connProv: func() (*sql.DB, error) {
				return nil, errors.New("foo")
			},
			beginTx: begin,
			onError: func(ctx context.Context, err error) {},
			brk:     breaker.NewBreaker(),
		}
		assert.Error(t, conn.Transact(func(session Session) error {
			return nil
		}))
	})

	runTxTest(t, func(conn SqlConn, mock sqlmock.Sqlmock) {
		_, err := conn.RawDB()
		assert.Equal(t, errNoRawDBFromTx, err)
		assert.Equal(t, errCantNestTx, conn.Transact(nil))
		assert.Equal(t, errCantNestTx, conn.TransactCtx(context.Background(), nil))
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		conn := NewSqlConnFromDB(db)
		assert.Error(t, conn.Transact(func(session Session) error {
			return errors.New("foo")
		}))
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectRollback().WillReturnError(errors.New("foo"))
		conn := NewSqlConnFromDB(db)
		assert.Error(t, conn.Transact(func(session Session) error {
			panic("foo")
		}))
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectRollback()
		conn := NewSqlConnFromDB(db)
		assert.Error(t, conn.Transact(func(session Session) error {
			panic(errors.New("foo"))
		}))
	})
}

func TestTxSession(t *testing.T) {
	runTxTest(t, func(conn SqlConn, mock sqlmock.Sqlmock) {
		mock.ExpectExec("any").WillReturnResult(sqlmock.NewResult(2, 3))
		res, err := conn.Exec("any")
		assert.NoError(t, err)
		last, err := res.LastInsertId()
		assert.NoError(t, err)
		assert.Equal(t, int64(2), last)
		affected, err := res.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(3), affected)

		mock.ExpectExec("any").WillReturnError(errors.New("foo"))
		_, err = conn.Exec("any")
		assert.Equal(t, "foo", err.Error())
	})

	runTxTest(t, func(conn SqlConn, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any")
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)
		assert.NotNil(t, stmt)

		mock.ExpectPrepare("any").WillReturnError(errors.New("foo"))
		_, err = conn.Prepare("any")
		assert.Equal(t, "foo", err.Error())
	})

	runTxTest(t, func(conn SqlConn, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"col"}).AddRow("foo")
		mock.ExpectQuery("any").WillReturnRows(rows)
		var val string
		err := conn.QueryRow(&val, "any")
		assert.NoError(t, err)
		assert.Equal(t, "foo", val)

		mock.ExpectQuery("any").WillReturnError(errors.New("foo"))
		err = conn.QueryRow(&val, "any")
		assert.Equal(t, "foo", err.Error())
	})

	runTxTest(t, func(conn SqlConn, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"col"}).AddRow("foo")
		mock.ExpectQuery("any").WillReturnRows(rows)
		var val string
		err := conn.QueryRowPartial(&val, "any")
		assert.NoError(t, err)
		assert.Equal(t, "foo", val)

		mock.ExpectQuery("any").WillReturnError(errors.New("foo"))
		err = conn.QueryRowPartial(&val, "any")
		assert.Equal(t, "foo", err.Error())
	})

	runTxTest(t, func(conn SqlConn, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"col"}).AddRow("foo").AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(rows)
		var val []string
		err := conn.QueryRows(&val, "any")
		assert.NoError(t, err)
		assert.Equal(t, []string{"foo", "bar"}, val)

		mock.ExpectQuery("any").WillReturnError(errors.New("foo"))
		err = conn.QueryRows(&val, "any")
		assert.Equal(t, "foo", err.Error())
	})

	runTxTest(t, func(conn SqlConn, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"col"}).AddRow("foo").AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(rows)
		var val []string
		err := conn.QueryRowsPartial(&val, "any")
		assert.NoError(t, err)
		assert.Equal(t, []string{"foo", "bar"}, val)

		mock.ExpectQuery("any").WillReturnError(errors.New("foo"))
		err = conn.QueryRowsPartial(&val, "any")
		assert.Equal(t, "foo", err.Error())
	})
}

func TestTxRollback(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("any").WillReturnResult(sqlmock.NewResult(2, 3))
		mock.ExpectQuery("foo").WillReturnError(errors.New("foo"))
		mock.ExpectRollback()

		conn := NewSqlConnFromDB(db)
		err := conn.Transact(func(session Session) error {
			c := NewSqlConnFromSession(session)
			_, err := c.Exec("any")
			assert.NoError(t, err)

			var val string
			return c.QueryRow(&val, "foo")
		})
		assert.Error(t, err)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("any").WillReturnError(errors.New("foo"))
		mock.ExpectRollback()

		conn := NewSqlConnFromDB(db)
		err := conn.Transact(func(session Session) error {
			c := NewSqlConnFromSession(session)
			if _, err := c.Exec("any"); err != nil {
				return err
			}

			var val string
			assert.NoError(t, c.QueryRow(&val, "foo"))
			return nil
		})
		assert.Error(t, err)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("any").WillReturnResult(sqlmock.NewResult(2, 3))
		mock.ExpectQuery("foo").WillReturnRows(sqlmock.NewRows([]string{"col"}).AddRow("bar"))
		mock.ExpectCommit()

		conn := NewSqlConnFromDB(db)
		err := conn.Transact(func(session Session) error {
			c := NewSqlConnFromSession(session)
			_, err := c.Exec("any")
			assert.NoError(t, err)

			var val string
			assert.NoError(t, c.QueryRow(&val, "foo"))
			assert.Equal(t, "bar", val)
			return nil
		})
		assert.NoError(t, err)
	})
}

func runTxTest(t *testing.T, f func(conn SqlConn, mock sqlmock.Sqlmock)) {
	dbtest.RunTxTest(t, func(tx *sql.Tx, mock sqlmock.Sqlmock) {
		sess := NewSessionFromTx(tx)
		conn := NewSqlConnFromSession(sess)
		f(conn, mock)
	})
}
