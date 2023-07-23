package model

import (
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
	studentFieldNames          = builder.RawFieldNames(&Student{})
	studentRows                = strings.Join(studentFieldNames, ",")
	studentRowsExpectAutoSet   = strings.Join(stringx.Remove(studentFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	studentRowsWithPlaceHolder = strings.Join(stringx.Remove(studentFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheStudentIdPrefix        = "cache#student#id#"
	cacheStudentClassNamePrefix = "cache#student#class#name#"
)

type (
	// StudentModel only for test
	StudentModel interface {
		Insert(data Student) (sql.Result, error)
		FindOne(id int64) (*Student, error)
		FindOneByClassName(class, name string) (*Student, error)
		Update(data Student) error
		// only for test
		Delete(id int64, className, studentName string) error
	}

	defaultStudentModel struct {
		sqlc.CachedConn
		table string
	}

	// Student only for test
	Student struct {
		Id         int64           `db:"id"`
		Class      string          `db:"class"`
		Name       string          `db:"name"`
		Age        sql.NullInt64   `db:"age"`
		Score      sql.NullFloat64 `db:"score"`
		CreateTime time.Time       `db:"create_time"`
		UpdateTime sql.NullTime    `db:"update_time"`
	}
)

// NewStudentModel only for test
func NewStudentModel(conn sqlx.SqlConn, c cache.CacheConf) StudentModel {
	return &defaultStudentModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`student`",
	}
}

func (m *defaultStudentModel) Insert(data Student) (sql.Result, error) {
	studentClassNameKey := fmt.Sprintf("%s%v%v", cacheStudentClassNamePrefix, data.Class, data.Name)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?)", m.table, studentRowsExpectAutoSet)
		return conn.Exec(query, data.Class, data.Name, data.Age, data.Score)
	}, studentClassNameKey)
	return ret, err
}

func (m *defaultStudentModel) FindOne(id int64) (*Student, error) {
	studentIdKey := fmt.Sprintf("%s%v", cacheStudentIdPrefix, id)
	var resp Student
	err := m.QueryRow(&resp, studentIdKey, func(conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", studentRows, m.table)
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

func (m *defaultStudentModel) FindOneByClassName(class, name string) (*Student, error) {
	studentClassNameKey := fmt.Sprintf("%s%v%v", cacheStudentClassNamePrefix, class, name)
	var resp Student
	err := m.QueryRowIndex(&resp, studentClassNameKey, m.formatPrimary, func(conn sqlx.SqlConn, v any) (i any, e error) {
		query := fmt.Sprintf("select %s from %s where `class` = ? and `name` = ? limit 1", studentRows, m.table)
		if err := conn.QueryRow(&resp, query, class, name); err != nil {
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

func (m *defaultStudentModel) Update(data Student) error {
	studentIdKey := fmt.Sprintf("%s%v", cacheStudentIdPrefix, data.Id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, studentRowsWithPlaceHolder)
		return conn.Exec(query, data.Class, data.Name, data.Age, data.Score, data.Id)
	}, studentIdKey)
	return err
}

func (m *defaultStudentModel) Delete(id int64, className, studentName string) error {
	studentIdKey := fmt.Sprintf("%s%v", cacheStudentIdPrefix, id)
	studentClassNameKey := fmt.Sprintf("%s%v%v", cacheStudentClassNamePrefix, className, studentName)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, studentIdKey, studentClassNameKey)
	return err
}

func (m *defaultStudentModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheStudentIdPrefix, primary)
}

func (m *defaultStudentModel) queryPrimary(conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", studentRows, m.table)
	return conn.QueryRow(v, query, primary)
}
