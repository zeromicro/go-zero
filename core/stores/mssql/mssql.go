package mssql

import (
	// imports the driver.
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

const mssqlDriverName = "mssql"

// New returns a mssql connection.
func NewMssqlConn(datasource string, opts ...sqlx.SqlOption) sqlx.SqlConn {
	return sqlx.NewSqlConn(mssqlDriverName, datasource, opts...)
}
