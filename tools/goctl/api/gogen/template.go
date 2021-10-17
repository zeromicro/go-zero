package gogen

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

const (
	category            = "api"
	configTemplateFile  = "config.tpl"
	contextTemplateFile = "context.tpl"
	etcTemplateFile     = "etc.tpl"
	handlerTemplateFile = "handler.tpl"
	logicTemplateFile   = "logic.tpl"
	mainTemplateFile    = "main.tpl"
)

var templates = map[string]string{
	configTemplateFile:  configTemplate,
	contextTemplateFile: contextTemplate,
	etcTemplateFile:     etcTemplate,
	handlerTemplateFile: handlerTemplate,
	logicTemplateFile:   logicTemplate,
	mainTemplateFile:    mainTemplate,
}

// Category returns the category of the api files.
func Category() string {
	return category
}

// Clean cleans the generated deployment files.
func Clean() error {
	return util.Clean(category)
}

// GenTemplates generates api template files.
func GenTemplates(_ *cli.Context) error {
	return util.InitTemplates(category, templates)
}

// RevertTemplate reverts the given template file to the default value.
func RevertTemplate(name string) error {
	content, ok := templates[name]
	if !ok {
		return fmt.Errorf("%s: no such file name", name)
	}
	return util.CreateTemplate(category, name, content)
}

// Update updates the template files to the templates built in current goctl.
func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return util.InitTemplates(category, templates)
}
