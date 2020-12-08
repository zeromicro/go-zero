package model

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/core/stores/cache"
	"github.com/tal-tech/go-zero/core/stores/sqlc"
	"github.com/tal-tech/go-zero/core/stores/sqlx"
	"github.com/tal-tech/go-zero/core/stringx"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/builderx"
)

var (
	bookFieldNames          = builderx.FieldNames(&Book{})
	bookRows                = strings.Join(bookFieldNames, ",")
	bookRowsExpectAutoSet   = strings.Join(stringx.Remove(bookFieldNames, "create_time", "update_time"), ",")
	bookRowsWithPlaceHolder = strings.Join(stringx.Remove(bookFieldNames, "book", "create_time", "update_time"), "=?,") + "=?"

	cacheBookPrefix = "cache#Book#book#"
)

type (
	BookModel interface {
		Insert(data Book) (sql.Result, error)
		FindOne(book string) (*Book, error)
		Update(data Book) error
		Delete(book string) error
	}

	defaultBookModel struct {
		sqlc.CachedConn
		table string
	}

	Book struct {
		Book  string `db:"book"`  // book name
		Price int64  `db:"price"` // book price
	}
)

func NewBookModel(conn sqlx.SqlConn, c cache.CacheConf) BookModel {
	return &defaultBookModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "book",
	}
}

func (m *defaultBookModel) Insert(data Book) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?)", m.table, bookRowsExpectAutoSet)
	ret, err := m.ExecNoCache(query, data.Book, data.Price)

	return ret, err
}

func (m *defaultBookModel) FindOne(book string) (*Book, error) {
	bookKey := fmt.Sprintf("%s%v", cacheBookPrefix, book)
	var resp Book
	err := m.QueryRow(&resp, bookKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where book = ? limit 1", bookRows, m.table)
		return conn.QueryRow(v, query, book)
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

func (m *defaultBookModel) Update(data Book) error {
	bookKey := fmt.Sprintf("%s%v", cacheBookPrefix, data.Book)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where book = ?", m.table, bookRowsWithPlaceHolder)
		return conn.Exec(query, data.Price, data.Book)
	}, bookKey)
	return err
}

func (m *defaultBookModel) Delete(book string) error {

	bookKey := fmt.Sprintf("%s%v", cacheBookPrefix, book)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where book = ?", m.table)
		return conn.Exec(query, book)
	}, bookKey)
	return err
}

func (m *defaultBookModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheBookPrefix, primary)
}

func (m *defaultBookModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where book = ? limit 1", bookRows, m.table)
	return conn.QueryRow(v, query, primary)
}
