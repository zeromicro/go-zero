package initlogic

import (
	_ "embed"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/tools/goctl/util/console"
)

//go:embed other.tpl
var otherTpl string

func OtherGen(g *CoreGenContext) error {
	var otherString strings.Builder
	otherTemplate, err := template.New("init_other").Parse(otherTpl)
	if err != nil {
		return errors.Wrap(err, "failed to create other init template")
	}

	err = otherTemplate.Execute(&otherString, map[string]any{
		"modelName":      g.ModelName,
		"modelNameSnake": strcase.ToSnake(g.ModelName),
		"modelNameUpper": strings.ToUpper(g.ModelName),
	})
	if err != nil {
		return err
	}

	console.Info(otherString.String())

	return err
}
