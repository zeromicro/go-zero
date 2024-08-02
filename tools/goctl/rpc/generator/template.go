package generator

import (
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
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
	callTemplateFile:          callTemplateText,
	configTemplateFileFile:    configTemplate,
	etcTemplateFileFile:       etcTemplate,
	logicTemplateFileFile:     logicTemplate,
	logicFuncTemplateFileFile: logicFunctionTemplate,
	mainTemplateFile:          mainTemplate,
	serverTemplateFile:        serverTemplate,
	serverFuncTemplateFile:    functionTemplate,
	svcTemplateFile:           svcTemplate,
	rpcTemplateFile:           rpcTemplateText,
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
