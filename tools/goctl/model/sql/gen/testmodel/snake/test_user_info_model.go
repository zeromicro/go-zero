package snake

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/cache"
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

	cacheTestUserInfoIdPrefix         = "cache#TestUserInfo#id#"
	cacheTestUserInfoNanosecondPrefix = "cache#TestUserInfo#nanosecond#"
)

type (
	TestUserInfoModel struct {
		sqlc.CachedConn
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

func NewTestUserInfoModel(conn sqlx.SqlConn, c cache.CacheConf) *TestUserInfoModel {
	return &TestUserInfoModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "test_user_info",
	}
}

func (m *TestUserInfoModel) Insert(data TestUserInfo) (sql.Result, error) {
	testUserInfoNanosecondKey := fmt.Sprintf("%s%v", cacheTestUserInfoNanosecondPrefix, data.Nanosecond)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, testUserInfoRowsExpectAutoSet)
		return conn.Exec(query, data.Nanosecond, data.Data)
	}, testUserInfoNanosecondKey)
	return ret, err
}

func (m *TestUserInfoModel) FindOne(id int64) (*TestUserInfo, error) {
	testUserInfoIdKey := fmt.Sprintf("%s%v", cacheTestUserInfoIdPrefix, id)
	var resp TestUserInfo
	err := m.QueryRow(&resp, testUserInfoIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where id = ? limit 1", testUserInfoRows, m.table)
		return conn.QueryRow(v, query, id)
	})
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
	testUserInfoNanosecondKey := fmt.Sprintf("%s%v", cacheTestUserInfoNanosecondPrefix, nanosecond)
	var resp TestUserInfo
	err := m.QueryRowIndex(&resp, testUserInfoNanosecondKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where nanosecond = ? limit 1", testUserInfoRows, m.table)
		if err := conn.QueryRow(&resp, query, nanosecond); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
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
	testUserInfoIdKey := fmt.Sprintf("%s%v", cacheTestUserInfoIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where id = ?", m.table, testUserInfoRowsWithPlaceHolder)
		return conn.Exec(query, data.Nanosecond, data.Data, data.Id)
	}, testUserInfoIdKey)
	return err
}

func (m *TestUserInfoModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	testUserInfoIdKey := fmt.Sprintf("%s%v", cacheTestUserInfoIdPrefix, id)
	testUserInfoNanosecondKey := fmt.Sprintf("%s%v", cacheTestUserInfoNanosecondPrefix, data.Nanosecond)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where id = ?", m.table)
		return conn.Exec(query, id)
	}, testUserInfoIdKey, testUserInfoNanosecondKey)
	return err
}

func (m *TestUserInfoModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheTestUserInfoIdPrefix, primary)
}

func (m *TestUserInfoModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where id = ? limit 1", testUserInfoRows, m.table)
	return conn.QueryRow(v, query, primary)
}
