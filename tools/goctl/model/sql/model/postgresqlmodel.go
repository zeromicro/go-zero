package model

import (
	"database/sql"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var p2m = map[string]string{
	"int8":        "bigint",
	"numeric":     "bigint",
	"float8":      "double",
	"float4":      "float",
	"int2":        "smallint",
	"int4":        "integer",
	"timestamptz": "timestamp",
}

// PostgreSqlModel gets table information from information_schema、pg_catalog
type PostgreSqlModel struct {
	conn sqlx.SqlConn
}

// PostgreColumn describes a column in table
type PostgreColumn struct {
	Num               sql.NullInt32  `db:"num"`
	Field             sql.NullString `db:"field"`
	Type              sql.NullString `db:"type"`
	NotNull           sql.NullBool   `db:"not_null"`
	Comment           sql.NullString `db:"comment"`
	ColumnDefault     sql.NullString `db:"column_default"`
	IdentityIncrement sql.NullInt32  `db:"identity_increment"`
}

// PostgreIndex describes an index for a column
type PostgreIndex struct {
	IndexName  sql.NullString `db:"index_name"`
	IndexId    sql.NullInt32  `db:"index_id"`
	IsUnique   sql.NullBool   `db:"is_unique"`
	IsPrimary  sql.NullBool   `db:"is_primary"`
	ColumnName sql.NullString `db:"column_name"`
	IndexSort  sql.NullInt32  `db:"index_sort"`
}

// NewPostgreSqlModel creates an instance and return
func NewPostgreSqlModel(conn sqlx.SqlConn) *PostgreSqlModel {
	return &PostgreSqlModel{
		conn: conn,
	}
}

// GetAllTables selects all tables from TABLE_SCHEMA
func (m *PostgreSqlModel) GetAllTables(schema string) ([]string, error) {
	query := `select table_name from information_schema.tables where table_schema = $1`
	var tables []string
	err := m.conn.QueryRows(&tables, query, schema)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

// FindColumns return columns in specified database and table
func (m *PostgreSqlModel) FindColumns(schema, table string) (*ColumnData, error) {
	querySql := `select t.num,t.field,t.type,t.not_null,t.comment, c.column_default, identity_increment
from (
         SELECT a.attnum AS num,
                c.relname,
                a.attname     AS field,
                t.typname     AS type,
                a.atttypmod   AS lengthvar,
                a.attnotnull  AS not_null,
                b.description AS comment,
                (c.relnamespace::regnamespace)::varchar AS schema_name
         FROM pg_class c,
              pg_attribute a
                  LEFT OUTER JOIN pg_description b ON a.attrelid = b.objoid AND a.attnum = b.objsubid,
              pg_type t
         WHERE c.relname = $1
           and a.attnum > 0
           and a.attrelid = c.oid
           and a.atttypid = t.oid
 		 GROUP BY a.attnum, c.relname, a.attname, t.typname, a.atttypmod, a.attnotnull, b.description, c.relnamespace::regnamespace
         ORDER BY a.attnum) AS t
         left join information_schema.columns AS c on t.relname = c.table_name and t.schema_name = c.table_schema
		and t.field = c.column_name
		where c.table_schema = $2`

	var reply []*PostgreColumn
	err := m.conn.QueryRowsPartial(&reply, querySql, table, schema)
	if err != nil {
		return nil, err
	}

	list, err := m.getColumns(schema, table, reply)
	if err != nil {
		return nil, err
	}

	var columnData ColumnData
	columnData.Db = schema
	columnData.Table = table
	columnData.Columns = list
	return &columnData, nil
}

func (m *PostgreSqlModel) getColumns(schema, table string, in []*PostgreColumn) ([]*Column, error) {
	index, err := m.getIndex(schema, table)
	if err != nil {
		return nil, err
	}

	var list []*Column
	for _, e := range in {
		var dft any
		if len(e.ColumnDefault.String) > 0 {
			dft = e.ColumnDefault
		}

		isNullAble := "YES"
		if e.NotNull.Bool {
			isNullAble = "NO"
		}

		var extra string
		// when identity is true, the column is auto increment
		if e.IdentityIncrement.Int32 == 1 {
			extra = "auto_increment"
		}

		// when type is serial, it's auto_increment. and the default value is tablename_columnname_seq
		if strings.Contains(e.ColumnDefault.String, table+"_"+e.Field.String+"_seq") {
			extra = "auto_increment"
		}

		if len(index[e.Field.String]) > 0 {
			for _, i := range index[e.Field.String] {
				list = append(list, &Column{
					DbColumn: &DbColumn{
						Name:            e.Field.String,
						DataType:        m.convertPostgreSqlTypeIntoMysqlType(e.Type.String),
						Extra:           extra,
						Comment:         e.Comment.String,
						ColumnDefault:   dft,
						IsNullAble:      isNullAble,
						OrdinalPosition: int(e.Num.Int32),
					},
					Index: i,
				})
			}
		} else {
			list = append(list, &Column{
				DbColumn: &DbColumn{
					Name:            e.Field.String,
					DataType:        m.convertPostgreSqlTypeIntoMysqlType(e.Type.String),
					Extra:           extra,
					Comment:         e.Comment.String,
					ColumnDefault:   dft,
					IsNullAble:      isNullAble,
					OrdinalPosition: int(e.Num.Int32),
				},
			})
		}
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

func (m *PostgreSqlModel) getIndex(schema, table string) (map[string][]*DbIndex, error) {
	indexes, err := m.FindIndex(schema, table)
	if err != nil {
		return nil, err
	}

	index := make(map[string][]*DbIndex)
	for _, e := range indexes {
		if e.IsPrimary.Bool {
			index[e.ColumnName.String] = append(index[e.ColumnName.String], &DbIndex{
				IndexName:  indexPri,
				SeqInIndex: int(e.IndexSort.Int32),
			})
			continue
		}

		nonUnique := 0
		if !e.IsUnique.Bool {
			nonUnique = 1
		}

		index[e.ColumnName.String] = append(index[e.ColumnName.String], &DbIndex{
			IndexName:  e.IndexName.String,
			NonUnique:  nonUnique,
			SeqInIndex: int(e.IndexSort.Int32),
		})
	}

	return index, nil
}

// FindIndex finds index with given schema, table and column.
func (m *PostgreSqlModel) FindIndex(schema, table string) ([]*PostgreIndex, error) {
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
  and E.SCHEMANAME = $1
  and E.RELNAME = $2
    order by C.INDEXRELID,G.attnum`

	var reply []*PostgreIndex
	err := m.conn.QueryRowsPartial(&reply, querySql, schema, table)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
