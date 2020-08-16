package model

import (
	"database/sql"
	"strings"
	"time"

	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	userCourseFieldNames          = builderx.FieldNames(&UserCourse{})
	userCourseRows                = strings.Join(userCourseFieldNames, ",")
	userCourseRowsExpectAutoSet   = strings.Join(stringx.Remove(userCourseFieldNames, "id", "create_time", "update_time"), ",")
	userCourseRowsWithPlaceHolder = strings.Join(stringx.Remove(userCourseFieldNames, "id", "create_time", "update_time"), "=?,") + "=?"
)

type (
	UserCourseModel struct {
		conn  sqlx.SqlConn
		table string
	}

	UserCourse struct {
		Id         int64     `db:"id"`
		UserId     int64     `db:"user_id"`     // 用户id
		CourseName string    `db:"course_name"` // 课程名称
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
	}
)

func NewUserCourseModel(conn sqlx.SqlConn, table string) *UserCourseModel {
	return &UserCourseModel{
		conn:  conn,
		table: table,
	}
}

func (m *UserCourseModel) Insert(data UserCourse) (sql.Result, error) {
	query := `insert into ` + m.table + `(` + userCourseRowsExpectAutoSet + `) value (?, ?)`
	return m.conn.Exec(query, data.UserId, data.CourseName)
}

func (m *UserCourseModel) FindOne(id int64) (*UserCourse, error) {
	query := `select ` + userCourseRows + ` from ` + m.table + ` where id = ? limit 1`
	var resp UserCourse
	err := m.conn.QueryRow(&resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *UserCourseModel) FindOneByUserId(userId int64) (*UserCourse, error) {
	var resp UserCourse
	query := `select ` + userCourseRows + ` from ` + m.table + ` where user_id limit 1`
	err := m.conn.QueryRow(&resp, query, userId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *UserCourseModel) FindOneByCourseName(courseName string) (*UserCourse, error) {
	var resp UserCourse
	query := `select ` + userCourseRows + ` from ` + m.table + ` where course_name limit 1`
	err := m.conn.QueryRow(&resp, query, courseName)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *UserCourseModel) Update(data UserCourse) error {
	query := `update ` + m.table + ` set ` + userCourseRowsWithPlaceHolder + ` where id = ?`
	_, err := m.conn.Exec(query, data.UserId, data.CourseName, data.Id)
	return err
}

func (m *UserCourseModel) Delete(id int64) error {
	query := `delete from ` + m.table + ` where id = ?`
	_, err := m.conn.Exec(query, id)
	return err
}
