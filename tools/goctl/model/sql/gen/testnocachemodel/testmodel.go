package testnocachemodel

import (
	"database/sql"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	testFieldNames          = builderx.FieldNames(&Test{})
	testRows                = strings.Join(testFieldNames, ",")
	testRowsExpectAutoSet   = strings.Join(stringx.Remove(testFieldNames, "id", "create_time", "update_time"), ",")
	testRowsWithPlaceHolder = strings.Join(stringx.Remove(testFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
)

type (
	TestModel struct {
		conn  sqlx.SqlConn
		table string
	}

	Test struct {
		Id         int64     `db:"id"`
		Nanosecond int64     `db:"nanosecond"`
		Data       string    `db:"data"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewTestModel(conn sqlx.SqlConn) *TestModel {
	return &TestModel{
		conn:  conn,
		table: "test",
	}
}

func (m *TestModel) Insert(data Test) (sql.Result, error) {
	query := `insert into ` + m.table + ` (` + testRowsExpectAutoSet + `) values (?, ?)`
	ret, err := m.conn.Exec(query, data.Nanosecond, data.Data)
	return ret, err
}

func (m *TestModel) FindOne(id int64) (*Test, error) {
	query := `select ` + testRows + ` from ` + m.table + ` where id = ? limit 1`
	var resp Test
	err := m.conn.QueryRow(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *TestModel) FindOneByNanosecond(nanosecond int64) (*Test, error) {
	var resp Test
	query := `select ` + testRows + ` from ` + m.table + ` where nanosecond = ? limit 1`
	err := m.conn.QueryRow(&resp, query, nanosecond)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *TestModel) Update(data Test) (sql.Result, error) {
	query := `update ` + m.table + ` set ` + testRowsWithPlaceHolder + ` where id = ?`
	ret, err := m.conn.Exec(query, data.Nanosecond, data.Data, data.Id)
	return ret, err
}

func (m *TestModel) Delete(id int64) error {
	query := `delete from ` + m.table + ` where id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}
