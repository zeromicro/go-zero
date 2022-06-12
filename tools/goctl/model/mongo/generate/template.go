package generate

import (
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/model/mongo/template"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	category                = "mongo"
	modelTemplateFile       = "model.tpl"
	modelCustomTemplateFile = "model_custom.tpl"
	modelTypesTemplateFile  = "model_types.tpl"
	errTemplateFile         = "err.tpl"
)

var templates = map[string]string{
	modelTemplateFile:       template.ModelText,
	modelCustomTemplateFile: template.ModelCustomText,
	modelTypesTemplateFile:  template.ModelTypesText,
	errTemplateFile:         template.Error,
}

// Category returns the mongo category.
func Category() string {
	return category
}

// Clean cleans the mongo templates.
func Clean() error {
	return pathx.Clean(category)
}

// Templates initializes the mongo templates.
func Templates() error {
	return pathx.InitTemplates(category, templates)
}

// RevertTemplate reverts the given template.
func RevertTemplate(name string) error {
	content, ok := templates[name]
	if !ok {
		return fmt.Errorf("%s: no such file name", name)
	}

	return pathx.CreateTemplate(category, name, content)
}

// Update cleans and updates the templates.
func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return pathx.InitTemplates(category, templates)
}
