// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	schoolFieldNames          = builder.RawFieldNames(&School{})
	schoolRows                = strings.Join(schoolFieldNames, ",")
	schoolRowsExpectAutoSet   = strings.Join(stringx.Remove(schoolFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	schoolRowsWithPlaceHolder = strings.Join(stringx.Remove(schoolFieldNames, "`id`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	schoolModel interface {
		Insert(ctx context.Context, data *School) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*School, error)
		Update(ctx context.Context, data *School) error
		Delete(ctx context.Context, id int64) error
	}

	defaultSchoolModel struct {
		conn  sqlx.SqlConn
		table string
	}

	School struct {
		Id       int64          `db:"id"`
		Name     sql.NullString `db:"name"`    // The username
		UserId   int64          `db:"user_id"` // The user id
		Type     int64          `db:"type"`    // The user type, 0:normal,1:vip, for test golang keyword
		CreateAt sql.NullTime   `db:"create_at"`
		UpdateAt time.Time      `db:"update_at"`
	}
)

func newSchoolModel(conn sqlx.SqlConn) *defaultSchoolModel {
	return &defaultSchoolModel{
		conn:  conn,
		table: "`school`",
	}
}

func (m *defaultSchoolModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *defaultSchoolModel) FindOne(ctx context.Context, id int64) (*School, error) {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", schoolRows, m.table)
	var resp School
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultSchoolModel) Insert(ctx context.Context, data *School) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, schoolRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Name, data.UserId, data.Type)
	return ret, err
}

func (m *defaultSchoolModel) Update(ctx context.Context, data *School) error {
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, schoolRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Name, data.UserId, data.Type, data.Id)
	return err
}

func (m *defaultSchoolModel) tableName() string {
	return m.table
}
