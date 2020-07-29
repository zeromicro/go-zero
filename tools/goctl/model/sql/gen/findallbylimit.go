package gen

import (
	"bytes"
	"strings"
	"text/template"

	sqltemplate "zero/tools/goctl/model/sql/template"
)

func genFindLimitByField(table *InnerTable) (string, error) {
	t, err := template.New("findLimitByField").Parse(sqltemplate.FindLimitByField)
	if err != nil {
		return "", err
	}
	list := make([]string, 0)
	for _, field := range table.Fields {
		if field.IsPrimaryKey {
			continue
		}
		if field.QueryType != QueryLimit {
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
		sortList := make([]string, 0)
		for _, item := range field.Sort {
			var sort = "ASC"
			if !item.Asc {
				sort = "DESC"
			}
			sortList = append(sortList, item.Field.SnakeCase+" "+sort)
		}
		err = t.Execute(fineOneByFieldBuffer, map[string]interface{}{
			"in":               strings.Join(in, ","),
			"upperObject":      table.UpperCamelCase,
			"upperFields":      strings.Join(upperFields, "And"),
			"lowerObject":      table.LowerCamelCase,
			"expression":       strings.Join(expressionFields, " AND "),
			"expressionValues": strings.Join(expressionValuesFields, ", "),
			"sortExpression":   strings.Join(sortList, ","),
			"containsCache":    table.ContainsCache,
		})
		if err != nil {
			return "", err
		}
		list = append(list, fineOneByFieldBuffer.String())
	}
	return strings.Join(list, ""), nil
}
