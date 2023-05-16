package sqlx

import (
	"database/sql"
	"io"
	"time"

	"github.com/zeromicro/go-zero/core/syncx"
)

const (
	maxIdleConns = 64
	maxOpenConns = 64
	maxLifetime  = time.Minute
)

var connManager = syncx.NewResourceManager()

func getCachedSqlConn(driverName, server string) (*sql.DB, error) {
	val, err := connManager.GetResource(server, func() (io.Closer, error) {
		conn, err := newDBConnection(driverName, server)
		if err != nil {
			return nil, err
		}

		return conn, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*sql.DB), nil
}

func getSqlConn(driverName, server string) (*sql.DB, error) {
	conn, err := getCachedSqlConn(driverName, server)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func newDBConnection(driverName, datasource string) (*sql.DB, error) {
	conn, err := sql.Open(driverName, datasource)
	if err != nil {
		return nil, err
	}

	// we need to do this until the issue https://github.com/golang/go/issues/9851 get fixed
	// discussed here https://github.com/go-sql-driver/mysql/issues/257
	// if the discussed SetMaxIdleTimeout methods added, we'll change this behavior
	// 8 means we can't have more than 8 goroutines to concurrently access the same database.
	conn.SetMaxIdleConns(maxIdleConns)
	conn.SetMaxOpenConns(maxOpenConns)
	conn.SetConnMaxLifetime(maxLifetime)

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}
