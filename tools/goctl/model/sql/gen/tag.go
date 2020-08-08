package gen

import (
	"bytes"
	"text/template"

	sqltemplate "github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
)

func genTag(in string) (string, error) {
	if in == "" {
		return in, nil
	}
	t, err := template.New("tag").Parse(sqltemplate.Tag)
	if err != nil {
		return "", err
	}
	var tagBuffer = new(bytes.Buffer)
	err = t.Execute(tagBuffer, map[string]interface{}{
		"field": in,
	})
	if err != nil {
		return "", err
	}
	return tagBuffer.String(), nil
}
