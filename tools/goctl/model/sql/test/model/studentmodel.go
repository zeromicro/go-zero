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
	studentFieldNames          = builderx.RawFieldNames(&Student{})
	studentRows                = strings.Join(studentFieldNames, ",")
	studentRowsExpectAutoSet   = strings.Join(stringx.Remove(studentFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	studentRowsWithPlaceHolder = strings.Join(stringx.Remove(studentFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheStudentIDPrefix = "cache#Student#id#"
)

type (
	// StudentModel defines a model for Student
	StudentModel interface {
		Insert(data Student) (sql.Result, error)
		FindOne(id int64) (*Student, error)
		Update(data Student) error
		Delete(id int64) error
	}

	defaultStudentModel struct {
		sqlc.CachedConn
		table string
	}

	// Student defines an data structure for mysql
	Student struct {
		ID         int64           `db:"id"`
		Name       string          `db:"name"`
		Age        sql.NullInt64   `db:"age"`
		Score      sql.NullFloat64 `db:"score"`
		CreateTime time.Time       `db:"create_time"`
		UpdateTime sql.NullTime    `db:"update_time"`
	}
)

// NewStudentModel creates an instance for StudentModel
func NewStudentModel(conn sqlx.SqlConn, c cache.CacheConf) StudentModel {
	return &defaultStudentModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`student`",
	}
}

func (m *defaultStudentModel) Insert(data Student) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?)", m.table, studentRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.Name, data.Age, data.Score)

	return ret, err
}

func (m *defaultStudentModel) FindOne(id int64) (*Student, error) {
	studentIDKey := fmt.Sprintf("%s%v", cacheStudentIDPrefix, id)
	var resp Student
	err := m.QueryRow(&resp, studentIDKey, func(conn sqlx.SqlConn, v interface{}) error {
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

func (m *defaultStudentModel) Update(data Student) error {
	studentIDKey := fmt.Sprintf("%s%v", cacheStudentIDPrefix, data.ID)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, studentRowsWithPlaceHolder)
		return conn.Exec(query, data.Name, data.Age, data.Score, data.ID)
	}, studentIDKey)
	return err
}

func (m *defaultStudentModel) Delete(id int64) error {

	studentIDKey := fmt.Sprintf("%s%v", cacheStudentIDPrefix, id)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, studentIDKey)
	return err
}

func (m *defaultStudentModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheStudentIDPrefix, primary)
}

func (m *defaultStudentModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", studentRows, m.table)
	return conn.QueryRow(v, query, primary)
}
