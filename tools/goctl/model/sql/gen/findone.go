package gen

import (
	"bytes"
	"text/template"

	sqltemplate "github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
)

func genFindOne(table *InnerTable) (string, error) {
	t, err := template.New("findOne").Parse(sqltemplate.FindOne)
	if err != nil {
		return "", err
	}
	fineOneBuffer := new(bytes.Buffer)
	err = t.Execute(fineOneBuffer, map[string]interface{}{
		"withCache":        table.PrimaryField.Cache,
		"upperObject":      table.UpperCamelCase,
		"lowerObject":      table.LowerCamelCase,
		"snakePrimaryKey":  table.PrimaryField.SnakeCase,
		"lowerPrimaryKey":  table.PrimaryField.LowerCamelCase,
		"dataType":         table.PrimaryField.DataType,
		"cacheKey":         table.CacheKey[table.PrimaryField.SnakeCase].Key,
		"cacheKeyVariable": table.CacheKey[table.PrimaryField.SnakeCase].KeyVariable,
	})
	if err != nil {
		return "", err
	}
	return fineOneBuffer.String(), nil
}
