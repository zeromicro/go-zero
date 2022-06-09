package clickhouse

import (
	// imports the driver, don't remove this comment, golint requires.
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const clickHouseDriverName = "clickhouse"

// New returns a clickhouse connection.
func New(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(clickHouseDriverName, datasource, opts...)
}
