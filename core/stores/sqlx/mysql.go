package sqlx

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

const (
	mysqlDriverName           = "mysql"
	duplicateEntryCode uint16 = 1062
)

// NewMysql returns a mysql connection.
func NewMysql(datasource string, opts ...SqlOption) SqlConn {
	opts = append([]SqlOption{withMysqlAcceptable()}, opts...)
	return NewSqlConn(mysqlDriverName, datasource, opts...)
}

func mysqlAcceptable(err error) bool {
	if err == nil {
		return true
	}

	var myerr *mysql.MySQLError
	ok := errors.As(err, &myerr)
	if !ok {
		return false
	}

	switch myerr.Number {
	case duplicateEntryCode:
		return true
	default:
		return false
	}
}

func withMysqlAcceptable() SqlOption {
	return WithAcceptable(mysqlAcceptable)
}
