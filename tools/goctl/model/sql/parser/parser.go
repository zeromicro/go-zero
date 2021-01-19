package parser

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tal-tech/go-zero/core/collection"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/converter"
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/model"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
	"github.com/xwb1989/sqlparser"
)

const timeImport = "time.Time"

type (
	Table struct {
		Name        stringx.String
		PrimaryKey  Primary
		UniqueIndex map[string][]*Field
		NormalIndex map[string][]*Field
		Fields      []*Field
	}

	Primary struct {
		Field
		AutoIncrement bool
	}

	Field struct {
		Name         stringx.String
		DataBaseType string
		DataType     string
		Comment      string
		SeqInIndex   int
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
	var primaryColumn string
	uniqueKeyMap := make(map[string][]string)
	spatialKeyMap := make(map[string][]string)
	normalKeyMap := make(map[string][]string)

	isCreateTimeOrUpdateTime := func(name string) bool {
		camelColumnName := stringx.From(name).ToCamel()
		// by default, createTime|updateTime findOne is not used.
		return camelColumnName == "CreateTime" || camelColumnName == "UpdateTime"
	}

	for _, index := range indexes {
		info := index.Info
		if info == nil {
			continue
		}

		indexName := index.Info.Name.String()
		if info.Primary {
			if len(index.Columns) > 1 {
				return nil, errPrimaryKey
			}
			columnName := index.Columns[0].Column.String()
			if isCreateTimeOrUpdateTime(columnName) {
				continue
			}

			primaryColumn = columnName
			continue
		} else if info.Unique {
			for _, each := range index.Columns {
				columnName := each.Column.String()
				if isCreateTimeOrUpdateTime(columnName) {
					break
				}

				uniqueKeyMap[indexName] = append(uniqueKeyMap[indexName], columnName)
			}
		} else if info.Spatial {
			for _, each := range index.Columns {
				columnName := each.Column.String()
				if isCreateTimeOrUpdateTime(columnName) {
					break
				}

				spatialKeyMap[indexName] = append(spatialKeyMap[indexName], each.Column.String())
			}
		} else {
			for _, each := range index.Columns {
				columnName := each.Column.String()
				if isCreateTimeOrUpdateTime(columnName) {
					break
				}

				normalKeyMap[indexName] = append(normalKeyMap[indexName], each.Column.String())
			}
		}
	}

	var (
		fields     []*Field
		primaryKey Primary
		fieldM     = make(map[string]*Field)
	)

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

		if field.Name.Source() == primaryColumn {
			primaryKey = Primary{
				Field:         field,
				AutoIncrement: bool(column.Type.Autoincrement),
			}
		}

		fields = append(fields, &field)
		fieldM[field.Name.Source()] = &field
	}

	var (
		uniqueIndex = make(map[string][]*Field)
		normalIndex = make(map[string][]*Field)
	)
	for indexName, each := range uniqueKeyMap {
		for _, columnName := range each {
			uniqueIndex[indexName] = append(uniqueIndex[indexName], fieldM[columnName])
		}
	}

	for indexName, each := range uniqueKeyMap {
		for _, columnName := range each {
			normalIndex[indexName] = append(normalIndex[indexName], fieldM[columnName])
		}
	}

	return &Table{
		Name:        stringx.From(tableName),
		PrimaryKey:  primaryKey,
		UniqueIndex: uniqueIndex,
		NormalIndex: normalIndex,
		Fields:      fields,
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

func ConvertDataType(table *model.Table) (*Table, error) {
	isPrimaryDefaultNull := table.PrimaryKey.ColumnDefault == nil && table.PrimaryKey.IsNullAble == "YES"
	primaryDataType, err := converter.ConvertDataType(table.PrimaryKey.DataType, isPrimaryDefaultNull)
	if err != nil {
		return nil, err
	}

	var reply Table
	reply.UniqueIndex = map[string][]*Field{}
	reply.NormalIndex = map[string][]*Field{}
	reply.Name = stringx.From(table.Table)
	reply.PrimaryKey = Primary{
		Field: Field{
			Name:         stringx.From(table.PrimaryKey.Name),
			DataBaseType: table.PrimaryKey.DataType,
			DataType:     primaryDataType,
			Comment:      table.PrimaryKey.Comment,
			SeqInIndex:   table.PrimaryKey.SeqInIndex,
		},
		AutoIncrement: strings.Contains(table.PrimaryKey.Extra, "auto_increment"),
	}

	fieldM := make(map[string]*Field)
	for _, each := range table.Columns {
		isDefaultNull := each.ColumnDefault == nil && each.IsNullAble == "YES"
		dt, err := converter.ConvertDataType(each.DataType, isDefaultNull)
		if err != nil {
			return nil, err
		}
		field := &Field{
			Name:         stringx.From(each.Name),
			DataBaseType: each.DataType,
			DataType:     dt,
			Comment:      each.Comment,
			SeqInIndex:   each.SeqInIndex,
		}
		fieldM[each.Name] = field
	}

	for _, each := range fieldM {
		reply.Fields = append(reply.Fields, each)
	}

	uniqueIndexSet := collection.NewSet()
	log := console.NewColorConsole()
	for indexName, each := range table.UniqueIndex {
		sort.Slice(each, func(i, j int) bool {
			return each[i].SeqInIndex < each[j].SeqInIndex
		})
		if len(each) == 1 {
			one := each[0]
			if one.Name == table.PrimaryKey.Name {
				log.Warning("duplicate unique index with primary key, %s", one.Name)
				continue
			}
		}

		var list []*Field
		var uniqueJoin []string
		for _, c := range each {
			list = append(list, fieldM[c.Name])
			uniqueJoin = append(uniqueJoin, c.Name)
		}

		uniqueKey := strings.Join(uniqueJoin, ",")
		if uniqueIndexSet.Contains(uniqueKey) {
			log.Warning("duplicate unique index, %s", uniqueKey)
			continue
		}

		reply.UniqueIndex[indexName] = list
	}

	for indexName, each := range table.NormalIndex {
		var list []*Field
		for _, c := range each {
			list = append(list, fieldM[c.Name])
		}

		sort.Slice(list, func(i, j int) bool {
			return list[i].SeqInIndex < list[j].SeqInIndex
		})

		reply.NormalIndex[indexName] = list
	}

	return &reply, nil
}
