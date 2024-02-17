package postgres

import (
	// imports the driver, don't remove this comment, golint requires.
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const postgresDriverName = "pgx"

// New returns a postgres connection.
func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(postgresDriverName, datasource, opts...)
}
