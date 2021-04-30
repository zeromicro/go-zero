package postgres

import (
	// imports the driver.
	_ "github.com/lib/pq"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

const postgresDriverName = "postgres"

// New returns a postgres connection.
func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(postgresDriverName, datasource, opts...)
}
