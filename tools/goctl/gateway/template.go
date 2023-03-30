package gateway

import (
	_ "embed"
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const (
	category            = "gateway"
	etcTemplateFileFile = "etc.tpl"
	mainTemplateFile    = "main.tpl"
)

//go:embed conf.yml
var etcTemplate string

//go:embed gateway.tpl
var mainTemplate string

var templates = map[string]string{
	etcTemplateFileFile: etcTemplate,
	mainTemplateFile:    mainTemplate,
}

// GenTemplates is the entry for command goctl template,
// it will create the specified category
func GenTemplates() error {
	return pathx.InitTemplates(category, templates)
}

// RevertTemplate restores the deleted template files
func RevertTemplate(name string) error {
	content, ok := templates[name]
	if !ok {
		return fmt.Errorf("%s: no such file name", name)
	}
	return pathx.CreateTemplate(category, name, content)
}

// Clean deletes all template files
func Clean() error {
	return pathx.Clean(category)
}

// Update is used to update the template files, it will delete the existing old templates at first,
// and then create the latest template files
func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return pathx.InitTemplates(category, templates)
}

// Category returns a const string value for rpc template category
func Category() string {
	return category
}
