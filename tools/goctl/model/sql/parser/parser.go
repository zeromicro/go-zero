package parser

import (
	"fmt"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/converter"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/xwb1989/sqlparser"
)

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

	for _, index := range indexes {
		info := index.Info
		if info == nil {
			continue
		}
		if info.Primary {
			if len(index.Columns) > 1 {
				return nil, errPrimaryKey
			}
			break
		}
	}
	var (
		fields     []Field
		primaryKey Primary
	)
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
		// see line 1194 https://github.com/xwb1989/sqlparser/blob/master/ast.go
		field.IsKey = column.Type.KeyOpt != 0
		field.IsPrimaryKey = column.Type.KeyOpt == 1
		fields = append(fields, field)
		// see line 1195 https://github.com/xwb1989/sqlparser/blob/master/ast.go
		if column.Type.KeyOpt == 1 {
			primaryKey.Field = field
			if column.Type.Autoincrement {
				primaryKey.AutoIncrement = true
			}
		}
	}
	return &Table{
		Name:       stringx.From(tableName),
		PrimaryKey: primaryKey,
		Fields:     fields,
	}, nil

}
