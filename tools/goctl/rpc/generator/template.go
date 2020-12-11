package generator

import (
	"fmt"

	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const (
	category                          = "rpc"
	callTemplateFile                  = "call.tpl"
	callInterfaceFunctionTemplateFile = "call-interface-func.tpl"
	callFunctionTemplateFile          = "call-func.tpl"
	configTemplateFileFile            = "config.tpl"
	etcTemplateFileFile               = "etc.tpl"
	logicTemplateFileFile             = "logic.tpl"
	logicFuncTemplateFileFile         = "logic-func.tpl"
	mainTemplateFile                  = "main.tpl"
	serverTemplateFile                = "server.tpl"
	serverFuncTemplateFile            = "server-func.tpl"
	svcTemplateFile                   = "svc.tpl"
	rpcTemplateFile                   = "template.tpl"
)

var templates = map[string]string{
	callTemplateFile:                  callTemplateText,
	callInterfaceFunctionTemplateFile: callInterfaceFunctionTemplate,
	callFunctionTemplateFile:          callFunctionTemplate,
	configTemplateFileFile:            configTemplate,
	etcTemplateFileFile:               etcTemplate,
	logicTemplateFileFile:             logicTemplate,
	logicFuncTemplateFileFile:         logicFunctionTemplate,
	mainTemplateFile:                  mainTemplate,
	serverTemplateFile:                serverTemplate,
	serverFuncTemplateFile:            functionTemplate,
	svcTemplateFile:                   svcTemplate,
	rpcTemplateFile:                   rpcTemplateText,
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

func Clean() error {
	return util.Clean(category)
}

func Update() error {
	err := Clean()
	if err != nil {
		return err
	}

	return util.InitTemplates(category, templates)
}

func Category() string {
	return category
}
