package gogen

import (
	"github.com/tal-tech/go-zero/tools/goctl/templatex"
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
	return templatex.InitTemplates(category, templates)
}
