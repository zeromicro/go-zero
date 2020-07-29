package gen

import (
	"bytes"
	"strings"
	"text/template"

	sqltemplate "zero/tools/goctl/model/sql/template"
)

func genInsert(table *InnerTable) (string, error) {
	t, err := template.New("insert").Parse(sqltemplate.Insert)
	if err != nil {
		return "", nil
	}
	insertBuffer := new(bytes.Buffer)
	expressions := make([]string, 0)
	expressionValues := make([]string, 0)
	for _, filed := range table.Fields {
		if filed.SnakeCase == "create_time" || filed.SnakeCase == "update_time" || filed.IsPrimaryKey {
			continue
		}
		expressions = append(expressions, "?")
		expressionValues = append(expressionValues, "data."+filed.UpperCamelCase)
	}
	err = t.Execute(insertBuffer, map[string]interface{}{
		"upperObject":      table.UpperCamelCase,
		"lowerObject":      table.LowerCamelCase,
		"expression":       strings.Join(expressions, ", "),
		"expressionValues": strings.Join(expressionValues, ", "),
		"containsCache":    table.ContainsCache,
	})
	if err != nil {
		return "", err
	}
	return insertBuffer.String(), nil
}
