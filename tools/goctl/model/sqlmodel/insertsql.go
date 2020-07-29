package sqlmodel

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"zero/core/stringx"
)

var insertTemplate = `
func ({{.pointer}} *{{.model}}Model) Insert(data {{.model}}) error {
	sql := ` + "`" + `insert into` + " ` + " + `{{.pointer}}.table ` + " + `(` + " + "{{.modelWithLowerStart}}RowsExpectAutoSet" + " + `" + `) value ({{.valueHolder}})` + "`\n" +
	`	_, err := {{.pointer}}.conn.Exec(sql, {{.values}})
	return err
}
`

func (s *structExp) genInsert() (string, error) {
	var valueHolder []string
	var values []string
	for _, field := range s.Fields {
		if stringx.Contains(s.ignoreFields, field.name) || stringx.Contains(s.ignoreFields, field.tag) {
			continue
		}
		valueHolder = append(valueHolder, "?")
		values = append(values, fmt.Sprintf("data.%s", field.name))
	}

	t := template.Must(template.New("insertTemplate").Parse(insertTemplate))
	var tmplBytes bytes.Buffer
	var columns = "rowsExpectAutoSet"
	err := t.Execute(&tmplBytes, map[string]string{
		"pointer":             "m",
		"model":               s.name,
		"columns":             columns,
		"valueHolder":         strings.Join(valueHolder, ","),
		"values":              strings.Join(values, ", "),
		"modelWithLowerStart": fmtUnderLine2Camel(s.name, false),
	})
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}
