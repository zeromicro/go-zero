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
	// Table describes a mysql table
	Table struct {
		Name        stringx.String
		PrimaryKey  Primary
		UniqueIndex map[string][]*Field
		NormalIndex map[string][]*Field
		Fields      []*Field
	}

	// Primary describes a primary key
	Primary struct {
		Field
		AutoIncrement bool
	}

	// Field describes a table field
	Field struct {
		Name            stringx.String
		DataBaseType    string
		DataType        string
		Comment         string
		SeqInIndex      int
		OrdinalPosition int
	}

	// KeyType types alias of int
	KeyType int
)

// Parse parses ddl into golang structure
func Parse(ddl string) (*Table, error) {
	stmt, err := sqlparser.ParseStrictDDL(ddl)
	if err != nil {
		return nil, err
	}

	ddlStmt, ok := stmt.(*sqlparser.DDL)
	if !ok {
		return nil, errUnsupportDDL
	}

	action := ddlStmt.Action
	if action != sqlparser.CreateStr {
		return nil, fmt.Errorf("expected [CREATE] action,but found: %s", action)
	}

	tableName := ddlStmt.NewName.Name.String()
	tableSpec := ddlStmt.TableSpec
	if tableSpec == nil {
		return nil, errTableBodyNotFound
	}

	columns := tableSpec.Columns
	indexes := tableSpec.Indexes
	primaryColumn, uniqueKeyMap, normalKeyMap, err := convertIndexes(indexes)
	if err != nil {
		return nil, err
	}

	fields, primaryKey, fieldM, err := convertColumns(columns, primaryColumn)
	if err != nil {
		return nil, err
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

	for indexName, each := range normalKeyMap {
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

func convertColumns(columns []*sqlparser.ColumnDefinition, primaryColumn string) ([]*Field, Primary, map[string]*Field, error) {
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
			return nil, Primary{}, nil, err
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
	return fields, primaryKey, fieldM, nil
}

func convertIndexes(indexes []*sqlparser.IndexDefinition) (string, map[string][]string, map[string][]string, error) {
	var primaryColumn string
	uniqueKeyMap := make(map[string][]string)
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
				return "", nil, nil, errPrimaryKey
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
			// do nothing
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
	return primaryColumn, uniqueKeyMap, normalKeyMap, nil
}

// ContainsTime returns true if contains golang type time.Time
func (t *Table) ContainsTime() bool {
	for _, item := range t.Fields {
		if item.DataType == timeImport {
			return true
		}
	}
	return false
}

// ConvertDataType converts mysql data type into golang data type
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
	seqInIndex := 0
	if table.PrimaryKey.Index != nil {
		seqInIndex = table.PrimaryKey.Index.SeqInIndex
	}

	reply.PrimaryKey = Primary{
		Field: Field{
			Name:            stringx.From(table.PrimaryKey.Name),
			DataBaseType:    table.PrimaryKey.DataType,
			DataType:        primaryDataType,
			Comment:         table.PrimaryKey.Comment,
			SeqInIndex:      seqInIndex,
			OrdinalPosition: table.PrimaryKey.OrdinalPosition,
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
		columnSeqInIndex := 0
		if each.Index != nil {
			columnSeqInIndex = each.Index.SeqInIndex
		}

		field := &Field{
			Name:            stringx.From(each.Name),
			DataBaseType:    each.DataType,
			DataType:        dt,
			Comment:         each.Comment,
			SeqInIndex:      columnSeqInIndex,
			OrdinalPosition: each.OrdinalPosition,
		}
		fieldM[each.Name] = field
	}

	for _, each := range fieldM {
		reply.Fields = append(reply.Fields, each)
	}
	sort.Slice(reply.Fields, func(i, j int) bool {
		return reply.Fields[i].OrdinalPosition < reply.Fields[j].OrdinalPosition
	})

	uniqueIndexSet := collection.NewSet()
	log := console.NewColorConsole()
	for indexName, each := range table.UniqueIndex {
		sort.Slice(each, func(i, j int) bool {
			if each[i].Index != nil {
				return each[i].Index.SeqInIndex < each[j].Index.SeqInIndex
			}
			return false
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
