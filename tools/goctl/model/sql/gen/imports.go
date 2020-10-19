package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/templatex"
)

func genImports(withCache, timeImport bool) (string, error) {
	if withCache {
		text, err := templatex.LoadTemplate(category, importsTemplateFile, template.Imports)
		if err != nil {
			return "", err
		}
		buffer, err := templatex.With("import").Parse(text).Execute(map[string]interface{}{
			"time": timeImport,
		})
		if err != nil {
			return "", err
		}
		return buffer.String(), nil
	} else {
		text, err := templatex.LoadTemplate(category, importsWithNoCacheTemplateFile, template.ImportsNoCache)
		if err != nil {
			return "", err
		}
		buffer, err := templatex.With("import").Parse(text).Execute(map[string]interface{}{
			"time": timeImport,
		})
		if err != nil {
			return "", err
		}
		return buffer.String(), nil
	}
}
