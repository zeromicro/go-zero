package gen

import (
	"bytes"
	"strings"
	"text/template"

	sqltemplate "github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
)

func genFields(fields []*InnerField) (string, error) {
	list := make([]string, 0)
	for _, field := range fields {
		result, err := genField(field)
		if err != nil {
			return "", err
		}
		list = append(list, result)
	}
	return strings.Join(list, "\r\n"), nil
}

func genField(field *InnerField) (string, error) {
	t, err := template.New("types").Parse(sqltemplate.Field)
	if err != nil {
		return "", nil
	}
	var typeBuffer = new(bytes.Buffer)
	err = t.Execute(typeBuffer, map[string]string{
		"name":    field.UpperCamelCase,
		"type":    field.DataType,
		"tag":     field.Tag,
		"comment": field.Comment,
	})
	if err != nil {
		return "", err
	}
	return typeBuffer.String(), nil
}
