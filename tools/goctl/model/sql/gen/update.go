package gen

import (
	"bytes"
	"strings"
	"text/template"

	sqltemplate "zero/tools/goctl/model/sql/template"
)

func genUpdate(table *InnerTable) (string, error) {
	t, err := template.New("update").Parse(sqltemplate.Update)
	if err != nil {
		return "", nil
	}
	updateBuffer := new(bytes.Buffer)
	expressionValues := make([]string, 0)
	for _, filed := range table.Fields {
		if filed.SnakeCase == "create_time" || filed.SnakeCase == "update_time" || filed.IsPrimaryKey {
			continue
		}
		expressionValues = append(expressionValues, "data."+filed.UpperCamelCase)
	}
	expressionValues = append(expressionValues, "data."+table.PrimaryField.UpperCamelCase)
	err = t.Execute(updateBuffer, map[string]interface{}{
		"containsCache":      table.ContainsCache,
		"upperObject":        table.UpperCamelCase,
		"primaryCacheKey":    table.CacheKey[table.PrimaryField.SnakeCase].DataKey,
		"primaryKeyVariable": table.CacheKey[table.PrimaryField.SnakeCase].KeyVariable,
		"lowerObject":        table.LowerCamelCase,
		"primarySnakeCase":   table.PrimaryField.SnakeCase,
		"expressionValues":   strings.Join(expressionValues, ", "),
	})
	if err != nil {
		return "", err
	}
	return updateBuffer.String(), nil
}
