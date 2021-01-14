package model

import "github.com/tal-tech/go-zero/core/stores/sqlx"

type (
	InformationSchemaModel struct {
		conn sqlx.SqlConn
	}

	Column struct {
		Name          string      `db:"COLUMN_NAME"`
		DataType      string      `db:"DATA_TYPE"`
		Key           string      `db:"COLUMN_KEY"`
		Extra         string      `db:"EXTRA"`
		Comment       string      `db:"COLUMN_COMMENT"`
		ColumnDefault interface{} `db:"COLUMN_DEFAULT"`
		IsNullAble    string      `db:"IS_NULLABLE"`
	}
)

func NewInformationSchemaModel(conn sqlx.SqlConn) *InformationSchemaModel {
	return &InformationSchemaModel{conn: conn}
}

func (m *InformationSchemaModel) GetAllTables(database string) ([]string, error) {
	query := `select TABLE_NAME from TABLES where TABLE_SCHEMA = ?`
	var tables []string
	err := m.conn.QueryRows(&tables, query, database)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func (m *InformationSchemaModel) FindByTableName(db, table string) ([]*Column, error) {
	querySql := `select COLUMN_NAME,COLUMN_DEFAULT,IS_NULLABLE,DATA_TYPE,COLUMN_KEY,EXTRA,COLUMN_COMMENT from COLUMNS where TABLE_SCHEMA = ? and TABLE_NAME = ?`
	var reply []*Column
	err := m.conn.QueryRows(&reply, querySql, db, table)
	return reply, err
}
