package sqlx

import (
	"database/sql"
	"fmt"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/core/logx"
)

const mockedDatasource = "sqlmock"

func init() {
	logx.Disable()
}

func TestClose(t *testing.T) {
	conn := NewMysql("root:password@tcp(localhost:3306)/db?charset=utf8mb4&parseTime=true")
	db, _ := conn.RawDB()
	db.Close()
	conn = NewMysql("root:password@tcp(localhost:3306)/db?charset=utf8mb4&parseTime=true")
	if _, err := conn.Exec("select 1"); err != nil {
		fmt.Println(err.Error())
	}
}

func TestOldClose(t *testing.T) {
	mock := buildConn()
	mock.ExpectClose()
	mock.ExpectExec("select 1").WillReturnResult(nil)
	conn := NewMysql(mockedDatasource)
	db, err := conn.RawDB()
	assert.Nil(t, err)
	assert.Nil(t, db.Close())
	conn = NewMysql(mockedDatasource)
	_, err = conn.Exec("select 1")
	assert.NotNil(t, err)
}

func TestNewClose(t *testing.T) {
	mock := buildConn()
	mock.ExpectClose()

	conn := NewMysql(mockedDatasource)
	err := conn.DBClose(mockedDatasource)
	assert.Nil(t, err)
	conn = NewMysql(mockedDatasource)
	mock = buildConn()
	db, err := conn.RawDB()
	err = db.Ping()
	assert.Nil(t, err)
}

func TestSqlConn(t *testing.T) {
	mock := buildConn()
	mock.ExpectExec("any")
	mock.ExpectQuery("any").WillReturnRows(sqlmock.NewRows([]string{"foo"}))
	conn := NewMysql(mockedDatasource)
	db, err := conn.RawDB()
	assert.Nil(t, err)
	rawConn := NewSqlConnFromDB(db, withMysqlAcceptable())
	badConn := NewMysql("badsql")
	_, err = conn.Exec("any", "value")
	assert.NotNil(t, err)
	_, err = badConn.Exec("any", "value")
	assert.NotNil(t, err)
	_, err = rawConn.Prepare("any")
	assert.NotNil(t, err)
	_, err = badConn.Prepare("any")
	assert.NotNil(t, err)
	var val string
	assert.NotNil(t, conn.QueryRow(&val, "any"))
	assert.NotNil(t, badConn.QueryRow(&val, "any"))
	assert.NotNil(t, conn.QueryRowPartial(&val, "any"))
	assert.NotNil(t, badConn.QueryRowPartial(&val, "any"))
	assert.NotNil(t, conn.QueryRows(&val, "any"))
	assert.NotNil(t, badConn.QueryRows(&val, "any"))
	assert.NotNil(t, conn.QueryRowsPartial(&val, "any"))
	assert.NotNil(t, badConn.QueryRowsPartial(&val, "any"))
	assert.NotNil(t, conn.Transact(func(session Session) error {
		return nil
	}))
	assert.NotNil(t, badConn.Transact(func(session Session) error {
		return nil
	}))
}

func buildConn() (mock sqlmock.Sqlmock) {
	connManager.GetResource(mockedDatasource, func() (io.Closer, error) {
		var db *sql.DB
		var err error
		db, mock, err = sqlmock.New()
		return &pingedDB{
			DB: db,
		}, err
	})

	return
}
