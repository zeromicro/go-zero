package generate

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/tools/goctl/model/mongo/template"
	"github.com/zeromicro/go-zero/tools/goctl/util"
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

// Category returns the mongo category.
func Category() string {
	return category
}

// Clean cleans the mongo templates.
func Clean() error {
	return util.Clean(category)
}

// Templates initializes the mongo templates.
func Templates(_ *cli.Context) error {
	return util.InitTemplates(category, templates)
}

// RevertTemplate reverts the given template.
func RevertTemplate(name string) error {
	content, ok := templates[name]
	if !ok {
		return fmt.Errorf("%s: no such file name", name)
	}

	return util.CreateTemplate(category, name, content)
}

// Update cleans and updates the templates.
func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return util.InitTemplates(category, templates)
}
