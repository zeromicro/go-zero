package sqlx

import (
	"database/sql"

	"github.com/tal-tech/go-zero/core/breaker"
)

// ErrNotFound is an alias of sql.ErrNoRows
var ErrNotFound = sql.ErrNoRows

type (
	// Session stands for raw connections or transaction sessions
	Session interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
		Prepare(query string) (StmtSession, error)
		QueryRow(v interface{}, query string, args ...interface{}) error
		QueryRowPartial(v interface{}, query string, args ...interface{}) error
		QueryRows(v interface{}, query string, args ...interface{}) error
		QueryRowsPartial(v interface{}, query string, args ...interface{}) error
	}

	// SqlConn only stands for raw connections, so Transact method can be called.
	SqlConn interface {
		Session
		Transact(func(session Session) error) error
		QueryRowsRaw(q string, args ...interface{}) (*[]map[string]interface{}, error)
	}

	// SqlOption defines the method to customize a sql connection.
	SqlOption func(*commonSqlConn)

	// StmtSession interface represents a session that can be used to execute statements.
	StmtSession interface {
		Close() error
		Exec(args ...interface{}) (sql.Result, error)
		QueryRow(v interface{}, args ...interface{}) error
		QueryRowPartial(v interface{}, args ...interface{}) error
		QueryRows(v interface{}, args ...interface{}) error
		QueryRowsPartial(v interface{}, args ...interface{}) error
	}

	// thread-safe
	// Because CORBA doesn't support PREPARE, so we need to combine the
	// query arguments into one string and do underlying query without arguments
	commonSqlConn struct {
		driverName string
		datasource string
		beginTx    beginnable
		brk        breaker.Breaker
		accept     func(error) bool
	}

	sessionConn interface {
		Exec(query string, args ...interface{}) (sql.Result, error)
		Query(query string, args ...interface{}) (*sql.Rows, error)
	}

	statement struct {
		stmt *sql.Stmt
	}

	stmtConn interface {
		Exec(args ...interface{}) (sql.Result, error)
		Query(args ...interface{}) (*sql.Rows, error)
	}
)

// NewSqlConn returns a SqlConn with given driver name and datasource.
func NewSqlConn(driverName, datasource string, opts ...SqlOption) SqlConn {
	conn := &commonSqlConn{
		driverName: driverName,
		datasource: datasource,
		beginTx:    begin,
		brk:        breaker.NewBreaker(),
	}
	for _, opt := range opts {
		opt(conn)
	}

	return conn
}

func (db *commonSqlConn) Exec(q string, args ...interface{}) (result sql.Result, err error) {
	err = db.brk.DoWithAcceptable(func() error {
		var conn *sql.DB
		conn, err = getSqlConn(db.driverName, db.datasource)
		if err != nil {
			logInstanceError(db.datasource, err)
			return err
		}

		result, err = exec(conn, q, args...)
		return err
	}, db.acceptable)

	return
}

func (db *commonSqlConn) Prepare(query string) (stmt StmtSession, err error) {
	err = db.brk.DoWithAcceptable(func() error {
		var conn *sql.DB
		conn, err = getSqlConn(db.driverName, db.datasource)
		if err != nil {
			logInstanceError(db.datasource, err)
			return err
		}

		st, err := conn.Prepare(query)
		if err != nil {
			return err
		}

		stmt = statement{
			stmt: st,
		}
		return nil
	}, db.acceptable)

	return
}

func (db *commonSqlConn) QueryRow(v interface{}, q string, args ...interface{}) error {
	return db.queryRows(func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, q, args...)
}

func (db *commonSqlConn) QueryRowPartial(v interface{}, q string, args ...interface{}) error {
	return db.queryRows(func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, q, args...)
}

func (db *commonSqlConn) QueryRows(v interface{}, q string, args ...interface{}) error {
	return db.queryRows(func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, q, args...)
}

func (db *commonSqlConn) QueryRowsPartial(v interface{}, q string, args ...interface{}) error {
	return db.queryRows(func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

func (db *commonSqlConn) Transact(fn func(Session) error) error {
	return db.brk.DoWithAcceptable(func() error {
		return transact(db, db.beginTx, fn)
	}, db.acceptable)
}

func (db *commonSqlConn) acceptable(err error) bool {
	ok := err == nil || err == sql.ErrNoRows || err == sql.ErrTxDone
	if db.accept == nil {
		return ok
	}

	return ok || db.accept(err)
}

func (db *commonSqlConn) queryRows(scanner func(*sql.Rows) error, q string, args ...interface{}) error {
	var qerr error
	return db.brk.DoWithAcceptable(func() error {
		conn, err := getSqlConn(db.driverName, db.datasource)
		if err != nil {
			logInstanceError(db.datasource, err)
			return err
		}

		return query(conn, func(rows *sql.Rows) error {
			qerr = scanner(rows)
			return qerr
		}, q, args...)
	}, func(err error) bool {
		return qerr == err || db.acceptable(err)
	})
}
func (db *commonSqlConn) QueryRowsRaw(q string, args ...interface{}) (*[]map[string]interface{}, error) {
	var v []map[string]interface{}
	return &v, db.queryRows(func(rows *sql.Rows) error {
		list := make([]map[string]interface{}, 0)
		for rows.Next() {
			if item, err := sqlRowToMap(rows); err != nil {
				return err
			} else {
				list = append(list, *item)
			}
		}
		v = list
		return nil
	}, q, args...)
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

func sqlRowToMap(rows rowsScanner) (*map[string]interface{}, error) {
	columns, _ := rows.Columns()
	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}
	_ = rows.Scan(cache...)
	item := make(map[string]interface{})
	for i, data := range cache {
		item[columns[i]] = readval(*data.(*interface{})) //取实际类型
	}
	return &item, nil
}

func readval(value interface{}) interface{} {
	var key string
	if value == nil {
		return key
	}
	switch value.(type) {
	case float64:
		return value.(float64)
	case float32:
		return value.(float32)
	case int:
		return value.(int)
	case uint:
		return value.(uint)
	case int8:
		return value.(int8)
	case uint8:
		return value.(uint8)
	case int16:
		return value.(int16)
	case uint16:
		return value.(uint16)
	case int32:
		return value.(int32)
	case uint32:
		return value.(uint32)
	case int64:
		return value.(int64)
	case uint64:
		return value.(uint64)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		return value
	}
	return key
}
