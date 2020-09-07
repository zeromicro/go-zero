package parser

import (
	"fmt"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/converter"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/xwb1989/sqlparser"
)

const (
	none = iota
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
		IsKey        bool
		IsPrimaryKey bool
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
		dataType, err := converter.ConvertDataType(column.Type.Type)
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
			field.IsKey = true
			field.IsPrimaryKey = key == primary
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
