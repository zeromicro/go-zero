package gen

import (
	"bytes"
	"text/template"

	sqltemplate "zero/tools/goctl/model/sql/template"
)

func genTypes(table *InnerTable) (string, error) {
	fields := table.Fields
	t, err := template.New("types").Parse(sqltemplate.Types)
	if err != nil {
		return "", nil
	}
	var typeBuffer = new(bytes.Buffer)
	fieldsString, err := genFields(fields)
	if err != nil {
		return "", err
	}
	err = t.Execute(typeBuffer, map[string]interface{}{
		"upperObject":   table.UpperCamelCase,
		"containsCache": table.ContainsCache,
		"fields":        fieldsString,
	})
	if err != nil {
		return "", err
	}
	return typeBuffer.String(), nil
}
