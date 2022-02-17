package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

func (mt *mockTx) Exec(q string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (mt *mockTx) ExecCtx(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (mt *mockTx) Prepare(query string) (StmtSession, error) {
	return nil, nil
}

func (mt *mockTx) PrepareCtx(ctx context.Context, query string) (StmtSession, error) {
	return nil, nil
}

func (mt *mockTx) QueryRow(v interface{}, q string, args ...interface{}) error {
	return nil
}

func (mt *mockTx) QueryRowCtx(ctx context.Context, v interface{}, query string, args ...interface{}) error {
	return nil
}

func (mt *mockTx) QueryRowPartial(v interface{}, q string, args ...interface{}) error {
	return nil
}

func (mt *mockTx) QueryRowPartialCtx(ctx context.Context, v interface{}, query string, args ...interface{}) error {
	return nil
}

func (mt *mockTx) QueryRows(v interface{}, q string, args ...interface{}) error {
	return nil
}

func (mt *mockTx) QueryRowsCtx(ctx context.Context, v interface{}, query string, args ...interface{}) error {
	return nil
}

func (mt *mockTx) QueryRowsPartial(v interface{}, q string, args ...interface{}) error {
	return nil
}

func (mt *mockTx) QueryRowsPartialCtx(ctx context.Context, v interface{}, query string, args ...interface{}) error {
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
