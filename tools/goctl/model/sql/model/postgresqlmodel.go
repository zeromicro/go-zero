package model

import (
	"strings"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

var (
	p2m = map[string]string{
		"int8":    "bigint",
		"numeric": "bigint",
		"float8":  "double",
		"float4":  "float",
		"int2":    "smallint",
		"int4":    "integer",
	}
)

// PostgreSqlModel gets table information from information_schema、pg_catalog
type PostgreSqlModel struct {
	conn sqlx.SqlConn
}

// PostgreColumn describes a column in table
type PostgreColumn struct {
	Num               int    `db:"num"`
	Field             string `db:"field"`
	Type              string `db:"type"`
	NotNull           bool   `db:"not_null"`
	Comment           string `db:"comment"`
	ColumnDefault     string `db:"column_default"`
	IdentityIncrement int    `db:"identity_increment"`
}

// PostgreIndex describes an index for a column
type PostgreIndex struct {
	IndexName  string `db:"index_name"`
	IndexId    int64  `db:"index_id"`
	IsUnique   bool   `db:"is_unique"`
	IsPrimary  bool   `db:"is_primary"`
	ColumnName string `db:"column_name"`
	IndexSort  int    `db:"index_sort"`
}

// NewPostgreSqlModel creates an instance and return
func NewPostgreSqlModel(conn sqlx.SqlConn) *PostgreSqlModel {
	return &PostgreSqlModel{
		conn: conn,
	}
}

// GetAllTables selects all tables from TABLE_SCHEMA
func (m *PostgreSqlModel) GetAllTables(database string) ([]string, error) {
	query := `select table_name from information_schema.tables where table_schema = ?;`
	var tables []string
	err := m.conn.QueryRows(&tables, query, database)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

// FindColumns return columns in specified database and table
func (m *PostgreSqlModel) FindColumns(db, table string) (*ColumnData, error) {
	querySql := `select t.num,t.field,t.type,t.not_null,t.comment, c.column_default, identity_increment
from (
         SELECT a.attnum AS num,
                c.relname,
                a.attname     AS field,
                t.typname     AS type,
                a.atttypmod   AS lengthvar,
                a.attnotnull  AS not_null,
                b.description AS comment
         FROM pg_class c,
              pg_attribute a
                  LEFT OUTER JOIN pg_description b ON a.attrelid = b.objoid AND a.attnum = b.objsubid,
              pg_type t
         WHERE c.relname = ?
           and a.attnum > 0
           and a.attrelid = c.oid
           and a.atttypid = t.oid
         ORDER BY a.attnum) AS t
         left join information_schema.columns AS c on t.relname = c.table_name 
		and t.field = c.column_name and c.table_schema = ?`

	var reply []*PostgreColumn
	err := m.conn.QueryRowsPartial(&reply, querySql, db, table)
	if err != nil {
		return nil, err
	}

	list, err := m.getColumns(db, table, reply)
	if err != nil {
		return nil, err
	}

	var columnData ColumnData
	columnData.Db = db
	columnData.Table = table
	columnData.Columns = list
	return &columnData, nil
}

func (m *PostgreSqlModel) getColumns(db, table string, in []*PostgreColumn) ([]*Column, error) {
	index, err := m.getIndex(db, table)
	if err != nil {
		return nil, err
	}
	var list []*Column
	for _, e := range in {
		var dft interface{}
		if len(e.ColumnDefault) > 0 {
			dft = e.ColumnDefault
		}

		isNullAble := "YES"
		if e.NotNull {
			isNullAble = "NO"
		}

		extra := "auto_increment"
		if e.IdentityIncrement != 1 {
			extra = ""
		}

		list = append(list, &Column{
			DbColumn: &DbColumn{
				Name:            e.Field,
				DataType:        m.convertPostgreSqlTypeIntoMysqlType(e.Type),
				Extra:           extra,
				Comment:         e.Comment,
				ColumnDefault:   dft,
				IsNullAble:      isNullAble,
				OrdinalPosition: e.Num,
			},
			Index: index[e.Field],
		})
	}

	return list, nil
}

func (m *PostgreSqlModel) convertPostgreSqlTypeIntoMysqlType(in string) string {
	r, ok := p2m[strings.ToLower(in)]
	if ok {
		return r
	}

	return in
}

func (m *PostgreSqlModel) getIndex(db, table string) (map[string]*DbIndex, error) {
	indexes, err := m.FindIndex(db, table)
	if err != nil {
		return nil, err
	}
	var index = make(map[string]*DbIndex)
	for _, e := range indexes {
		if e.IsPrimary {
			index[e.ColumnName] = &DbIndex{
				IndexName:  indexPri,
				SeqInIndex: e.IndexSort,
			}
			continue
		}

		nonUnique := 0
		if !e.IsUnique {
			nonUnique = 1
		}

		index[e.ColumnName] = &DbIndex{
			IndexName:  e.IndexName,
			NonUnique:  nonUnique,
			SeqInIndex: e.IndexSort,
		}
	}
	return index, nil
}

// FindIndex finds index with given db, table and column.
func (m *PostgreSqlModel) FindIndex(db, table string) ([]*PostgreIndex, error) {
	querySql := `select A.INDEXNAME AS index_name,
       C.INDEXRELID AS index_id,
       C.INDISUNIQUE AS is_unique,
       C.INDISPRIMARY AS is_primary,
       G.ATTNAME AS column_name,
       G.attnum AS index_sort
from PG_AM B
         left join PG_CLASS F on
    B.OID = F.RELAM
         left join PG_STAT_ALL_INDEXES E on
    F.OID = E.INDEXRELID
         left join PG_INDEX C on
    E.INDEXRELID = C.INDEXRELID
         left outer join PG_DESCRIPTION D on
    C.INDEXRELID = D.OBJOID,
     PG_INDEXES A,
     pg_attribute G
where A.SCHEMANAME = E.SCHEMANAME
  and A.TABLENAME = E.RELNAME
  and A.INDEXNAME = E.INDEXRELNAME
  and F.oid = G.attrelid
  and E.SCHEMANAME = ?
  and E.RELNAME = ?
    order by C.INDEXRELID,G.attnum`

	var reply []*PostgreIndex
	err := m.conn.QueryRowsPartial(&reply, querySql, db, table)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
