// Code generated by goctl. DO NOT EDIT!

package cache

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userFieldNames          = builder.RawFieldNames(&User{})
	userRows                = strings.Join(userFieldNames, ",")
	userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), ",")
	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "`id`", "`create_time`", "`update_time`", "`create_at`", "`update_at`"), "=?,") + "=?"

	cacheUserIdPrefix     = "cache:user:id:"
	cacheUserMobilePrefix = "cache:user:mobile:"
	cacheUserNamePrefix   = "cache:user:name:"
	cacheUserTypePrefix   = "cache:user:type:"
	cacheUserUserPrefix   = "cache:user:user:"
)

type (
	userModel interface {
		Insert(ctx context.Context, data *User) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*User, error)
		FindOneByMobile(ctx context.Context, mobile string) (*User, error)
		FindOneByName(ctx context.Context, name sql.NullString) (*User, error)
		FindOneByType(ctx context.Context, tp int64) (*User, error)
		FindOneByUser(ctx context.Context, user string) (*User, error)
		Update(ctx context.Context, data *User) error
		Delete(ctx context.Context, id int64) error
	}

	defaultUserModel struct {
		sqlc.CachedConn
		table string
	}

	User struct {
		Id         int64          `db:"id"`
		User       string         `db:"user"`     // 用户
		Name       sql.NullString `db:"name"`     // 用户\t名称
		Age        int64          `db:"age"`      // 年龄
		Password   string         `db:"password"` // 用户密码
		Mobile     string         `db:"mobile"`   // 手机号
		Gender     string         `db:"gender"`   // 男｜女｜未公开
		Nickname   string         `db:"nickname"` // 用户昵称
		Type       int64          `db:"type"`     // 用户类型
		CreateTime sql.NullTime   `db:"create_time"`
		UpdateTime time.Time      `db:"update_time"`
	}
)

func newUserModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUserModel {
	return &defaultUserModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user`",
	}
}

func (m *defaultUserModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	userTypeKey := fmt.Sprintf("%s%v", cacheUserTypePrefix, data.Type)
	userUserKey := fmt.Sprintf("%s%v", cacheUserUserPrefix, data.User)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, userIdKey, userMobileKey, userNameKey, userTypeKey, userUserKey)
	return err
}

func (m *defaultUserModel) FindOne(ctx context.Context, id int64) (*User, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, id)
	var resp User
	err := m.QueryRowCtx(ctx, &resp, userIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
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

func (m *defaultUserModel) FindOneByMobile(ctx context.Context, mobile string) (*User, error) {
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, mobile)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userMobileKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `mobile` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, mobile); err != nil {
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

func (m *defaultUserModel) FindOneByName(ctx context.Context, name sql.NullString) (*User, error) {
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, name)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userNameKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `name` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, name); err != nil {
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

func (m *defaultUserModel) FindOneByType(ctx context.Context, tp int64) (*User, error) {
	userTypeKey := fmt.Sprintf("%s%v", cacheUserTypePrefix, tp)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userTypeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `type` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, tp); err != nil {
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

func (m *defaultUserModel) FindOneByUser(ctx context.Context, user string) (*User, error) {
	userUserKey := fmt.Sprintf("%s%v", cacheUserUserPrefix, user)
	var resp User
	err := m.QueryRowIndexCtx(ctx, &resp, userUserKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user` = ? limit 1", userRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, user); err != nil {
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

func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	userTypeKey := fmt.Sprintf("%s%v", cacheUserTypePrefix, data.Type)
	userUserKey := fmt.Sprintf("%s%v", cacheUserUserPrefix, data.User)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, data.User, data.Name, data.Age, data.Password, data.Mobile, data.Gender, data.Nickname, data.Type)
	}, userIdKey, userMobileKey, userNameKey, userTypeKey, userUserKey)
	return ret, err
}

func (m *defaultUserModel) Update(ctx context.Context, newData *User) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	userIdKey := fmt.Sprintf("%s%v", cacheUserIdPrefix, data.Id)
	userMobileKey := fmt.Sprintf("%s%v", cacheUserMobilePrefix, data.Mobile)
	userNameKey := fmt.Sprintf("%s%v", cacheUserNamePrefix, data.Name)
	userTypeKey := fmt.Sprintf("%s%v", cacheUserTypePrefix, data.Type)
	userUserKey := fmt.Sprintf("%s%v", cacheUserUserPrefix, data.User)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query, newData.User, newData.Name, newData.Age, newData.Password, newData.Mobile, newData.Gender, newData.Nickname, newData.Type, newData.Id)
	}, userIdKey, userMobileKey, userNameKey, userTypeKey, userUserKey)
	return err
}

func (m *defaultUserModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheUserIdPrefix, primary)
}

func (m *defaultUserModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultUserModel) tableName() string {
	return m.table
}
