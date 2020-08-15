package model

import (
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
	query := `select ` + userCamelRows + ` from ` + m.table + ` where id = ? limit 1`
	var resp UserCamel
	err := m.QueryRowNoCache(&resp, query, id)
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
	var resp UserCamel
	query := `select ` + userCamelRows + ` from ` + m.table + ` where name limit 1`
	err := m.QueryRowNoCache(&resp, query, name)
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
	var resp UserCamel
	query := `select ` + userCamelRows + ` from ` + m.table + ` where mobile limit 1`
	err := m.QueryRowNoCache(&resp, query, mobile)
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
	query := `update ` + m.table + ` set ` + userCamelRowsWithPlaceHolder + ` where id = ?`
	_, err := m.ExecNoCache(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, data.Id)
	return err
}

func (m *UserCamelModel) Delete(id int64) error {
	query := `delete from ` + m.table + ` where id = ?`
	_, err := m.ExecNoCache(query, id)
	return err
}
