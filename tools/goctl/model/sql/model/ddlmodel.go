package model

import (
	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type (
	DDLModel struct {
		conn sqlx.SqlConn
	}
	DDL struct {
		Table string `db:"Table"`
		DDL   string `db:"Create Table"`
	}
)

func NewDDLModel(conn sqlx.SqlConn) *DDLModel {
	return &DDLModel{conn: conn}
}

func (m *DDLModel) ShowDDL(table ...string) ([]string, error) {
	var ddl []string
	for _, t := range table {
		query := `show create table ` + t
		var resp DDL
		err := m.conn.QueryRow(&resp, query)
		if err != nil {
			return nil, err
		}
		ddl = append(ddl, resp.DDL)
	}
	return ddl, nil
}
