package sqlmodel

import (
	"bytes"
	"strings"
	"text/template"
)

var queryOneTemplate = `
func ({{.pointer}} *{{.model}}Model) FindOne({{.conditions}}) (*{{.model}}, error) {
	sql :=` + " `" + "select " + "` +" + ` {{.modelWithLowerStart}}Rows + ` + "`" + ` from ` + "` + " + `{{.pointer}}.table +` + " ` " + `where {{.valueConditions}} limit 1` + "`" + `
	var resp {{.model}}
	err := {{.pointer}}.conn.QueryRow(&resp, sql, {{.values}})
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &resp, nil
}
`

var queryListTemplate = `
func ({{.pointer}} *{{.model}}Model) Find({{.conditions}}) ([]{{.model}}, error) {
	sql :=` + " `" + "select " + "` +" + ` {{.modelWithLowerStart}}Rows + ` + "`" + ` from ` + "` + " + `{{.pointer}}.table +` + " ` " + `where {{.valueConditions}}` + "`" + `
	var resp []{{.model}}
	err := {{.pointer}}.conn.QueryRows(&resp, sql, {{.values}})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
`

func (s *structExp) genQueryOne() (string, error) {
	conditionExp, valueConditions := s.buildCondition()

	t := template.Must(template.New("queryOneTemplate").Parse(queryOneTemplate))
	var tmplBytes bytes.Buffer
	err := t.Execute(&tmplBytes, map[string]string{
		"pointer":             "m",
		"model":               s.name,
		"conditions":          conditionExp,
		"valueConditions":     valueConditions,
		"values":              strings.Join(s.conditionNames(), ", "),
		"modelWithLowerStart": fmtUnderLine2Camel(s.name, false),
	})
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}

func (s *structExp) genQueryList() (string, error) {
	if len(s.conditions) == 1 && s.conditions[0] == s.primaryKey {
		return "", nil
	}
	conditionExp, valueConditions := s.buildCondition()

	t := template.Must(template.New("queryListTemplate").Parse(queryListTemplate))
	var tmplBytes bytes.Buffer
	err := t.Execute(&tmplBytes, map[string]string{
		"pointer":             "m",
		"model":               s.name,
		"conditions":          conditionExp,
		"valueConditions":     valueConditions,
		"values":              strings.Join(s.conditionNames(), ", "),
		"modelWithLowerStart": fmtUnderLine2Camel(s.name, false),
	})
	if err != nil {
		return "", err
	}
	return tmplBytes.String(), nil
}
