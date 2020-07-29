package gen

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"zero/tools/goctl/model/sql/util"
)

func TableConvert(outerTable OuterTable) (*InnerTable, error) {
	var table InnerTable
	table.CreateNotFound = outerTable.CreateNotFound
	tableSnakeCase, tableUpperCamelCase, tableLowerCamelCase := util.FormatField(outerTable.Table)
	table.SnakeCase = tableSnakeCase
	table.UpperCamelCase = tableUpperCamelCase
	table.LowerCamelCase = tableLowerCamelCase
	fields := make([]*InnerField, 0)
	var primaryField *InnerField
	conflict := make(map[string]struct{})
	var containsCache bool
	for _, field := range outerTable.Fields {
		if field.Cache && !containsCache {
			containsCache = true
		}
		fieldSnakeCase, fieldUpperCamelCase, fieldLowerCamelCase := util.FormatField(field.Name)
		tag, err := genTag(fieldSnakeCase)
		if err != nil {
			return nil, err
		}
		var comment string
		if field.Comment != "" {
			comment = fmt.Sprintf("// %s", field.Comment)
		}
		withFields := make([]InnerWithField, 0)
		unique := make([]string, 0)
		unique = append(unique, fmt.Sprintf("%v", field.QueryType))
		unique = append(unique, field.Name)

		for _, item := range field.WithFields {
			unique = append(unique, item.Name)
			withFieldSnakeCase, withFieldUpperCamelCase, withFieldLowerCamelCase := util.FormatField(item.Name)
			withFields = append(withFields, InnerWithField{
				Case: Case{
					SnakeCase:      withFieldSnakeCase,
					LowerCamelCase: withFieldLowerCamelCase,
					UpperCamelCase: withFieldUpperCamelCase,
				},
				DataType: commonMysqlDataTypeMap[item.DataBaseType],
			})
		}
		sort.Strings(unique)
		uniqueKey := strings.Join(unique, "#")
		if _, ok := conflict[uniqueKey]; ok {
			return nil, ErrCircleQuery
		} else {
			conflict[uniqueKey] = struct{}{}
		}
		sortFields := make([]InnerSort, 0)
		for _, sortField := range field.OuterSort {
			sortSnake, sortUpperCamelCase, sortLowerCamelCase := util.FormatField(sortField.Field)
			sortFields = append(sortFields, InnerSort{
				Field: Case{
					SnakeCase:      sortSnake,
					LowerCamelCase: sortUpperCamelCase,
					UpperCamelCase: sortLowerCamelCase,
				},
				Asc: sortField.Asc,
			})
		}
		innerField := &InnerField{
			IsPrimaryKey: field.IsPrimaryKey,
			InnerWithField: InnerWithField{
				Case: Case{
					SnakeCase:      fieldSnakeCase,
					LowerCamelCase: fieldLowerCamelCase,
					UpperCamelCase: fieldUpperCamelCase,
				},
				DataType: commonMysqlDataTypeMap[field.DataBaseType],
			},
			DataBaseType: field.DataBaseType,
			Tag:          tag,
			Comment:      comment,
			Cache:        field.Cache,
			QueryType:    field.QueryType,
			WithFields:   withFields,
			Sort:         sortFields,
		}
		if field.IsPrimaryKey {
			primaryField = innerField
		}
		fields = append(fields, innerField)
	}
	if primaryField == nil {
		return nil, errors.New("please ensure that primary exists")
	}
	table.ContainsCache = containsCache
	primaryField.Cache = containsCache
	table.PrimaryField = primaryField
	table.Fields = fields
	cacheKey, err := genCacheKeys(&table)
	if err != nil {
		return nil, err
	}
	table.CacheKey = cacheKey
	return &table, nil
}
