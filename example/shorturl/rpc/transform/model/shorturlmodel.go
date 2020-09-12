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
	shorturlFieldNames          = builderx.FieldNames(&Shorturl{})
	shorturlRows                = strings.Join(shorturlFieldNames, ",")
	shorturlRowsExpectAutoSet   = strings.Join(stringx.Remove(shorturlFieldNames, "create_time", "update_time"), ",")
	shorturlRowsWithPlaceHolder = strings.Join(stringx.Remove(shorturlFieldNames, "shorten", "create_time", "update_time"), "=?,") + "=?"

	cacheShorturlShortenPrefix = "cache#Shorturl#shorten#"
)

type (
	ShorturlModel struct {
		sqlc.CachedConn
		table string
	}

	Shorturl struct {
		Shorten string `db:"shorten"` // shorten key
		Url     string `db:"url"`     // original url
	}
)

func NewShorturlModel(conn sqlx.SqlConn, c cache.CacheConf, table string) *ShorturlModel {
	return &ShorturlModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      table,
	}
}

func (m *ShorturlModel) Insert(data Shorturl) (sql.Result, error) {
	query := `insert into ` + m.table + ` (` + shorturlRowsExpectAutoSet + `) values (?, ?)`
	return m.ExecNoCache(query, data.Shorten, data.Url)
}

func (m *ShorturlModel) FindOne(shorten string) (*Shorturl, error) {
	shorturlShortenKey := fmt.Sprintf("%s%v", cacheShorturlShortenPrefix, shorten)
	var resp Shorturl
	err := m.QueryRow(&resp, shorturlShortenKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := `select ` + shorturlRows + ` from ` + m.table + ` where shorten = ? limit 1`
		return conn.QueryRow(v, query, shorten)
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

func (m *ShorturlModel) Update(data Shorturl) error {
	shorturlShortenKey := fmt.Sprintf("%s%v", cacheShorturlShortenPrefix, data.Shorten)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `update ` + m.table + ` set ` + shorturlRowsWithPlaceHolder + ` where shorten = ?`
		return conn.Exec(query, data.Url, data.Shorten)
	}, shorturlShortenKey)
	return err
}

func (m *ShorturlModel) Delete(shorten string) error {
	_, err := m.FindOne(shorten)
	if err != nil {
		return err
	}

	shorturlShortenKey := fmt.Sprintf("%s%v", cacheShorturlShortenPrefix, shorten)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := `delete from ` + m.table + ` where shorten = ?`
		return conn.Exec(query, shorten)
	}, shorturlShortenKey)
	return err
}
