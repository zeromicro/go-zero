package main

import (
	"database/sql"
	"fmt"

	"zero/core/stores/cache"
	"zero/core/stores/sqlc"
	"zero/core/stores/sqlx"
	"zero/kq"
)

var (
	userRows = "id, mobile, name, sex"

	cacheUserMobilePrefix = "cache#user#mobile#"
	cacheUserIdPrefix     = "cache#user#id#"

	ErrNotFound = sqlc.ErrNotFound
)

type (
	User struct {
		Id     int64  `db:"id" json:"id,omitempty"`
		Mobile string `db:"mobile" json:"mobile,omitempty"`
		Name   string `db:"name" json:"name,omitempty"`
		Sex    int    `db:"sex" json:"sex,omitempty"`
	}

	UserModel struct {
		sqlc.CachedConn
		// sqlx.SqlConn
		table string

		// kafka use kq not kmq
		push *kq.Pusher
	}
)

func NewUserModel(db sqlx.SqlConn, c cache.CacheConf, table string, pusher *kq.Pusher) *UserModel {
	return &UserModel{
		CachedConn: sqlc.NewConn(db, c),
		table:      table,
		push:       pusher,
	}
}

func (um *UserModel) FindOne(id int64) (*User, error) {
	key := fmt.Sprintf("%s%d", cacheUserIdPrefix, id)
	var user User
	err := um.QueryRow(&user, key, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("SELECT %s FROM user WHERE id=?", userRows)
		return conn.QueryRow(v, query, id)
	})
	switch err {
	case nil:
		return &user, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (um *UserModel) FindByMobile(mobile string) (*User, error) {
	var user User
	key := fmt.Sprintf("%s%s", cacheUserMobilePrefix, mobile)
	err := um.QueryRowIndex(&user, key, func(primary interface{}) string {
		return fmt.Sprintf("%s%d", cacheUserIdPrefix, primary.(int64))
	}, func(conn sqlx.SqlConn, v interface{}) (interface{}, error) {
		query := fmt.Sprintf("SELECT %s FROM user WHERE mobile=?", userRows)
		if err := conn.QueryRow(&user, query, mobile); err != nil {
			return nil, err
		}
		return user.Id, nil
	}, func(conn sqlx.SqlConn, v interface{}, primary interface{}) error {
		return conn.QueryRow(v, "SELECT * FROM user WHERE id=?", primary)
	})
	switch err {
	case nil:
		return &user, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Count for no cache
func (um *UserModel) Count() (int64, error) {
	var count int64
	err := um.QueryRowNoCache(&count, "SELECT count(1) FROM user")
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Query rows
func (um *UserModel) FindByName(name string) ([]*User, error) {
	var users []*User
	query := fmt.Sprintf("SELECT %s FROM user WHERE name=?", userRows)
	err := um.QueryRowsNoCache(&userRows, query, name)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (um *UserModel) UpdateSexById(sex int, id int64) error {
	key := fmt.Sprintf("%s%d", cacheUserIdPrefix, id)
	_, err := um.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("UPDATE user SET sex=? WHERE id=?")
		return conn.Exec(query, sex, id)
	}, key)
	return err
}

func (um *UserModel) UpdateMobileById(mobile string, id int64) error {
	idKey := fmt.Sprintf("%s%d", cacheUserIdPrefix, id)
	mobileKey := fmt.Sprintf("%s%s", cacheUserMobilePrefix, mobile)
	_, err := um.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("UPDATE user SET mobile=? WHERE id=?")
		return conn.Exec(query, mobile, id)
	}, idKey, mobileKey)
	return err
}

func (um *UserModel) Update(u *User) error {
	oldUser, err := um.FindOne(u.Id)
	if err != nil {
		return err
	}

	idKey := fmt.Sprintf("%s%d", cacheUserIdPrefix, oldUser.Id)
	mobileKey := fmt.Sprintf("%s%s", cacheUserMobilePrefix, oldUser.Mobile)
	_, err = um.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("UPDATE user SET mobile=?, name=?, sex=? WHERE id=?")
		return conn.Exec(query, u.Mobile, u.Name, u.Sex, u.Id)
	}, idKey, mobileKey)
	return err
}
