package gen

import (
	"bytes"
	"strings"
	"text/template"

	sqltemplate "zero/tools/goctl/model/sql/template"
)

func genFineOneByField(table *InnerTable) (string, error) {
	t, err := template.New("findOneByField").Parse(sqltemplate.FindOneByField)
	if err != nil {
		return "", err
	}
	list := make([]string, 0)
	for _, field := range table.Fields {
		if field.IsPrimaryKey {
			continue
		}
		if field.QueryType != QueryOne {
			continue
		}
		fineOneByFieldBuffer := new(bytes.Buffer)
		upperFields := make([]string, 0)
		in := make([]string, 0)
		expressionFields := make([]string, 0)
		expressionValuesFields := make([]string, 0)
		upperFields = append(upperFields, field.UpperCamelCase)
		in = append(in, field.LowerCamelCase+" "+field.DataType)
		expressionFields = append(expressionFields, field.SnakeCase+" = ?")
		expressionValuesFields = append(expressionValuesFields, field.LowerCamelCase)
		for _, withField := range field.WithFields {
			upperFields = append(upperFields, withField.UpperCamelCase)
			in = append(in, withField.LowerCamelCase+" "+withField.DataType)
			expressionFields = append(expressionFields, withField.SnakeCase+" = ?")
			expressionValuesFields = append(expressionValuesFields, withField.LowerCamelCase)
		}
		err = t.Execute(fineOneByFieldBuffer, map[string]interface{}{
			"in":                    strings.Join(in, ","),
			"upperObject":           table.UpperCamelCase,
			"upperFields":           strings.Join(upperFields, "And"),
			"onlyOneFiled":          len(field.WithFields) == 0,
			"withCache":             field.Cache,
			"containsCache":         table.ContainsCache,
			"lowerObject":           table.LowerCamelCase,
			"lowerField":            field.LowerCamelCase,
			"snakeField":            field.SnakeCase,
			"lowerPrimaryKey":       table.PrimaryField.LowerCamelCase,
			"UpperPrimaryKey":       table.PrimaryField.UpperCamelCase,
			"primaryKeyDefine":      table.CacheKey[table.PrimaryField.SnakeCase].Define,
			"primarySnakeField":     table.PrimaryField.SnakeCase,
			"primaryDataType":       table.PrimaryField.DataType,
			"primaryDataTypeString": table.PrimaryField.DataType == "string",
			"upperObjectKey":        table.PrimaryField.UpperCamelCase,
			"cacheKey":              table.CacheKey[field.SnakeCase].Key,
			"cacheKeyVariable":      table.CacheKey[field.SnakeCase].KeyVariable,
			"expression":            strings.Join(expressionFields, " AND "),
			"expressionValues":      strings.Join(expressionValuesFields, ", "),
		})
		if err != nil {
			return "", err
		}
		list = append(list, fineOneByFieldBuffer.String())
	}
	return strings.Join(list, ""), nil
}
