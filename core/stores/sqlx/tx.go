package sqlx

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/breaker"
)

type (
	beginnable func(*sql.DB, ...TxOption) (trans, error)

	trans interface {
		Session
		Commit() error
		Rollback() error
	}

	txConn struct {
		Session
	}

	txSession struct {
		*sql.Tx
		unmarshalRowHandler  UnmarshalRowHandler
		unmarshalRowsHandler UnmarshalRowsHandler
	}

	TxOption func(*txSession)
)

func (s txConn) RawDB() (*sql.DB, error) {
	return nil, errNoRawDBFromTx
}

func (s txConn) Transact(_ func(Session) error) error {
	return errCantNestTx
}

func (s txConn) TransactCtx(_ context.Context, _ func(context.Context, Session) error) error {
	return errCantNestTx
}

// WithTxUnmarshalRowHandler sets the UnmarshalRowHandler for the txSession.
// It's used to customize the unmarshal behavior for QueryRow and QueryRowPartial.
func WithTxUnmarshalRowHandler(handler UnmarshalRowHandler) TxOption {
	return func(ts *txSession) {
		ts.unmarshalRowHandler = handler
	}
}

// WithTxUnmarshalRowsHandler sets the UnmarshalRowsHandler for the txSession.
// It's used to customize the unmarshal behavior for QueryRows and QueryRowsPartial.
func WithTxUnmarshalRowsHandler(handler UnmarshalRowsHandler) TxOption {
	return func(ts *txSession) {
		ts.unmarshalRowsHandler = handler
	}
}

// NewSessionFromTx returns a Session with the given sql.Tx.
// Use it with caution, it's provided for other ORM to interact with.
func NewSessionFromTx(tx *sql.Tx, opts ...TxOption) Session {
	ts := txSession{Tx: tx}
	for _, opt := range opts {
		opt(&ts)
	}
	return ts
}

func (t txSession) Exec(q string, args ...any) (sql.Result, error) {
	return t.ExecCtx(context.Background(), q, args...)
}

func (t txSession) ExecCtx(ctx context.Context, q string, args ...any) (result sql.Result, err error) {
	ctx, span := startSpan(ctx, "Exec")
	defer func() {
		endSpan(span, err)
	}()

	result, err = exec(ctx, t.Tx, q, args...)

	return
}

func (t txSession) Prepare(q string) (StmtSession, error) {
	return t.PrepareCtx(context.Background(), q)
}

func (t txSession) PrepareCtx(ctx context.Context, q string) (stmtSession StmtSession, err error) {
	ctx, span := startSpan(ctx, "Prepare")
	defer func() {
		endSpan(span, err)
	}()

	stmt, err := t.Tx.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}

	return statement{
		query:                q,
		stmt:                 stmt,
		brk:                  breaker.NopBreaker(),
		unmarshalRowHandler:  t.unmarshalRowHandler,
		unmarshalRowsHandler: t.unmarshalRowsHandler,
	}, nil
}

func (t txSession) QueryRow(v any, q string, args ...any) error {
	return t.QueryRowCtx(context.Background(), v, q, args...)
}

func (t txSession) QueryRowCtx(ctx context.Context, v any, q string, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRow")
	defer func() {
		endSpan(span, err)
	}()

	return query(ctx, t.Tx, func(rows *sql.Rows) error {
		if t.unmarshalRowHandler != nil {
			return t.unmarshalRowHandler(v, rows, true)
		}
		return unmarshalRow(v, rows, true)
	}, q, args...)
}

func (t txSession) QueryRowPartial(v any, q string, args ...any) error {
	return t.QueryRowPartialCtx(context.Background(), v, q, args...)
}

func (t txSession) QueryRowPartialCtx(ctx context.Context, v any, q string,
	args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowPartial")
	defer func() {
		endSpan(span, err)
	}()

	return query(ctx, t.Tx, func(rows *sql.Rows) error {
		if t.unmarshalRowHandler != nil {
			return t.unmarshalRowHandler(v, rows, false)
		}
		return unmarshalRow(v, rows, false)
	}, q, args...)
}

func (t txSession) QueryRows(v any, q string, args ...any) error {
	return t.QueryRowsCtx(context.Background(), v, q, args...)
}

func (t txSession) QueryRowsCtx(ctx context.Context, v any, q string, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRows")
	defer func() {
		endSpan(span, err)
	}()

	return query(ctx, t.Tx, func(rows *sql.Rows) error {
		if t.unmarshalRowsHandler != nil {
			return t.unmarshalRowsHandler(v, rows, true)
		}
		return unmarshalRows(v, rows, true)
	}, q, args...)
}

func (t txSession) QueryRowsPartial(v any, q string, args ...any) error {
	return t.QueryRowsPartialCtx(context.Background(), v, q, args...)
}

func (t txSession) QueryRowsPartialCtx(ctx context.Context, v any, q string,
	args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowsPartial")
	defer func() {
		endSpan(span, err)
	}()

	return query(ctx, t.Tx, func(rows *sql.Rows) error {
		if t.unmarshalRowsHandler != nil {
			return t.unmarshalRowsHandler(v, rows, false)
		}
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

func begin(db *sql.DB, opts ...TxOption) (trans, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	ts := txSession{
		Tx: tx,
	}

	for _, opt := range opts {
		opt(&ts)
	}

	return ts, nil
}

func transact(ctx context.Context, db *commonSqlConn, b beginnable,
	fn func(context.Context, Session) error) (err error) {
	conn, err := db.connProv()
	if err != nil {
		db.onError(ctx, err)
		return err
	}

	return transactOnConn(ctx, conn, func(d *sql.DB, _ ...TxOption) (trans, error) {
		return b(d,
			WithTxUnmarshalRowHandler(db.unmarshalRowHandler),
			WithTxUnmarshalRowsHandler(db.unmarshalRowsHandler),
		)
	}, fn)
}

func transactOnConn(ctx context.Context, conn *sql.DB, b beginnable,
	fn func(context.Context, Session) error) (err error) {
	var tx trans
	tx, err = b(conn)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("recover from %#v, rollback failed: %w", p, e)
			} else {
				err = fmt.Errorf("recover from %#v", p)
			}
		} else if err != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("transaction failed: %s, rollback failed: %w", err, e)
			}
		} else {
			err = tx.Commit()
		}
	}()

	return fn(ctx, tx)
}
