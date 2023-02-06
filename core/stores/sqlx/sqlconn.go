package sqlx

import (
	"context"
	"database/sql"

	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
)

// spanName is used to identify the span name for the SQL execution.
const spanName = "sql"

// ErrNotFound is an alias of sql.ErrNoRows
var ErrNotFound = sql.ErrNoRows

type (
	// Session stands for raw connections or transaction sessions
	Session interface {
		Exec(query string, args ...any) (sql.Result, error)
		ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error)
		Prepare(query string) (StmtSession, error)
		PrepareCtx(ctx context.Context, query string) (StmtSession, error)
		QueryRow(v any, query string, args ...any) error
		QueryRowCtx(ctx context.Context, v any, query string, args ...any) error
		QueryRowPartial(v any, query string, args ...any) error
		QueryRowPartialCtx(ctx context.Context, v any, query string, args ...any) error
		QueryRows(v any, query string, args ...any) error
		QueryRowsCtx(ctx context.Context, v any, query string, args ...any) error
		QueryRowsPartial(v any, query string, args ...any) error
		QueryRowsPartialCtx(ctx context.Context, v any, query string, args ...any) error
	}

	// SqlConn only stands for raw connections, so Transact method can be called.
	SqlConn interface {
		Session
		// RawDB is for other ORM to operate with, use it with caution.
		// Notice: don't close it.
		RawDB() (*sql.DB, error)
		Transact(fn func(Session) error) error
		TransactCtx(ctx context.Context, fn func(context.Context, Session) error) error
	}

	// SqlOption defines the method to customize a sql connection.
	SqlOption func(*commonSqlConn)

	// StmtSession interface represents a session that can be used to execute statements.
	StmtSession interface {
		Close() error
		Exec(args ...any) (sql.Result, error)
		ExecCtx(ctx context.Context, args ...any) (sql.Result, error)
		QueryRow(v any, args ...any) error
		QueryRowCtx(ctx context.Context, v any, args ...any) error
		QueryRowPartial(v any, args ...any) error
		QueryRowPartialCtx(ctx context.Context, v any, args ...any) error
		QueryRows(v any, args ...any) error
		QueryRowsCtx(ctx context.Context, v any, args ...any) error
		QueryRowsPartial(v any, args ...any) error
		QueryRowsPartialCtx(ctx context.Context, v any, args ...any) error
	}

	// thread-safe
	// Because CORBA doesn't support PREPARE, so we need to combine the
	// query arguments into one string and do underlying query without arguments
	commonSqlConn struct {
		connProv   connProvider
		onError    func(error)
		beginTx    beginnable
		accept     func(error) bool
		picker     picker
		fnSlaves   fnSlaves
		driverName string
	}

	slave struct {
		datasource string
		driverName string
		brk        breaker.Breaker
	}

	connProvider func(ctx context.Context, onlyRead bool) (brk breaker.Breaker, db *sql.DB, err error)

	sessionConn interface {
		Exec(query string, args ...any) (sql.Result, error)
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		Query(query string, args ...any) (*sql.Rows, error)
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	}

	statement struct {
		query string
		stmt  *sql.Stmt
	}

	stmtConn interface {
		Exec(args ...any) (sql.Result, error)
		ExecContext(ctx context.Context, args ...any) (sql.Result, error)
		Query(args ...any) (*sql.Rows, error)
		QueryContext(ctx context.Context, args ...any) (*sql.Rows, error)
	}
)

// NewSqlConn returns a SqlConn with given driver name and datasource.
func NewSqlConn(driverName, datasource string, opts ...SqlOption) SqlConn {
	conn := &commonSqlConn{
		onError: func(err error) {
			logInstanceError(datasource, err)
		},
		beginTx:    begin,
		driverName: driverName,
	}

	brk := breaker.NewBreaker()
	conn.connProv = func(ctx context.Context, onlyRead bool) (breaker.Breaker, *sql.DB, error) {
		if onlyRead && conn.picker != nil && conn.fnSlaves != nil {
			salve, err := conn.picker.pick()
			if err == nil {
				var db *sql.DB
				db, err = salve.getDB()

				return salve.getBreaker(), db, err
			}

			logx.WithContext(ctx).Error(err)
		}

		db, err := getSqlConn(driverName, datasource)
		if err != nil {
			return brk, nil, err
		}

		return brk, db, nil
	}

	for _, opt := range opts {
		opt(conn)
	}

	return conn
}

// NewSqlConnFromDB returns a SqlConn with the given sql.DB.
// Use it with caution, it's provided for other ORM to interact with.
func NewSqlConnFromDB(db *sql.DB, opts ...SqlOption) SqlConn {
	brk := breaker.NewBreaker()
	conn := &commonSqlConn{
		connProv: func(ctx context.Context, onlyRead bool) (breaker.Breaker, *sql.DB, error) {
			return brk, db, nil
		},
		onError: func(err error) {
			logx.Errorf("Error on getting sql instance: %v", err)
		},
		beginTx: begin,
	}

	for _, opt := range opts {
		opt(conn)
	}

	return conn
}

func (db *commonSqlConn) Exec(q string, args ...any) (result sql.Result, err error) {
	return db.ExecCtx(context.Background(), q, args...)
}

func (db *commonSqlConn) ExecCtx(ctx context.Context, q string, args ...any) (
	result sql.Result, err error) {
	ctx, span := startSpan(ctx, "Exec")
	defer func() {
		endSpan(span, err)
	}()

	brk, conn, err := db.connProv(ctx, false)
	err = brk.DoWithAcceptable(func() error {
		if err != nil {
			db.onError(err)
			return err
		}

		result, err = exec(ctx, conn, q, args...)
		return err
	}, db.acceptable)

	if err == breaker.ErrServiceUnavailable {
		metricReqErr.Inc("Exec", "breaker")
	}

	return
}

func (db *commonSqlConn) Prepare(query string) (stmt StmtSession, err error) {
	return db.PrepareCtx(context.Background(), query)
}

func (db *commonSqlConn) PrepareCtx(ctx context.Context, query string) (stmt StmtSession, err error) {
	ctx, span := startSpan(ctx, "Prepare")
	defer func() {
		endSpan(span, err)
	}()

	brk, conn, err := db.connProv(ctx, false)
	err = brk.DoWithAcceptable(func() error {
		if err != nil {
			db.onError(err)
			return err
		}

		st, err := conn.PrepareContext(ctx, query)
		if err != nil {
			return err
		}

		stmt = statement{
			query: query,
			stmt:  st,
		}
		return nil
	}, db.acceptable)
	if err == breaker.ErrServiceUnavailable {
		metricReqErr.Inc("Prepare", "breaker")
	}

	return
}

func (db *commonSqlConn) QueryRow(v any, q string, args ...any) error {
	return db.QueryRowCtx(context.Background(), v, q, args...)
}

func (db *commonSqlConn) QueryRowCtx(ctx context.Context, v any, q string,
	args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRow")
	defer func() {
		endSpan(span, err)
	}()

	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, q, args...)
}

func (db *commonSqlConn) QueryRowPartial(v any, q string, args ...any) error {
	return db.QueryRowPartialCtx(context.Background(), v, q, args...)
}

func (db *commonSqlConn) QueryRowPartialCtx(ctx context.Context, v any,
	q string, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowPartial")
	defer func() {
		endSpan(span, err)
	}()

	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, q, args...)
}

func (db *commonSqlConn) QueryRows(v any, q string, args ...any) error {
	return db.QueryRowsCtx(context.Background(), v, q, args...)
}

func (db *commonSqlConn) QueryRowsCtx(ctx context.Context, v any, q string,
	args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRows")
	defer func() {
		endSpan(span, err)
	}()

	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, q, args...)
}

func (db *commonSqlConn) QueryRowsPartial(v any, q string, args ...any) error {
	return db.QueryRowsPartialCtx(context.Background(), v, q, args...)
}

func (db *commonSqlConn) QueryRowsPartialCtx(ctx context.Context, v any,
	q string, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowsPartial")
	defer func() {
		endSpan(span, err)
	}()

	return db.queryRows(ctx, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, q, args...)
}

func (db *commonSqlConn) RawDB() (*sql.DB, error) {
	_, conn, err := db.connProv(context.Background(), false)
	return conn, err
}

func (db *commonSqlConn) Transact(fn func(Session) error) error {
	return db.TransactCtx(context.Background(), func(_ context.Context, session Session) error {
		return fn(session)
	})
}

func (db *commonSqlConn) TransactCtx(ctx context.Context, fn func(context.Context, Session) error) (err error) {
	ctx, span := startSpan(ctx, "Transact")
	defer func() {
		endSpan(span, err)
	}()

	brk, conn, err := db.connProv(ctx, false)
	err = brk.DoWithAcceptable(func() error {
		if err != nil {
			db.onError(err)
			return err
		}
		return transactOnConn(ctx, conn, db.beginTx, fn)
	}, db.acceptable)
	if err == breaker.ErrServiceUnavailable {
		metricReqErr.Inc("Transact", "breaker")
	}

	return
}

func (db *commonSqlConn) acceptable(err error) bool {
	ok := err == nil || err == sql.ErrNoRows || err == sql.ErrTxDone || err == context.Canceled
	if db.accept == nil {
		return ok
	}

	return ok || db.accept(err)
}

func (db *commonSqlConn) queryRows(ctx context.Context, scanner func(*sql.Rows) error,
	q string, args ...any) (err error) {
	var qerr error
	brk, conn, err := db.connProv(ctx, true)
	err = brk.DoWithAcceptable(func() error {
		if err != nil {
			db.onError(err)
			return err
		}

		return query(ctx, conn, func(rows *sql.Rows) error {
			qerr = scanner(rows)
			return qerr
		}, q, args...)
	}, func(err error) bool {
		return qerr == err || db.acceptable(err)
	})
	if err == breaker.ErrServiceUnavailable {
		metricReqErr.Inc("queryRows", "breaker")
	}

	return
}

func (s statement) Close() error {
	return s.stmt.Close()
}

func (s statement) Exec(args ...any) (sql.Result, error) {
	return s.ExecCtx(context.Background(), args...)
}

func (s statement) ExecCtx(ctx context.Context, args ...any) (result sql.Result, err error) {
	ctx, span := startSpan(ctx, "Exec")
	defer func() {
		endSpan(span, err)
	}()

	return execStmt(ctx, s.stmt, s.query, args...)
}

func (s statement) QueryRow(v any, args ...any) error {
	return s.QueryRowCtx(context.Background(), v, args...)
}

func (s statement) QueryRowCtx(ctx context.Context, v any, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRow")
	defer func() {
		endSpan(span, err)
	}()

	return queryStmt(ctx, s.stmt, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, true)
	}, s.query, args...)
}

func (s statement) QueryRowPartial(v any, args ...any) error {
	return s.QueryRowPartialCtx(context.Background(), v, args...)
}

func (s statement) QueryRowPartialCtx(ctx context.Context, v any, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowPartial")
	defer func() {
		endSpan(span, err)
	}()

	return queryStmt(ctx, s.stmt, func(rows *sql.Rows) error {
		return unmarshalRow(v, rows, false)
	}, s.query, args...)
}

func (s statement) QueryRows(v any, args ...any) error {
	return s.QueryRowsCtx(context.Background(), v, args...)
}

func (s statement) QueryRowsCtx(ctx context.Context, v any, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRows")
	defer func() {
		endSpan(span, err)
	}()

	return queryStmt(ctx, s.stmt, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, true)
	}, s.query, args...)
}

func (s statement) QueryRowsPartial(v any, args ...any) error {
	return s.QueryRowsPartialCtx(context.Background(), v, args...)
}

func (s statement) QueryRowsPartialCtx(ctx context.Context, v any, args ...any) (err error) {
	ctx, span := startSpan(ctx, "QueryRowsPartial")
	defer func() {
		endSpan(span, err)
	}()

	return queryStmt(ctx, s.stmt, func(rows *sql.Rows) error {
		return unmarshalRows(v, rows, false)
	}, s.query, args...)
}

func newSlaves(driverName string, datasourceGroup []string) []slave {
	slaves := make([]slave, 0, len(datasourceGroup))
	for _, datasource := range datasourceGroup {
		slaves = append(slaves, slave{
			datasource: datasource,
			brk:        breaker.NewBreaker(breaker.WithName(datasource)),
			driverName: driverName,
		})
	}

	return slaves
}

func (s *slave) getBreaker() breaker.Breaker {
	return s.brk
}

func (s *slave) getDB() (*sql.DB, error) {
	return getSqlConn(s.driverName, s.datasource)
}

// WithSlaves returns a SqlOption that contains the slave database source.
func WithSlaves(dataSourceGroup []string) SqlOption {
	return func(conn *commonSqlConn) {
		conn.fnSlaves = func() []slave {
			return newSlaves(conn.driverName, dataSourceGroup)
		}
	}
}

// WithRandomPicker returns a SqlOption that contains the randomPicker algorithm.
func WithRandomPicker() SqlOption {
	return func(conn *commonSqlConn) {
		conn.picker = newRandomPicker(conn.fnSlaves)
	}
}

// WithWeightRandomPicker returns a SqlOption that contains the weightRandomPicker algorithm.
func WithWeightRandomPicker(weights []int) SqlOption {
	return func(conn *commonSqlConn) {
		conn.picker = newWeightRandomPicker(weights, conn.fnSlaves)
	}
}

// WithRoundRobinPicker returns a SqlOption that contains the roundRobinPicker algorithm.
func WithRoundRobinPicker() SqlOption {
	return func(conn *commonSqlConn) {
		conn.picker = newRoundRobinPicker(conn.fnSlaves)
	}
}

// WithWeightRoundRobinPicker returns a SqlOption that contains the weightRoundRobinPicker algorithm.
func WithWeightRoundRobinPicker(weights []int) SqlOption {
	return func(conn *commonSqlConn) {
		conn.picker = newWeightRoundRobinPicker(weights, conn.fnSlaves)
	}
}