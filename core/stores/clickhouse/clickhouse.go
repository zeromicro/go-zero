package clickhouse

import (
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

const clickHouseDriverName = "clickhouse"

func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(clickHouseDriverName, datasource, opts...)
}
