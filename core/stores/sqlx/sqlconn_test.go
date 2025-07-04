package sqlx

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/dbtest"
	"github.com/zeromicro/go-zero/core/trace/tracetest"
)

const mockedDatasource = "sqlmock"

func init() {
	logx.Disable()
}

func TestSqlConn(t *testing.T) {
	me := tracetest.NewInMemoryExporter(t)
	mock, err := buildConn()
	assert.Nil(t, err)
	mock.ExpectExec("any")
	mock.ExpectQuery("any").WillReturnRows(sqlmock.NewRows([]string{"foo"}))
	conn := NewMysql(mockedDatasource)
	db, err := conn.RawDB()
	assert.Nil(t, err)
	rawConn := NewSqlConnFromDB(db, withMysqlAcceptable())
	badConn := NewMysql("badsql")
	_, err = conn.Exec("any", "value")
	assert.NotNil(t, err)
	_, err = badConn.Exec("any", "value")
	assert.NotNil(t, err)
	_, err = rawConn.Prepare("any")
	assert.NotNil(t, err)
	_, err = badConn.Prepare("any")
	assert.NotNil(t, err)
	var val string
	assert.NotNil(t, conn.QueryRow(&val, "any"))
	assert.NotNil(t, badConn.QueryRow(&val, "any"))
	assert.NotNil(t, conn.QueryRowPartial(&val, "any"))
	assert.NotNil(t, badConn.QueryRowPartial(&val, "any"))
	assert.NotNil(t, conn.QueryRows(&val, "any"))
	assert.NotNil(t, badConn.QueryRows(&val, "any"))
	assert.NotNil(t, conn.QueryRowsPartial(&val, "any"))
	assert.NotNil(t, badConn.QueryRowsPartial(&val, "any"))
	assert.NotNil(t, conn.Transact(func(session Session) error {
		return nil
	}))
	assert.NotNil(t, badConn.Transact(func(session Session) error {
		return nil
	}))
	assert.Equal(t, 14, len(me.GetSpans()))
}

func TestSqlConn_RawDB(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(rows)
		conn := NewSqlConnFromDB(db)
		var val string
		assert.NoError(t, conn.QueryRow(&val, "any"))
		assert.Equal(t, "bar", val)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(rows)
		conn := NewSqlConnFromDB(db)
		var val string
		assert.NoError(t, conn.QueryRowPartial(&val, "any"))
		assert.Equal(t, "bar", val)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"any"}).AddRow("foo").AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(rows)
		conn := NewSqlConnFromDB(db)
		var vals []string
		assert.NoError(t, conn.QueryRows(&vals, "any"))
		assert.ElementsMatch(t, []string{"foo", "bar"}, vals)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"any"}).AddRow("foo").AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(rows)
		conn := NewSqlConnFromDB(db)
		var vals []string
		assert.NoError(t, conn.QueryRowsPartial(&vals, "any"))
		assert.ElementsMatch(t, []string{"foo", "bar"}, vals)
	})
}

func TestSqlConn_Errors(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		conn := NewSqlConnFromDB(db)
		conn.(*commonSqlConn).connProv = func(ctx context.Context) (*sql.DB, error) {
			return nil, errors.New("error")
		}
		_, err := conn.Prepare("any")
		assert.Error(t, err)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec("any").WillReturnError(breaker.ErrServiceUnavailable)
		conn := NewSqlConnFromDB(db)
		_, err := conn.Exec("any")
		assert.Error(t, err)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any").WillReturnError(breaker.ErrServiceUnavailable)
		conn := NewSqlConnFromDB(db)
		_, err := conn.Prepare("any")
		assert.Error(t, err)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectRollback()
		conn := NewSqlConnFromDB(db)
		err := conn.Transact(func(session Session) error {
			return breaker.ErrServiceUnavailable
		})
		assert.Equal(t, breaker.ErrServiceUnavailable, err)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectQuery("any").WillReturnError(breaker.ErrServiceUnavailable)
		conn := NewSqlConnFromDB(db)
		var vals []string
		err := conn.QueryRows(&vals, "any")
		assert.Equal(t, breaker.ErrServiceUnavailable, err)
	})
}

func TestStatement(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any").WillBeClosed()

		conn := NewSqlConnFromDB(db)
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)
		assert.NoError(t, stmt.Close())
	})

	dbtest.RunTxTest(t, func(tx *sql.Tx, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any").WillBeClosed()

		stmt, err := tx.Prepare("any")
		assert.NoError(t, err)
		st := statement{
			query: "foo",
			stmt:  stmt,
			brk:   breaker.NopBreaker(),
		}
		assert.NoError(t, st.Close())
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any")
		mock.ExpectExec("any").WillReturnResult(sqlmock.NewResult(2, 3))

		conn := NewSqlConnFromDB(db)
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)
		res, err := stmt.Exec()
		assert.NoError(t, err)
		lastInsertID, err := res.LastInsertId()
		assert.NoError(t, err)
		assert.Equal(t, int64(2), lastInsertID)
		rowsAffected, err := res.RowsAffected()
		assert.NoError(t, err)
		assert.Equal(t, int64(3), rowsAffected)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any")
		row := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(row)

		conn := NewSqlConnFromDB(db)
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)

		var val string
		err = stmt.QueryRow(&val)
		assert.NoError(t, err)
		assert.Equal(t, "bar", val)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any")
		row := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(row)

		conn := NewSqlConnFromDB(db)
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)

		var val string
		err = stmt.QueryRowPartial(&val)
		assert.NoError(t, err)
		assert.Equal(t, "bar", val)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any")
		rows := sqlmock.NewRows([]string{"any"}).AddRow("foo").AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(rows)

		conn := NewSqlConnFromDB(db)
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)

		var vals []string
		assert.NoError(t, stmt.QueryRows(&vals))
		assert.ElementsMatch(t, []string{"foo", "bar"}, vals)
	})

	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		mock.ExpectPrepare("any")
		rows := sqlmock.NewRows([]string{"any"}).AddRow("foo").AddRow("bar")
		mock.ExpectQuery("any").WillReturnRows(rows)

		conn := NewSqlConnFromDB(db)
		stmt, err := conn.Prepare("any")
		assert.NoError(t, err)

		var vals []string
		assert.NoError(t, stmt.QueryRowsPartial(&vals))
		assert.ElementsMatch(t, []string{"foo", "bar"}, vals)
	})
}

func TestBreakerWithFormatError(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		conn := NewSqlConnFromDB(db, withMysqlAcceptable())
		for i := 0; i < 1000; i++ {
			var val string
			if !assert.NotEqual(t, breaker.ErrServiceUnavailable,
				conn.QueryRow(&val, "any ?, ?", "foo")) {
				break
			}
		}
	})
}

func TestBreakerWithScanError(t *testing.T) {
	dbtest.RunTest(t, func(db *sql.DB, mock sqlmock.Sqlmock) {
		conn := NewSqlConnFromDB(db, withMysqlAcceptable())
		for i := 0; i < 1000; i++ {
			rows := sqlmock.NewRows([]string{"foo"}).AddRow("bar")
			mock.ExpectQuery("any").WillReturnRows(rows)
			var val int
			if !assert.NotEqual(t, breaker.ErrServiceUnavailable, conn.QueryRow(&val, "any")) {
				break
			}
		}
	})
}

func TestWithAcceptable(t *testing.T) {
	var (
		acceptableErr  = errors.New("acceptable")
		acceptableErr2 = errors.New("acceptable2")
		acceptableErr3 = errors.New("acceptable3")
	)
	opts := []SqlOption{
		WithAcceptable(func(err error) bool {
			if err == nil {
				return true
			}
			return errors.Is(err, acceptableErr)
		}),
		WithAcceptable(func(err error) bool {
			if err == nil {
				return true
			}
			return errors.Is(err, acceptableErr2)
		}),
		WithAcceptable(func(err error) bool {
			if err == nil {
				return true
			}
			return errors.Is(err, acceptableErr3)
		}),
	}

	var conn = &commonSqlConn{}
	for _, opt := range opts {
		opt(conn)
	}

	assert.True(t, conn.accept(nil))
	assert.False(t, conn.accept(assert.AnError))
	assert.True(t, conn.accept(acceptableErr))
	assert.True(t, conn.accept(acceptableErr2))
	assert.True(t, conn.accept(acceptableErr3))
}

func TestProvider(t *testing.T) {
	defer func() {
		_ = connManager.Close()
	}()

	mainDSN := "main:password@tcp(127.0.0.1:3306)/main_db"
	replicasDSN := []string{
		"replica_one:pwd@tcp(localhost:3306)/replica_one",
		"replica_two:pwd@tcp(localhost:3306)/replica_two",
		"replica_three:pwd@tcp(localhost:3306)/replica_three",
	}

	mainDB, err := connManager.GetResource(mainDSN, func() (io.Closer, error) { return sql.Open(mysqlDriverName, mainDSN) })
	assert.Nil(t, err)
	assert.NotNil(t, mainDB)
	replicaOneDB, err := connManager.GetResource(replicasDSN[0], func() (io.Closer, error) { return sql.Open(mysqlDriverName, replicasDSN[0]) })
	assert.Nil(t, err)
	assert.NotNil(t, replicaOneDB)
	replicaTwoDB, err := connManager.GetResource(replicasDSN[1], func() (io.Closer, error) { return sql.Open(mysqlDriverName, replicasDSN[1]) })
	assert.Nil(t, err)
	assert.NotNil(t, replicaTwoDB)
	replicaThreeDB, err := connManager.GetResource(replicasDSN[2], func() (io.Closer, error) { return sql.Open(mysqlDriverName, replicasDSN[2]) })
	assert.Nil(t, err)
	assert.NotNil(t, replicaThreeDB)

	sc := &commonSqlConn{
		policy: PolicyRoundRobin,
		index:  0,
	}
	sc.connProv = getConnProvider(sc, mysqlDriverName, mainDSN)

	ctx := context.Background()
	db, err := sc.connProv(ctx)
	assert.Nil(t, err)
	assert.Equal(t, mainDB, db)

	ctx = WithWriteMode(ctx)
	db, err = sc.connProv(ctx)
	assert.Nil(t, err)
	assert.Equal(t, mainDB, db)

	ctx = WithReadMode(ctx)
	sc.replicas = []string{replicasDSN[0]}
	db, err = sc.connProv(ctx)
	assert.Nil(t, err)
	assert.Equal(t, replicaOneDB, db)

	// default policy is round-robin
	sc.replicas = replicasDSN
	replicas := []io.Closer{replicaOneDB, replicaTwoDB, replicaThreeDB}
	for i := 0; i < len(sc.replicas); i++ {
		db, err = sc.connProv(ctx)
		assert.Nil(t, err)
		assert.Equal(t, replicas[i], db)
	}

	// random policy
	sc.policy = PolicyRandom
	for i := 0; i < len(sc.replicas); i++ {
		db, err = sc.connProv(ctx)
		assert.Nil(t, err)
		assert.Contains(t, replicas, db)
	}

	// unknown policy
	sc.policy = "unknown"
	_, err = sc.connProv(ctx)
	assert.NotNil(t, err)

	// empty policy transforms to round-robin
	sc.policy = ""
	for i := 0; i < len(sc.replicas); i++ {
		db, err = sc.connProv(ctx)
		assert.Nil(t, err)
		assert.Equal(t, replicas[i], db)
	}
}

func TestWithPolicy(t *testing.T) {
	replicas := []string{"replica1", "replica2", "replica3"}
	opts := []SqlOption{
		WithReplicas(replicas...),
	}
	conn := &commonSqlConn{}
	for _, opt := range opts {
		opt(conn)
	}
	assert.Equal(t, conn.replicas, replicas)
}

func TestWithReplicas(t *testing.T) {
	opts := []SqlOption{
		WithPolicy(PolicyRandom),
	}
	conn := &commonSqlConn{}
	for _, opt := range opts {
		opt(conn)
	}
	assert.Equal(t, PolicyRandom, conn.policy)
}

func buildConn() (mock sqlmock.Sqlmock, err error) {
	_, err = connManager.GetResource(mockedDatasource, func() (io.Closer, error) {
		var db *sql.DB
		var err error
		db, mock, err = sqlmock.New()
		return db, err
	})
	return
}
