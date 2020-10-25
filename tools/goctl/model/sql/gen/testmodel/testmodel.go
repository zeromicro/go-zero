package testmodel

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
	testFieldNames          = builderx.FieldNames(&Test{})
	testRows                = strings.Join(testFieldNames, ",")
	testRowsExpectAutoSet   = strings.Join(stringx.Remove(testFieldNames, "id", "create_time", "update_time"), ",")
	testRowsWithPlaceHolder = strings.Join(stringx.Remove(testFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"

	cacheTestIdPrefix         = "cache#Test#id#"
	cacheTestNanosecondPrefix = "cache#Test#nanosecond#"
)

type (
	TestModel struct {
		sqlc.CachedConn
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

func NewTestModel(conn sqlx.SqlConn, c cache.CacheConf) *TestModel {
	return &TestModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "test",
	}
}

func (m *TestModel) Insert(data Test) (sql.Result, error) {
	testNanosecondKey := fmt.Sprintf("%s%v", cacheTestNanosecondPrefix, data.Nanosecond)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `insert into ` + m.table + ` (` + testRowsExpectAutoSet + `) values (?, ?)`
		return conn.Exec(query, data.Nanosecond, data.Data)
	}, testNanosecondKey)
	return ret, err
}

func (m *TestModel) FindOne(id int64) (*Test, error) {
	testIdKey := fmt.Sprintf("%s%v", cacheTestIdPrefix, id)
	var resp Test
	err := m.QueryRow(&resp, testIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := `select ` + testRows + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *TestModel) FindOneByNanosecond(nanosecond int64) (*Test, error) {
	testNanosecondKey := fmt.Sprintf("%s%v", cacheTestNanosecondPrefix, nanosecond)
	var resp Test
	err := m.QueryRowIndex(&resp, testNanosecondKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := `select ` + testRows + ` from ` + m.table + ` where nanosecond = ? limit 1`
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

func (m *TestModel) Update(data Test) (sql.Result, error) {
	testIdKey := fmt.Sprintf("%s%v", cacheTestIdPrefix, data.Id)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + testRowsWithPlaceHolder + ` where id = ?`
		return conn.Exec(query, data.Nanosecond, data.Data, data.Id)
	}, testIdKey)
	return ret, err
}

func (m *TestModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	testNanosecondKey := fmt.Sprintf("%s%v", cacheTestNanosecondPrefix, data.Nanosecond)
	testIdKey := fmt.Sprintf("%s%v", cacheTestIdPrefix, id)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `delete from ` + m.table + ` where id = ?`
		return conn.Exec(query, id)
	}, testNanosecondKey, testIdKey)
	return err
}

func (m *TestModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheTestIdPrefix, primary)
}

func (m *TestModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := `select ` + testRows + ` from ` + m.table + ` where id = ? limit 1`
	return conn.QueryRow(v, query, primary)
}
