package clickhouse

import (
	// imports the driver, don't remove this comment, golint requires.
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const clickHouseDriverName = "clickhouse"

// New returns a clickhouse connection.
func New(datasource string, dialProvider sqlx.DialProvider, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(clickHouseDriverName, datasource, dialProvider, opts...)
}
