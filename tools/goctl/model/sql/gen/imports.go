package gen

import (
	"bytes"
	"text/template"

	sqltemplate "github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
)

func genImports(table *InnerTable) (string, error) {
	t, err := template.New("imports").Parse(sqltemplate.Imports)
	if err != nil {
		return "", err
	}
	importBuffer := new(bytes.Buffer)
	err = t.Execute(importBuffer, map[string]interface{}{
		"containsCache": table.ContainsCache,
	})
	if err != nil {
		return "", err
	}
	return importBuffer.String(), nil
}
