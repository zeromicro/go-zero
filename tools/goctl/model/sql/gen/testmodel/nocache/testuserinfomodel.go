package nocache

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	testUserInfoFieldNames          = builderx.FieldNames(&TestUserInfo{})
	testUserInfoRows                = strings.Join(testUserInfoFieldNames, ",")
	testUserInfoRowsExpectAutoSet   = strings.Join(stringx.Remove(testUserInfoFieldNames, "id", "create_time", "update_time"), ",")
	testUserInfoRowsWithPlaceHolder = strings.Join(stringx.Remove(testUserInfoFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
)

type (
	TestUserInfoModel struct {
		conn  sqlx.SqlConn
		table string
	}

	TestUserInfo struct {
		Id         int64     `db:"id"`
		Nanosecond int64     `db:"nanosecond"`
		Data       string    `db:"data"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewTestUserInfoModel(conn sqlx.SqlConn) *TestUserInfoModel {
	return &TestUserInfoModel{
		conn:  conn,
		table: "test_user_info",
	}
}

func (m *TestUserInfoModel) Insert(data TestUserInfo) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, testUserInfoRowsExpectAutoSet)
	ret, err := m.conn.Exec(query, data.Nanosecond, data.Data)
	return ret, err
}

func (m *TestUserInfoModel) FindOne(id int64) (*TestUserInfo, error) {
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", testUserInfoRows, m.table)
	var resp TestUserInfo
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

func (m *TestUserInfoModel) FindOneByNanosecond(nanosecond int64) (*TestUserInfo, error) {
	var resp TestUserInfo
	query := fmt.Sprintf("select %s from %s where nanosecond = ? limit 1", testUserInfoRows, m.table)
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

func (m *TestUserInfoModel) Update(data TestUserInfo) error {
	query := fmt.Sprintf("update %s set %s where id = ?", m.table, testUserInfoRowsWithPlaceHolder)
	_, err := m.conn.Exec(query, data.Nanosecond, data.Data, data.Id)
	return err
}

func (m *TestUserInfoModel) Delete(id int64) error {
	query := fmt.Sprintf("delete from %s where id = ?", m.table)
	_, err := m.conn.Exec(query, id)
	return err
}
