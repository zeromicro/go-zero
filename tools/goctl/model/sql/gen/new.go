package gen

import (
	"bytes"
	"text/template"

	sqltemplate "github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
)

func genNew(table *InnerTable) (string, error) {
	t, err := template.New("new").Parse(sqltemplate.New)
	if err != nil {
		return "", err
	}
	newBuffer := new(bytes.Buffer)
	err = t.Execute(newBuffer, map[string]interface{}{
		"containsCache": table.ContainsCache,
		"upperObject":   table.UpperCamelCase,
	})
	if err != nil {
		return "", err
	}
	return newBuffer.String(), nil
}
