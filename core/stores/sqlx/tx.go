package sqlx

import (
	"database/sql"
	"fmt"
)

type (
	beginnable func(*sql.DB) (trans, error)

	trans interface {
		Session
		Commit() error
		Rollback() error
	}

	txSession struct {
		*sql.Tx
	}
)

func (t txSession) Exec(q string, args ...interface{}) (sql.Result, error) {
	return exec(t.Tx, q, args...)
}

func (t txSession) Prepare(q string) (StmtSession, error) {
	if stmt, err := t.Tx.Prepare(q); err != nil {
		return nil, err
	} else {
		return statement{
			stmt: stmt,
		}, nil
	}
}

func (t txSession) QueryRow(v interface{}, q string, args ...interface{}) error {
	return query(t.Tx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, q, args...)
}

func (t txSession) QueryRowPartial(v interface{}, q string, args ...interface{}) error {
	return query(t.Tx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, q, args...)
}

func (t txSession) QueryRows(v interface{}, q string, args ...interface{}) error {
	return query(t.Tx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, q, args...)
}

func (t txSession) QueryRowsPartial(v interface{}, q string, args ...interface{}) error {
	return query(t.Tx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

func begin(db *sql.DB) (trans, error) {
	if tx, err := db.Begin(); err != nil {
		return nil, err
	} else {
		return txSession{
			Tx: tx,
		}, nil
	}
}

func transact(db *commonSqlConn, b beginnable, fn func(Session) error) (err error) {
	conn, err := getSqlConn(db.driverName, db.datasource)
	if err != nil {
		logInstanceError(db.datasource, err)
		return err
	}

	return transactOnConn(conn, b, fn)
}

func transactOnConn(conn *sql.DB, b beginnable, fn func(Session) error) (err error) {
	var tx trans
	tx, err = b(conn)
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("recover from %#v, rollback failed: %s", p, e)
			} else {
				err = fmt.Errorf("recoveer from %#v", p)
			}
		} else if err != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("transaction failed: %s, rollback failed: %s", err, e)
			}
		} else {
			err = tx.Commit()
		}
	}()

	return fn(tx)
}
