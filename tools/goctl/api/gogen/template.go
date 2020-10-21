package gogen

import (
	"fmt"

	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
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

func GenTemplates(_ *cli.Context) error {
	return util.InitTemplates(category, templates)
}

func RevertTemplate(name string) error {
	content, ok := templates[name]
	if !ok {
		return fmt.Errorf("%s: no such file name", name)
	}
	return util.CreateTemplate(category, name, content)
}

func Update(category string) error {
	err := Clean()
	if err != nil {
		return err
	}
	return util.InitTemplates(category, templates)
}

func Clean() error {
	return util.Clean(category)
}

func GetCategory() string {
	return category
}
