package clickhouse

import (
	"zero/core/stores/sqlx"

	_ "github.com/kshvakov/clickhouse"
)

const clickHouseDriverName = "clickhouse"

func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(clickHouseDriverName, datasource, opts...)
}
