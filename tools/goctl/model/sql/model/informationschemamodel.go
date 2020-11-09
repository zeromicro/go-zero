package model

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type (
	InformationSchemaModel struct {
		conn sqlx.SqlConn
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
