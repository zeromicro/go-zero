package model

import (
	"fmt"
	"sort"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/tools/goctl/model/sql/util"
)

const indexPri = "PRIMARY"

type (
	// InformationSchemaModel defines information schema model
	InformationSchemaModel struct {
		conn sqlx.SqlConn
	}

	// Column defines column in table
	Column struct {
		*DbColumn
		Index *DbIndex
	}

	// DbColumn defines column info of columns
	DbColumn struct {
		Name            string `db:"COLUMN_NAME"`
		DataType        string `db:"DATA_TYPE"`
		ColumnType      string `db:"COLUMN_TYPE"`
		Extra           string `db:"EXTRA"`
		Comment         string `db:"COLUMN_COMMENT"`
		ColumnDefault   any    `db:"COLUMN_DEFAULT"`
		IsNullAble      string `db:"IS_NULLABLE"`
		OrdinalPosition int    `db:"ORDINAL_POSITION"`
	}

	// DbIndex defines index of columns in information_schema.statistic
	DbIndex struct {
		IndexName  string `db:"INDEX_NAME"`
		NonUnique  int    `db:"NON_UNIQUE"`
		SeqInIndex int    `db:"SEQ_IN_INDEX"`
	}

	// ColumnData describes the columns of table
	ColumnData struct {
		Db      string
		Table   string
		Columns []*Column
	}

	// Table describes mysql table which contains database name, table name, columns, keys
	Table struct {
		Db      string
		Table   string
		Columns []*Column
		// Primary key not included
		UniqueIndex map[string][]*Column
		PrimaryKey  *Column
		NormalIndex map[string][]*Column
	}

	// IndexType describes an alias of string
	IndexType string

	// Index describes a column index
	Index struct {
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

// FindColumns return columns in specified database and table
func (m *InformationSchemaModel) FindColumns(db, table string) (*ColumnData, error) {
	querySql := `SELECT c.COLUMN_NAME,c.DATA_TYPE,c.COLUMN_TYPE,EXTRA,c.COLUMN_COMMENT,c.COLUMN_DEFAULT,c.IS_NULLABLE,c.ORDINAL_POSITION from COLUMNS c WHERE c.TABLE_SCHEMA = ? and c.TABLE_NAME = ? `
	var reply []*DbColumn
	err := m.conn.QueryRowsPartial(&reply, querySql, db, table)
	if err != nil {
		return nil, err
	}

	var list []*Column
	for _, item := range reply {
		index, err := m.FindIndex(db, table, item.Name)
		if err != nil {
			if err != sqlx.ErrNotFound {
				return nil, err
			}
			continue
		}

		if len(index) > 0 {
			for _, i := range index {
				list = append(list, &Column{
					DbColumn: item,
					Index:    i,
				})
			}
		} else {
			list = append(list, &Column{
				DbColumn: item,
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].OrdinalPosition < list[j].OrdinalPosition
	})

	var columnData ColumnData
	columnData.Db = db
	columnData.Table = table
	columnData.Columns = list
	return &columnData, nil
}

// FindIndex finds index with given db, table and column.
func (m *InformationSchemaModel) FindIndex(db, table, column string) ([]*DbIndex, error) {
	querySql := `SELECT s.INDEX_NAME,s.NON_UNIQUE,s.SEQ_IN_INDEX from  STATISTICS s  WHERE  s.TABLE_SCHEMA = ? and s.TABLE_NAME = ? and s.COLUMN_NAME = ?`
	var reply []*DbIndex
	err := m.conn.QueryRowsPartial(&reply, querySql, db, table, column)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Convert converts column data into Table
func (c *ColumnData) Convert() (*Table, error) {
	var table Table
	table.Table = c.Table
	table.Db = c.Db
	table.Columns = c.Columns
	table.UniqueIndex = map[string][]*Column{}
	table.NormalIndex = map[string][]*Column{}

	m := make(map[string][]*Column)
	for _, each := range c.Columns {
		each.Comment = util.TrimNewLine(each.Comment)
		if each.Index != nil {
			m[each.Index.IndexName] = append(m[each.Index.IndexName], each)
		}
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

		for _, one := range columns {
			if one.Index != nil {
				if one.Index.NonUnique == 0 {
					table.UniqueIndex[indexName] = columns
				} else {
					table.NormalIndex[indexName] = columns
				}
			}
		}
	}

	return &table, nil
}
