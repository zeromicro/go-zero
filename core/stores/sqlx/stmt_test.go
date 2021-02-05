package sqlx

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var errMockedPlaceholder = errors.New("placeholder")

func TestStmt_exec(t *testing.T) {
	tests := []struct {
		name         string
		args         []interface{}
		delay        bool
		formatError  bool
		hasError     bool
		err          error
		lastInsertId int64
		rowsAffected int64
	}{
		{
			name:         "normal",
			args:         []interface{}{1},
			lastInsertId: 1,
			rowsAffected: 2,
		},
		{
			name:        "wrong format",
			args:        []interface{}{1, 2},
			formatError: true,
			hasError:    true,
		},
		{
			name:     "exec error",
			args:     []interface{}{1},
			hasError: true,
			err:      errors.New("exec"),
		},
		{
			name:         "slowcall",
			args:         []interface{}{1},
			delay:        true,
			lastInsertId: 1,
			rowsAffected: 2,
		},
	}

	for _, test := range tests {
		test := test
		fns := []func(args ...interface{}) (sql.Result, error){
			func(args ...interface{}) (sql.Result, error) {
				return exec(&mockedSessionConn{
					lastInsertId: test.lastInsertId,
					rowsAffected: test.rowsAffected,
					err:          test.err,
					delay:        test.delay,
				}, "select user from users where id=?", args...)
			},
			func(args ...interface{}) (sql.Result, error) {
				return execStmt(&mockedStmtConn{
					lastInsertId: test.lastInsertId,
					rowsAffected: test.rowsAffected,
					err:          test.err,
					delay:        test.delay,
				}, args...)
			},
		}

		for i, fn := range fns {
			i := i
			fn := fn
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				res, err := fn(test.args...)
				if i == 0 && test.formatError {
					assert.NotNil(t, err)
					return
				}
				if !test.formatError && test.hasError {
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
		name        string
		args        []interface{}
		delay       bool
		formatError bool
		hasError    bool
		err         error
	}{
		{
			name: "normal",
			args: []interface{}{1},
		},
		{
			name:        "wrong format",
			args:        []interface{}{1, 2},
			formatError: true,
			hasError:    true,
		},
		{
			name:     "query error",
			args:     []interface{}{1},
			hasError: true,
			err:      errors.New("exec"),
		},
		{
			name:  "slowcall",
			args:  []interface{}{1},
			delay: true,
		},
	}

	for _, test := range tests {
		test := test
		fns := []func(args ...interface{}) error{
			func(args ...interface{}) error {
				return query(&mockedSessionConn{
					err:   test.err,
					delay: test.delay,
				}, func(rows *sql.Rows) error {
					return nil
				}, "select user from users where id=?", args...)
			},
			func(args ...interface{}) error {
				return queryStmt(&mockedStmtConn{
					err:   test.err,
					delay: test.delay,
				}, func(rows *sql.Rows) error {
					return nil
				}, args...)
			},
		}

		for i, fn := range fns {
			i := i
			fn := fn
			t.Run(test.name, func(t *testing.T) {
				t.Parallel()

				err := fn(test.args...)
				if i == 0 && test.formatError {
					assert.NotNil(t, err)
					return
				}
				if !test.formatError && test.hasError {
					assert.NotNil(t, err)
					return
				}

				assert.Equal(t, errMockedPlaceholder, err)
			})
		}
	}
}

type mockedSessionConn struct {
	lastInsertId int64
	rowsAffected int64
	err          error
	delay        bool
}

func (m *mockedSessionConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.delay {
		time.Sleep(slowThreshold + time.Millisecond)
	}
	return mockedResult{
		lastInsertId: m.lastInsertId,
		rowsAffected: m.rowsAffected,
	}, m.err
}

func (m *mockedSessionConn) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.delay {
		time.Sleep(slowThreshold + time.Millisecond)
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

func (m *mockedStmtConn) Exec(args ...interface{}) (sql.Result, error) {
	if m.delay {
		time.Sleep(slowThreshold + time.Millisecond)
	}
	return mockedResult{
		lastInsertId: m.lastInsertId,
		rowsAffected: m.rowsAffected,
	}, m.err
}

func (m *mockedStmtConn) Query(args ...interface{}) (*sql.Rows, error) {
	if m.delay {
		time.Sleep(slowThreshold + time.Millisecond)
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
