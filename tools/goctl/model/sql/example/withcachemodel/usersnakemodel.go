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
	userSnakeFieldNames          = builderx.FieldNames(&UserSnake{})
	userSnakeRows                = strings.Join(userSnakeFieldNames, ",")
	userSnakeRowsExpectAutoSet   = strings.Join(stringx.Remove(userSnakeFieldNames, "id", "create_time", "update_time"), ",")
	userSnakeRowsWithPlaceHolder = strings.Join(stringx.Remove(userSnakeFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"

	cacheUserSnakeIdPrefix     = "cache#UserSnake#id#"
	cacheUserSnakeNamePrefix   = "cache#UserSnake#name#"
	cacheUserSnakeMobilePrefix = "cache#UserSnake#mobile#"
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
	userSnakeIdKey := fmt.Sprintf("%s%v", cacheUserSnakeIdPrefix, id)
	var resp UserSnake
	err := m.QueryRow(&resp, userSnakeIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := `select ` + userSnakeRows + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *UserSnakeModel) FindOneByName(name string) (*UserSnake, error) {
	userSnakeNameKey := fmt.Sprintf("%s%v", cacheUserSnakeNamePrefix, name)
	var resp UserSnake
	err := m.QueryRowIndex(&resp, userSnakeNameKey, func(primary interface{}) string {
		return fmt.Sprintf("%s%v", cacheUserSnakeIdPrefix, primary)
	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := `select ` + userSnakeRows + ` from ` + m.table + ` where name = ? limit 1`
		if err := conn.QueryRow(&resp, query, name); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
		query := `select ` + userSnakeRows + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *UserSnakeModel) FindOneByMobile(mobile string) (*UserSnake, error) {
	userSnakeMobileKey := fmt.Sprintf("%s%v", cacheUserSnakeMobilePrefix, mobile)
	var resp UserSnake
	err := m.QueryRowIndex(&resp, userSnakeMobileKey, func(primary interface{}) string {
		return fmt.Sprintf("%s%v", cacheUserSnakeIdPrefix, primary)
	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := `select ` + userSnakeRows + ` from ` + m.table + ` where mobile = ? limit 1`
		if err := conn.QueryRow(&resp, query, mobile); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
		query := `select ` + userSnakeRows + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *UserSnakeModel) Update(data UserSnake) error {
	userSnakeIdKey := fmt.Sprintf("%s%v", cacheUserSnakeIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + userSnakeRowsWithPlaceHolder + ` where id = ?`
		return conn.Exec(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, data.Id)
	}, userSnakeIdKey)
	return err
}

func (m *UserSnakeModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	userSnakeIdKey := fmt.Sprintf("%s%v", cacheUserSnakeIdPrefix, id)
	userSnakeNameKey := fmt.Sprintf("%s%v", cacheUserSnakeNamePrefix, data.Name)
	userSnakeMobileKey := fmt.Sprintf("%s%v", cacheUserSnakeMobilePrefix, data.Mobile)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `delete from ` + m.table + ` where id = ?`
		return conn.Exec(query, id)
	}, userSnakeMobileKey, userSnakeIdKey, userSnakeNameKey)
	return err
}
