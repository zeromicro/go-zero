package postgres

import (
	"github.com/3Rivers/go-zero/core/stores/sqlx"
	_ "github.com/lib/pq"
)

const postgresDriverName = "postgres"

func NewPostgres(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(postgresDriverName, datasource, opts...)
}
