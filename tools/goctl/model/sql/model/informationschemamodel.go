package model

import (
	"fmt"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

const indexPri = "PRIMARY"

type (
	// InformationSchemaModel defines information schema model
	InformationSchemaModel struct {
		conn sqlx.SqlConn
	}

	// Column defines column in table
	Column struct {
		Name          string      `db:"COLUMN_NAME"`
		DataType      string      `db:"DATA_TYPE"`
		Extra         string      `db:"EXTRA"`
		Comment       string      `db:"COLUMN_COMMENT"`
		ColumnDefault interface{} `db:"COLUMN_DEFAULT"`
		IsNullAble    string      `db:"IS_NULLABLE"`
		IndexName     string      `db:"INDEX_NAME"`
		NonUnique     int         `db:"NON_UNIQUE"`
		SeqInIndex    int         `db:"SEQ_IN_INDEX"`
	}

	ColumnData struct {
		Db      string
		Table   string
		Columns []*Column
	}

	Table struct {
		Db      string
		Table   string
		Columns []*Column
		// Primary key not included
		UniqueIndex map[string][]*Column
		PrimaryKey  *Column
		NormalIndex map[string][]*Column
	}

	IndexType string
	Index     struct {
		IndexType IndexType
		Columns   []*Column
	}
)

// NewInformationSchemaModel creates an instance for InformationSchemaModel
func NewInformationSchemaModel(conn sqlx.SqlConn) *InformationSchemaModel {
	return &InformationSchemaModel{conn: conn}
}

// GetAllTables selects all tables from TABLE_SCHEMA
func (m *InformationSchemaModel) GetAllTables(database string) ([]string, error) {
	query := `select TABLE_NAME from TABLES where TABLE_SCHEMA = ?`
	var tables []string
	err := m.conn.QueryRows(&tables, query, database)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func (m *InformationSchemaModel) FindColumns(db, table string) (*ColumnData, error) {
	querySql := `SELECT c.COLUMN_NAME,c.DATA_TYPE,EXTRA,c.COLUMN_COMMENT,c.COLUMN_DEFAULT,c.IS_NULLABLE,s.INDEX_NAME,s.NON_UNIQUE,s.SEQ_IN_INDEX from COLUMNS c LEFT JOIN STATISTICS s  ON  c.COLUMN_NAME=s.COLUMN_NAME WHERE  c.TABLE_SCHEMA = ? and c.TABLE_NAME = ? AND s.TABLE_SCHEMA = c.TABLE_SCHEMA and s.TABLE_NAME = c.TABLE_NAME`
	var reply []*Column
	err := m.conn.QueryRowsPartial(&reply, querySql, db, table)
	if err != nil {
		return nil, err
	}

	return &ColumnData{
		Db:      db,
		Table:   table,
		Columns: reply,
	}, err
}

func (c *ColumnData) Convert() (*Table, error) {
	var table Table
	table.Table = c.Table
	table.Db = c.Db
	table.Columns = c.Columns
	table.UniqueIndex = map[string][]*Column{}
	table.NormalIndex = map[string][]*Column{}

	m := make(map[string][]*Column)
	for _, each := range c.Columns {
		m[each.IndexName] = append(m[each.IndexName], each)
	}

	primaryColumns := m[indexPri]
	if len(primaryColumns) == 0 {
		return nil, fmt.Errorf("db:%s, table:%s, missing primary key", c.Db, c.Table)
	}

	if len(primaryColumns) > 1 {
		return nil, fmt.Errorf("db:%s, table:%s, joint primary key is not supported", c.Db, c.Table)
	}

	table.PrimaryKey = primaryColumns[0]
	for indexName, columns := range m {
		if indexName == indexPri {
			continue
		}

		one := columns[0]
		if one.NonUnique == 0 {
			table.UniqueIndex[indexName] = columns
		} else {
			table.NormalIndex[indexName] = columns
		}
	}

	return &table, nil
}
