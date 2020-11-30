package parser

import (
	"fmt"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/converter"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/xwb1989/sqlparser"
)

const (
	_ = iota
	primary
	unique
	normal
	spatial
)

const timeImport = "time.Time"

type (
	Table struct {
		Name       stringx.String
		PrimaryKey Primary
		Fields     []Field
	}

	Primary struct {
		Field
		AutoIncrement bool
	}

	Field struct {
		Name         stringx.String
		DataBaseType string
		DataType     string
		IsPrimaryKey bool
		IsUniqueKey  bool
		Comment      string
	}

	KeyType int
)

func Parse(ddl string) (*Table, error) {
	stmt, err := sqlparser.ParseStrictDDL(ddl)
	if err != nil {
		return nil, err
	}

	ddlStmt, ok := stmt.(*sqlparser.DDL)
	if !ok {
		return nil, unSupportDDL
	}

	action := ddlStmt.Action
	if action != sqlparser.CreateStr {
		return nil, fmt.Errorf("expected [CREATE] action,but found: %s", action)
	}

	tableName := ddlStmt.NewName.Name.String()
	tableSpec := ddlStmt.TableSpec
	if tableSpec == nil {
		return nil, tableBodyIsNotFound
	}

	columns := tableSpec.Columns
	indexes := tableSpec.Indexes
	keyMap := make(map[string]KeyType)
	for _, index := range indexes {
		info := index.Info
		if info == nil {
			continue
		}
		if info.Primary {
			if len(index.Columns) > 1 {
				return nil, errPrimaryKey
			}

			keyMap[index.Columns[0].Column.String()] = primary
			continue
		}
		// can optimize
		if len(index.Columns) > 1 {
			continue
		}
		column := index.Columns[0]
		columnName := column.Column.String()
		camelColumnName := stringx.From(columnName).ToCamel()
		// by default, createTime|updateTime findOne is not used.
		if camelColumnName == "CreateTime" || camelColumnName == "UpdateTime" {
			continue
		}
		if info.Unique {
			keyMap[columnName] = unique
		} else if info.Spatial {
			keyMap[columnName] = spatial
		} else {
			keyMap[columnName] = normal
		}
	}

	var fields []Field
	var primaryKey Primary
	for _, column := range columns {
		if column == nil {
			continue
		}
		var comment string
		if column.Type.Comment != nil {
			comment = string(column.Type.Comment.Val)
		}
		var isDefaultNull = true
		if column.Type.NotNull {
			isDefaultNull = false
		} else {
			if column.Type.Default == nil {
				isDefaultNull = false
			} else if string(column.Type.Default.Val) != "null" {
				isDefaultNull = false
			}
		}
		dataType, err := converter.ConvertDataType(column.Type.Type, isDefaultNull)
		if err != nil {
			return nil, err
		}

		var field Field
		field.Name = stringx.From(column.Name.String())
		field.DataBaseType = column.Type.Type
		field.DataType = dataType
		field.Comment = comment
		key, ok := keyMap[column.Name.String()]
		if ok {
			field.IsPrimaryKey = key == primary
			field.IsUniqueKey = key == unique
			if field.IsPrimaryKey {
				primaryKey.Field = field
				if column.Type.Autoincrement {
					primaryKey.AutoIncrement = true
				}
			}
		}
		fields = append(fields, field)
	}

	return &Table{
		Name:       stringx.From(tableName),
		PrimaryKey: primaryKey,
		Fields:     fields,
	}, nil
}

func (t *Table) ContainsTime() bool {
	for _, item := range t.Fields {
		if item.DataType == timeImport {
			return true
		}
	}
	return false
}

func ConvertColumn(db, table string, in []*model.Column) (*Table, error) {
	var reply Table
	reply.Name = stringx.From(table)
	keyMap := make(map[string][]*model.Column)

	for _, column := range in {
		keyMap[column.Key] = append(keyMap[column.Key], column)
	}
	primaryColumns := keyMap["PRI"]
	if len(primaryColumns) == 0 {
		return nil, fmt.Errorf("database:%s, table %s: missing primary key", db, table)
	}

	if len(primaryColumns) > 1 {
		return nil, fmt.Errorf("database:%s, table %s: only one primary key expected", db, table)
	}

	primaryColumn := primaryColumns[0]
	isDefaultNull := primaryColumn.ColumnDefault == nil && primaryColumn.IsNullAble == "YES"
	primaryFt, err := converter.ConvertDataType(primaryColumn.DataType, isDefaultNull)
	if err != nil {
		return nil, err
	}

	primaryField := Field{
		Name:         stringx.From(primaryColumn.Name),
		DataBaseType: primaryColumn.DataType,
		DataType:     primaryFt,
		IsUniqueKey:  true,
		IsPrimaryKey: true,
		Comment:      primaryColumn.Comment,
	}
	reply.PrimaryKey = Primary{
		Field:         primaryField,
		AutoIncrement: strings.Contains(primaryColumn.Extra, "auto_increment"),
	}
	for key, columns := range keyMap {
		for _, item := range columns {
			isColumnDefaultNull := item.ColumnDefault == nil && item.IsNullAble == "YES"
			dt, err := converter.ConvertDataType(item.DataType, isColumnDefaultNull)
			if err != nil {
				return nil, err
			}

			f := Field{
				Name:         stringx.From(item.Name),
				DataBaseType: item.DataType,
				DataType:     dt,
				IsPrimaryKey: primaryColumn.Name == item.Name,
				Comment:      item.Comment,
			}
			if key == "UNI" {
				f.IsUniqueKey = true
			}
			reply.Fields = append(reply.Fields, f)
		}
	}

	return &reply, nil
}
