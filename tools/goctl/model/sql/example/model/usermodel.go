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
	userFieldNames          = builderx.FieldNames(&User{})
	userRows                = strings.Join(userFieldNames, ",")
	userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "id", "create_time", "update_time"), ",")
	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"

	cacheUserMobilePrefix = "cache#User#mobile#"
	cacheUserIdPrefix     = "cache#User#id#"
	cacheUserNamePrefix   = "cache#User#name#"
)

type (
	UserModel struct {
		sqlc.CachedConn
		table string
	}

	User struct {
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

func NewUserModel(conn sqlx.SqlConn, c cache.CacheConf, table string) *UserModel {
	return &UserModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      table,
	}
}

func (m *UserModel) Insert(data User) (sql.Result, error) {
	query := `insert into ` + m.table + `(` + userRowsExpectAutoSet + `) value (?, ?, ?, ?, ?)`
	return m.ExecNoCache(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname)
}

func (m *UserModel) FindOne(id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRow(&resp, userIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := `select ` + userRows + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *UserModel) FindOneByName(name string) (*User, error) {
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, name)
	var resp User
	err := m.QueryRowIndex(&resp, userNameKey, func(primary interface{}) string {
		return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := `select ` + userRows + ` from ` + m.table + ` where name = ? limit 1`
		if err := conn.QueryRow(&resp, query, name); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
		query := `select ` + userRows + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *UserModel) FindOneByMobile(mobile string) (*User, error) {
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, mobile)
	var resp User
	err := m.QueryRowIndex(&resp, userMobileKey, func(primary interface{}) string {
		return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
	}, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := `select ` + userRows + ` from ` + m.table + ` where mobile = ? limit 1`
		if err := conn.QueryRow(&resp, query, mobile); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, func(conn sqlx.SqlConn, v, primary interface{}) error {
		query := `select ` + userRows + ` from ` + m.table + ` where id = ? limit 1`
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

func (m *UserModel) Update(data User) error {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + userRowsWithPlaceHolder + ` where id = ?`
		return conn.Exec(query, data.Name, data.Password, data.Mobile, data.Gender, data.Nickname, data.Id)
	}, userIdKey)
	return err
}

func (m *UserModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `delete from ` + m.table + ` where id = ?`
		return conn.Exec(query, id)
	}, userIdKey, userNameKey, userMobileKey)
	return err
}
