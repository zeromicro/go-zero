package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/stores/dbtest"
)

var errMockedPlaceholder = errors.New("placeholder")

func TestStmt_exec(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		args         []any
		delay        bool
		hasError     bool
		err          error
		lastInsertId int64
		rowsAffected int64
	}{
		{
			name:         "normal",
			query:        "select user from users where id=?",
			args:         []any{1},
			lastInsertId: 1,
			rowsAffected: 2,
		},
		{
			name:     "exec error",
			query:    "select user from users where id=?",
			args:     []any{1},
			hasError: true,
			err:      errors.New("exec"),
		},
		{
			name:     "exec more args error",
			query:    "select user from users where id=? and name=?",
			args:     []any{1},
			hasError: true,
			err:      errors.New("exec"),
		},
		{
			name:         "slowcall",
			query:        "select user from users where id=?",
			args:         []any{1},
			delay:        true,
			lastInsertId: 1,
			rowsAffected: 2,
		},
	}

	for _, test := range tests {
		test := test
		fns := []func(args ...any) (sql.Result, error){
			func(args ...any) (sql.Result, error) {
				return exec(context.Background(), &mockedSessionConn{
					lastInsertId: test.lastInsertId,
					rowsAffected: test.rowsAffected,
					err:          test.err,
					delay:        test.delay,
				}, test.query, args...)
			},
			func(args ...any) (sql.Result, error) {
				return execStmt(context.Background(), &mockedStmtConn{
					lastInsertId: test.lastInsertId,
					rowsAffected: test.rowsAffected,
					err:          test.err,
					delay:        test.delay,
				}, test.query, args...)
			},
		}

		for _, fn := range fns {
			fn := fn
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				res, err := fn(test.args...)
				if test.hasError {
					assert.NotNil(t, err)
					return
				}

				assert.Nil(t, err)
				lastInsertId, err := res.LastInsertId()
				assert.Nil(t, err)
				assert.Equal(t, test.lastInsertId, lastInsertId)
				rowsAffected, err := res.RowsAffected()
				assert.Nil(t, err)
				assert.Equal(t, test.rowsAffected, rowsAffected)
			})
		}
	}
}

func TestStmt_query(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		args     []any
		delay    bool
		hasError bool
		err      error
	}{
		{
			name:  "normal",
			query: "select user from users where id=?",
			args:  []any{1},
		},
		{
			name:     "query error",
			query:    "select user from users where id=?",
			args:     []any{1},
			hasError: true,
			err:      errors.New("exec"),
		},
		{
			name:     "query more args error",
			query:    "select user from users where id=? and name=?",
			args:     []any{1},
			hasError: true,
			err:      errors.New("exec"),
		},
		{
			name:  "slowcall",
			query: "select user from users where id=?",
			args:  []any{1},
			delay: true,
		},
	}

	for _, test := range tests {
		test := test
		fns := []func(args ...any) error{
			func(args ...any) error {
				return query(context.Background(), &mockedSessionConn{
					err:   test.err,
					delay: test.delay,
				}, func(rows *sql.Rows) error {
					return nil
				}, test.query, args...)
			},
			func(args ...any) error {
				return queryStmt(context.Background(), &mockedStmtConn{
					err:   test.err,
					delay: test.delay,
				}, func(rows *sql.Rows) error {
					return nil
				}, test.query, args...)
			},
		}

		for _, fn := range fns {
			fn := fn
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				err := fn(test.args...)
				if test.hasError {
					assert.NotNil(t, err)
					return
				}

				assert.NotNil(t, err)
			})
		}
	}
}

func TestSetSlowThreshold(t *testing.T) {
	assert.Equal(t, defaultSlowThreshold, slowThreshold.Load())
	SetSlowThreshold(time.Second)
	assert.Equal(t, time.Second, slowThreshold.Load())
}

func TestDisableLog(t *testing.T) {
	assert.True(t, logSql.True())
	assert.True(t, logSlowSql.True())
	defer func() {
		logSql.Set(true)
		logSlowSql.Set(true)
	}()

	DisableLog()
	assert.False(t, logSql.True())
	assert.False(t, logSlowSql.True())
}

func TestDisableStmtLog(t *testing.T) {
	assert.True(t, logSql.True())
	assert.True(t, logSlowSql.True())
	defer func() {
		logSql.Set(true)
		logSlowSql.Set(true)
	}()

	DisableStmtLog()
	assert.False(t, logSql.True())
	assert.True(t, logSlowSql.True())
}

func TestNilGuard(t *testing.T) {
	assert.True(t, logSql.True())
	assert.True(t, logSlowSql.True())
	defer func() {
		logSql.Set(true)
		logSlowSql.Set(true)
	}()

	DisableLog()
	guard := newGuard("any")
	assert.Nil(t, guard.start("foo", "bar"))
	guard.finish(context.Background(), nil)
	assert.Equal(t, nilGuard{}, guard)
}

func TestStmtBreaker(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any")

		conn := NewSqlConnFromDB(db)
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)

		var val struct {
			Foo int
			Bar string
		}
		for i := 0; i < 1000; i++ {
			row := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
			mock.ExpectQuery("any").WillReturnRows(row)
			err := stmt.QueryRow(&val)
			assert.Error(t, err)
			assert.NotErrorIs(t, err, breaker.ErrServiceUnavailable)
		}
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any")
		conn := NewSqlConnFromDB(db)
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)

		for i := 0; i < 1000; i++ {
			assert.Error(t, conn.Transact(func(session Session) error {
				return nil
			}))
		}

		var breakerTriggered bool
		for i := 0; i < 1000; i++ {
			_, err = stmt.Exec("any")
			if errors.Is(err, breaker.ErrServiceUnavailable) {
				breakerTriggered = true
				break
			}
		}
		assert.True(t, breakerTriggered)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any")
		conn := NewSqlConnFromDB(db)
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)

		for i := 0; i < 1000; i++ {
			assert.Error(t, conn.Transact(func(session Session) error {
				return nil
			}))
		}

		var breakerTriggered bool
		for i := 0; i < 1000; i++ {
			err = stmt.QueryRows(&struct{}{}, "any")
			if errors.Is(err, breaker.ErrServiceUnavailable) {
				breakerTriggered = true
				break
			}
		}
		assert.True(t, breakerTriggered)
	})
}

func TestQueryRowsScanTimeout(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"foo"})
		for i := 0; i < 10000; i++ {
			rows = rows.AddRow("bar" + strconv.Itoa(i))
		}
		mock.ExpectQuery("any").WillReturnRows(rows)
		var val []struct {
			Foo string
		}
		conn := NewSqlConnFromDB(db)
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*2)
		err := conn.QueryRowsCtx(ctx, &val, "any")
		assert.ErrorIs(t, err, context.DeadlineExceeded)
		cancel()
	})
}

type mockedSessionConn struct {
	lastInsertId int64
	rowsAffected int64
	err          error
	delay        bool
}

func (m *mockedSessionConn) Exec(query string, args ...any) (sql.Result, error) {
	return m.ExecContext(context.Background(), query, args...)
}

func (m *mockedSessionConn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if m.delay {
		time.Sleep(defaultSlowThreshold + time.Millisecond)
	}
	return mockedResult{
		lastInsertId: m.lastInsertId,
		rowsAffected: m.rowsAffected,
	}, m.err
}

func (m *mockedSessionConn) Query(query string, args ...any) (*sql.Rows, error) {
	return m.QueryContext(context.Background(), query, args...)
}

func (m *mockedSessionConn) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if m.delay {
		time.Sleep(defaultSlowThreshold + time.Millisecond)
	}

	err := errMockedPlaceholder
	if m.err != nil {
		err = m.err
	}
	return new(sql.Rows), err
}

type mockedStmtConn struct {
	lastInsertId int64
	rowsAffected int64
	err          error
	delay        bool
}

func (m *mockedStmtConn) Exec(args ...any) (sql.Result, error) {
	return m.ExecContext(context.Background(), args...)
}

func (m *mockedStmtConn) ExecContext(_ context.Context, _ ...any) (sql.Result, error) {
	if m.delay {
		time.Sleep(defaultSlowThreshold + time.Millisecond)
	}
	return mockedResult{
		lastInsertId: m.lastInsertId,
		rowsAffected: m.rowsAffected,
	}, m.err
}

func (m *mockedStmtConn) Query(args ...any) (*sql.Rows, error) {
	return m.QueryContext(context.Background(), args...)
}

func (m *mockedStmtConn) QueryContext(_ context.Context, _ ...any) (*sql.Rows, error) {
	if m.delay {
		time.Sleep(defaultSlowThreshold + time.Millisecond)
	}

	err := errMockedPlaceholder
	if m.err != nil {
		err = m.err
	}
	return new(sql.Rows), err
}

type mockedResult struct {
	lastInsertId int64
	rowsAffected int64
}

func (m mockedResult) LastInsertId() (int64, error) {
	return m.lastInsertId, nil
}

func (m mockedResult) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}
