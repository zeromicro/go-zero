package generate

import (
	"fmt"

	"github.com/tal-tech/go-zero/tools/goctl/model/mongo/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const (
	category          = "mongo"
	modelTemplateFile = "model.tpl"
	errTemplateFile   = "err.tpl"
)

var templates = map[string]string{
	modelTemplateFile: template.Text,
	errTemplateFile:   template.Error,
}

func Category() string {
	return category
}

func Clean() error {
	return util.Clean(category)
}

func Templates(_ *cli.Context) error {
	return util.InitTemplates(category, templates)
}

func RevertTemplate(name string) error {
	content, ok := templates[name]
	if !ok {
		return fmt.Errorf("%s: no such file name", name)
	}

	return util.CreateTemplate(category, name, content)
}

func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return util.InitTemplates(category, templates)
}
