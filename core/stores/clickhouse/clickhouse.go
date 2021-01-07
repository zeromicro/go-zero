package clickhouse

import (
	"github.com/3Rivers/go-zero/core/stores/sqlx"
	_ "github.com/ClickHouse/clickhouse-go"
)

const clickHouseDriverName = "clickhouse"

func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(clickHouseDriverName, datasource, opts...)
}
