package sqlx

import (
	"database/sql"
	"io"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

const mockedDatasource = "sqlmock"

func init() {
	logx.Disable()
}

func TestSqlConn(t *testing.T) {
	mock, err := buildConn()
	assert.Nil(t, err)
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

func buildConn() (mock sqlmock.Sqlmock, err error) {
	_, err = connManager.GetResource(mockedDatasource, func() (io.Closer, error) {
		var db *sql.DB
		var err error
		db, mock, err = sqlmock.New()
		return &pingedDB{
			DB: db,
		}, err
	})

	return
}

func TestSalve(t *testing.T) {
	salve1 := "salve1"
	db, salve1Mock, err := sqlmock.New()
	_, err = connManager.GetResource(salve1, func() (io.Closer, error) {
		return &pingedDB{
			DB: db,
		}, err
	})
	assert.NoError(t, err)

	salve1Mock.ExpectQuery("any").WillReturnRows(
		sqlmock.NewRows([]string{"a"}).AddRow("foo"),
	)

	mysql := NewMysql("", WithSlaves([]string{salve1}), WithRandomPicker())

	var result string
	err = mysql.QueryRow(&result, "any")
	assert.NoError(t, err)
	assert.EqualValues(t, "foo", result)

	_, err = mysql.Exec("any")
	assert.Error(t, err)
}

func TestWithSlaves(t *testing.T) {
	slaves := WithSlaves([]string{"1", "2"})
	assert.NotNil(t, slaves)
}

func TestWithRandomPicker(t *testing.T) {
	slaves := WithRandomPicker()
	assert.NotNil(t, slaves)
}

func TestWithRoundRobinPicker(t *testing.T) {
	slaves := WithRoundRobinPicker()
	assert.NotNil(t, slaves)
}

func TestWithWeightRandomPicker(t *testing.T) {
	slaves := WithWeightRandomPicker([]int{1, 2})
	assert.NotNil(t, slaves)
}

func TestWithWeightRoundRobinPicker(t *testing.T) {
	slaves := WithWeightRoundRobinPicker([]int{1, 2})
	assert.NotNil(t, slaves)
}
