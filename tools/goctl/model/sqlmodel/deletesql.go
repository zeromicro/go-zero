package sqlmodel

import (
	"bytes"
	"strings"
	"text/template"
)

var deleteTemplate = `
func ({{.pointer}} *{{.model}}Model) Delete({{.conditions}}) error {
	sql := ` + "`" + `delete from` + " ` + " + `{{.pointer}}.table ` + "+ `" + ` where {{.valueConditions}}` + "`\n" +
	`	_, err := {{.pointer}}.conn.Exec(sql, {{.values}})
	return err
}
`

func (s *structExp) genDelete() (string, error) {
	idType := "string"
	if s.idAutoIncrement {
		idType = "int64"
	}
	conditionExp, valueConditions := s.buildCondition()
	t := template.Must(template.New("deleteTemplate").Parse(deleteTemplate))
	var tmplBytes bytes.Buffer
	err := t.Execute(&tmplBytes, map[string]string{
		"pointer":         "m",
		"model":           s.name,
		"idType":          idType,
		"valueConditions": valueConditions,
		"conditions":      conditionExp,
		"values":          strings.Join(s.conditionNames(), ", "),
	})
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}
