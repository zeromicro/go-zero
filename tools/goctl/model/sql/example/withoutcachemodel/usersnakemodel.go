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
	userSnakeFieldNames          = builderx.FieldNames(&UserSnake{})
	userSnakeRows                = strings.Join(userSnakeFieldNames, ",")
	userSnakeRowsExpectAutoSet   = strings.Join(stringx.Remove(userSnakeFieldNames, "id", "create_time", "update_time"), ",")
	userSnakeRowsWithPlaceHolder = strings.Join(stringx.Remove(userSnakeFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
)

type (
	UserSnakeModel struct {
		sqlc.CachedConn
		table string
	}

	UserSnake struct {
		Id         int64     `db:"id"`
		Name       string    `db:"name"`     // 用户名称
		Password   string    `db:"password"` // 用户密码
		Mobile     string    `db:"mobile"`   // 手机号
		Gender     string    `db:"gender"`   // 男｜女｜未公开
		Nickname   string    `db:"nickname"` // 用户昵称
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewUserSnakeModel(conn sqlx.SqlConn, c cache.CacheConf, table string) *UserSnakeModel {
	return &UserSnakeModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      table,
	}
}

func (m *UserSnakeModel) Insert(data UserSnake) error {
	query := `insert into ` + m.table + `(` + userSnakeRowsExpectAutoSet + `) value (?, ?, ?, ?, ?)`
	_, err := m.ExecNoCache(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname)
	return err
}

func (m *UserSnakeModel) FindOne(id int64) (*UserSnake, error) {
	query := `select ` + userSnakeRows + ` from ` + m.table + ` where id = ? limit 1`
	var resp UserSnake
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

func (m *UserSnakeModel) FindOneByName(name string) (*UserSnake, error) {
	var resp UserSnake
	query := `select ` + userSnakeRows + ` from ` + m.table + ` where name limit 1`
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

func (m *UserSnakeModel) FindOneByMobile(mobile string) (*UserSnake, error) {
	var resp UserSnake
	query := `select ` + userSnakeRows + ` from ` + m.table + ` where mobile limit 1`
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

func (m *UserSnakeModel) Update(data UserSnake) error {
	query := `update ` + m.table + ` set ` + userSnakeRowsWithPlaceHolder + ` where id = ?`
	_, err := m.ExecNoCache(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, data.Id)
	return err
}

func (m *UserSnakeModel) Delete(id int64) error {
	query := `delete from ` + m.table + ` where id = ?`
	_, err := m.ExecNoCache(query, id)
	return err
}
