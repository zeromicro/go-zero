package sqlmodel

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/core/stringx"
)

var updateTemplate = `
func ({{.pointer}} *{{.model}}Model) Update(data {{.model}}{{.conditions}}) error {
	sql := ` + "`" + `update ` + "` + " + `{{.pointer}}.table` + " + `" + ` set ` + "` +" + ` {{.modelWithLowerStart}}RowsWithPlaceHolder + ` + "` where" + ` {{.valueConditions}}` + "`\n" +
	`	_, err := {{.pointer}}.conn.Exec(sql, {{.values}})
	return err
}
`

func (s *structExp) genUpdate() (string, error) {
	var updateValues []string
	var conditionsValues []string
	for _, field := range s.Fields {
		key := fmt.Sprintf("data.%s", field.name)
		if stringx.Contains(s.conditions, field.tag) ||
			stringx.Contains(s.conditions, field.name) {
			conditionsValues = append(conditionsValues, key)
		} else if !stringx.Contains(s.ignoreFields, field.name) && !stringx.Contains(s.ignoreFields, field.tag) {
			updateValues = append(updateValues, key)
		}
	}
	conditionExp, valueConditions := s.buildCondition()
	if len(s.conditions) == 1 && s.conditions[0] == s.primaryKey {
		conditionExp = ""
	} else {
		conditionExp = ", " + conditionExp
	}
	t := template.Must(template.New("updateTemplate").Parse(updateTemplate))
	var tmplBytes bytes.Buffer
	err := t.Execute(&tmplBytes, map[string]string{
		"pointer":             "m",
		"model":               s.name,
		"values":              strings.Join(append(updateValues, conditionsValues...), ", "),
		"primaryKey":          s.primaryKey,
		"valueConditions":     valueConditions,
		"conditions":          conditionExp,
		"modelWithLowerStart": fmtUnderLine2Camel(s.name, false),
	})
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}
