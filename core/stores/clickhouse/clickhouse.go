package clickhouse

import (
	// imports the driver.
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

const clickHouseDriverName = "clickhouse"

// New returns a clickhouse connection.
func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(clickHouseDriverName, datasource, opts...)
}
