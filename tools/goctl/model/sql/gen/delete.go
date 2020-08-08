package gen

import (
	"bytes"
	"strings"
	"text/template"

	sqltemplate "github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
)

func genDelete(table *InnerTable) (string, error) {
	t, err := template.New("delete").Parse(sqltemplate.Delete)
	if err != nil {
		return "", nil
	}
	deleteBuffer := new(bytes.Buffer)
	keys := make([]string, 0)
	keyValues := make([]string, 0)
	for snake, key := range table.CacheKey {
		if snake == table.PrimaryField.SnakeCase {
			keys = append(keys, key.Key)
		} else {
			keys = append(keys, key.DataKey)
		}
		keyValues = append(keyValues, key.KeyVariable)
	}
	var isOnlyPrimaryKeyCache = true
	for _, item := range table.Fields {
		if item.IsPrimaryKey {
			continue
		}
		if item.Cache {
			isOnlyPrimaryKeyCache = false
			break
		}
	}
	err = t.Execute(deleteBuffer, map[string]interface{}{
		"upperObject":     table.UpperCamelCase,
		"containsCache":   table.ContainsCache,
		"isNotPrimaryKey": !isOnlyPrimaryKeyCache,
		"lowerPrimaryKey": table.PrimaryField.LowerCamelCase,
		"dataType":        table.PrimaryField.DataType,
		"keys":            strings.Join(keys, "\r\n"),
		"snakePrimaryKey": table.PrimaryField.SnakeCase,
		"keyValues":       strings.Join(keyValues, ", "),
	})
	if err != nil {
		return "", err
	}
	return deleteBuffer.String(), nil
}
