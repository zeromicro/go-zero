package sqlx

import (
	"context"
	"database/sql"
	"fmt"
)

type (
	beginnable func(ctx context.Context, logOpt logOption, db *sql.DB) (trans, error)

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
		logOpt logOption
	}
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

// NewSessionFromTx returns a Session with the given sql.Tx.
// Use it with caution, it's provided for other ORM to interact with.
func NewSessionFromTx(tx *sql.Tx) Session {
	return txSession{Tx: tx}
}

func (t txSession) Exec(q string, args ...any) (sql.Result, error) {
	return t.ExecCtx(context.Background(), q, args...)
}

func (t txSession) ExecCtx(ctx context.Context, q string, args ...any) (result sql.Result, err error) {
	ctx, span := startSpan(ctx, "Exec")
	defer func() {
		endSpan(span, err)
	}()

	result, err = exec(ctx, t.logOpt, t.Tx, q, args...)

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

	guard := newGuard(t.logOpt, "prepare")
	_ = guard.start(q)

	stmt, err := t.Tx.PrepareContext(ctx, q)
	guard.finish(ctx, err)
	if err != nil {
		return nil, err
	}

	return statement{
		query:     q,
		stmt:      stmt,
		logOption: t.logOpt,
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

	return query(ctx, t.logOpt, t.Tx, func(rows *sql.Rows) error {
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

	return query(ctx, t.logOpt, t.Tx, func(rows *sql.Rows) error {
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

	return query(ctx, t.logOpt, t.Tx, func(rows *sql.Rows) error {
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

	return query(ctx, t.logOpt, t.Tx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

func commit(ctx context.Context, logOpt logOption, tx trans) (err error) {
	ctx, span := startSpan(ctx, "Commit")
	defer func() {
		endSpan(span, err)
	}()

	guard := newGuard(logOpt, "transact")
	_ = guard.start("COMMIT")

	err = tx.Commit()
	guard.finish(ctx, err)

	return
}

func rollback(ctx context.Context, logOpt logOption, tx trans) (err error) {
	ctx, span := startSpan(ctx, "Rollback")
	defer func() {
		endSpan(span, err)
	}()

	guard := newGuard(logOpt, "transact")
	_ = guard.start("ROLLBACK")

	err = tx.Rollback()
	guard.finish(ctx, err)

	return
}

func begin(ctx context.Context, logOpt logOption, db *sql.DB) (t trans, err error) {
	ctx, span := startSpan(ctx, "Begin")
	defer func() {
		endSpan(span, err)
	}()

	guard := newGuard(logOpt, "transact")
	_ = guard.start("BEGIN")

	var tx *sql.Tx
	tx, err = db.Begin()

	guard.finish(ctx, err)

	if err != nil {
		return
	}

	return txSession{
		Tx:     tx,
		logOpt: logOpt,
	}, nil
}

func transact(ctx context.Context, logOpt logOption, db *commonSqlConn, b beginnable,
	fn func(context.Context, Session) error) (err error) {
	conn, err := db.connProv()
	if err != nil {
		db.onError(ctx, err)
		return err
	}

	return transactOnConn(ctx, logOpt, conn, b, fn)
}

func transactOnConn(ctx context.Context, logOpt logOption, conn *sql.DB, b beginnable,
	fn func(context.Context, Session) error) (err error) {
	var tx trans
	tx, err = b(ctx, logOpt, conn)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			if e := rollback(ctx, logOpt, tx); e != nil {
				err = fmt.Errorf("recover from %#v, rollback failed: %w", p, e)
			} else {
				err = fmt.Errorf("recover from %#v", p)
			}
		} else if err != nil {
			if e := rollback(ctx, logOpt, tx); e != nil {
				err = fmt.Errorf("transaction failed: %s, rollback failed: %w", err, e)
			}
		} else {
			err = commit(ctx, logOpt, tx)
		}
	}()

	return fn(ctx, tx)
}
