package model

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
	userCamelFieldNames          = builderx.FieldNames(&UserCamel{})
	userCamelRows                = strings.Join(userCamelFieldNames, ",")
	userCamelRowsExpectAutoSet   = strings.Join(stringx.Remove(userCamelFieldNames, "create_time", "update_time"), ",")
	userCamelRowsWithPlaceHolder = strings.Join(stringx.Remove(userCamelFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"

	cacheuserCamelIdPrefix     = "cache#userCamel#id#"
	cacheuserCamelNamePrefix   = "cache#userCamel#name#"
	cacheuserCamelMobilePrefix = "cache#userCamel#mobile#"

	ErrNotFound = sqlx.ErrNotFound
)

type (
	UserCamelModel struct {
		sqlc.CachedConn
		table string
	}

	UserCamel struct {
		Id         int64     `db:"id"`
		Name       string    `db:"name"`     // 用户名称
		Password   string    `db:"password"` // 用户密码
		Mobile     string    `db:"mobile"`   // 手机号
		Gender     string    `db:"gender"`   // 男｜女｜未公开
		Nickname   string    `db:"nickname"` // 用户昵称
		CreateTime time.Time `db:"createTime"`
		UpdateTime time.Time `db:"updateTime"`
	}
)

func NewUserCamelModel(conn sqlx.SqlConn, c cache.CacheConf, table string) *UserCamelModel {
	return &UserCamelModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      table,
	}
}

func (m *UserCamelModel) Insert(data UserCamel) error {
	query := `insert into ` + m.table + `(` + userCamelRowsExpectAutoSet + `) value (?, ?, ?, ?, ?, ?)`
	_, err := m.ExecNoCache(query, data.Id, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname)
	return err
}

func (m *UserCamelModel) FindOne(id int64) (*UserCamel, error) {
	userCamelIdKey := fmt.Sprintf("cache#userCamel#id#%v", id)
	var resp UserCamel
	err := m.QueryRow(&resp, userCamelIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := `select ` + userCamelRows + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *UserCamelModel) FindOneByName(name string) (*UserCamel, error) {
	userCamelNameKey := fmt.Sprintf("cache#userCamel#name#%v", name)
	var resp UserCamel
	err := m.QueryRowIndex(&resp, userCamelNameKey, func(primary interface{}) string {
		return fmt.Sprintf("%s%v", cacheuserCamelIdPrefix, primary)
	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := `select ` + userCamelRows + ` from ` + m.table + ` where name = ? limit 1`
		if err := conn.QueryRow(&resp, query, name); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
		query := `select ` + userCamelRows + ` from ` + m.table + ` where id = ? limit 1`
		return conn.QueryRow(v, query, primary)
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

func (m *UserCamelModel) FindOneByMobile(mobile string) (*UserCamel, error) {
	userCamelMobileKey := fmt.Sprintf("cache#userCamel#mobile#%v", mobile)
	var resp UserCamel
	err := m.QueryRowIndex(&resp, userCamelMobileKey, func(primary interface{}) string {
		return fmt.Sprintf("%s%v", cacheuserCamelIdPrefix, primary)
	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := `select ` + userCamelRows + ` from ` + m.table + ` where mobile = ? limit 1`
		if err := conn.QueryRow(&resp, query, mobile); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
		query := `select ` + userCamelRows + ` from ` + m.table + ` where id = ? limit 1`
		return conn.QueryRow(v, query, primary)
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

func (m *UserCamelModel) Update(data UserCamel) error {
	userCamelIdKey := fmt.Sprintf("cache#userCamel#id#%v", data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + userCamelRowsWithPlaceHolder + ` where id = ?`
		return conn.Exec(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, data.Id)
	}, userCamelIdKey)
	return err
}

func (m *UserCamelModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	userCamelIdKey := fmt.Sprintf("cache#userCamel#id#%v", id)
	userCamelNameKey := fmt.Sprintf("cache#userCamel#name#%v", data.Name)
	userCamelMobileKey := fmt.Sprintf("cache#userCamel#mobile#%v", data.Mobile)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `delete from ` + m.table + ` where id = ?`
		return conn.Exec(query, id)
	}, userCamelIdKey, userCamelNameKey, userCamelMobileKey)
	return err
}
