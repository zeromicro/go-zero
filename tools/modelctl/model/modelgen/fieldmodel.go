package modelgen

import (
	"zero/core/stores/sqlx"
)

type (
	FieldModel struct {
		dataSource string
		conn       sqlx.SqlConn
		table      string
	}
	Field struct {
		// 字段名称,下划线
		Name string `db:"name"`
		// 字段数据类型
		Type string `db:"type"`
		// 字段顺序
		Position int `db:"position"`
		// 字段注释
		Comment string `db:"comment"`
		// key
		Primary string `db:"k"`
	}

	Table struct {
		Name string `db:"name"`
	}
)

func NewFieldModel(dataSource, table string) *FieldModel {
	return &FieldModel{conn: sqlx.NewMysql(dataSource), table: table}
}

func (fm *FieldModel) findTables() ([]string, error) {
	querySql := `select TABLE_NAME AS name from COLUMNS where TABLE_SCHEMA = ? GROUP BY TABLE_NAME`
	var tables []*Table
	err := fm.conn.QueryRows(&tables, querySql, fm.table)
	if err != nil {
		return nil, err
	}
	tableList := make([]string, 0)
	for _, item := range tables {
		tableList = append(tableList, item.Name)
	}
	return tableList, nil
}

func (fm *FieldModel) findColumns(tableName string) ([]*Field, error) {
	querySql := `select ` + queryRows + ` from COLUMNS where TABLE_SCHEMA = ? and TABLE_NAME = ?`
	var resp []*Field
	err := fm.conn.QueryRows(&resp, querySql, fm.table, tableName)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
