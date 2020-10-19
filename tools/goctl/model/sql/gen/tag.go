package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/templatex"
)

func genTag(in string) (string, error) {
	if in == "" {
		return in, nil
	}
	text, err := templatex.LoadTemplate(category, tagTemplateFile, template.Tag)
	if err != nil {
		return "", err
	}
	output, err := templatex.With("tag").
		Parse(text).
		Execute(map[string]interface{}{
			"field": in,
		})
	if err != nil {
		return "", err
	}
	return output.String(), nil
}
