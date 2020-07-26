package postgres

import (
	"zero/core/stores/sqlx"

	_ "github.com/lib/pq"
)

const postgreDriverName = "postgres"

func NewPostgre(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(postgreDriverName, datasource, opts...)
}
