package model

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type (
	ColumnModel struct {
		conn sqlx.SqlConn
	}
	Column struct {
		Name     string `db:"COLUMN_NAME"`
		DataType string `db:"DATA_TYPE"`
		Key      string `db:"COLUMN_KEY"`
		Extra    string `db:"EXTRA"`
		Comment  string `db:"COLUMN_COMMENT"`
	}
)

func NewColumnModel(conn sqlx.SqlConn) *ColumnModel {
	return &ColumnModel{
		conn: conn,
	}
}

func (m *ColumnModel) FindByTableName(table string) ([]*Column, error) {
	querySql := `select COLUMN_NAME,DATA_TYPE,COLUMN_KEY,EXTRA,COLUMN_COMMENT from information_schema where TABLE_NAME = ?`
	var reply []*Column
	err := m.conn.QueryRows(&reply, querySql, table)
	return reply, err
}
