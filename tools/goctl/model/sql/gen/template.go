package gen

import (
	"fmt"

	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const (
	category                              = "model"
	deleteTemplateFile                    = "delete.tpl"
	deleteMethodTemplateFile              = "interface-delete.tpl"
	fieldTemplateFile                     = "filed.tpl"
	findOneTemplateFile                   = "find-one.tpl"
	findOneMethodTemplateFile             = "interface-find-one.tpl"
	findOneByFieldTemplateFile            = "find-one-by-field.tpl"
	findOneByFieldMethodTemplateFile      = "interface-find-one-by-field.tpl"
	findOneByFieldExtraMethodTemplateFile = "find-one-by-filed-extra-method.tpl"
	importsTemplateFile                   = "import.tpl"
	importsWithNoCacheTemplateFile        = "import-no-cache.tpl"
	insertTemplateFile                    = "insert.tpl"
	insertTemplateMethodFile              = "interface-insert.tpl"
	modelTemplateFile                     = "model.tpl"
	modelNewTemplateFile                  = "model-new.tpl"
	tagTemplateFile                       = "tag.tpl"
	typesTemplateFile                     = "types.tpl"
	updateTemplateFile                    = "update.tpl"
	updateMethodTemplateFile              = "interface-update.tpl"
	varTemplateFile                       = "var.tpl"
	errTemplateFile                       = "err.tpl"
)

var templates = map[string]string{
	deleteTemplateFile:                    template.Delete,
	deleteMethodTemplateFile:              template.DeleteMethod,
	fieldTemplateFile:                     template.Field,
	findOneTemplateFile:                   template.FindOne,
	findOneMethodTemplateFile:             template.FindOneMethod,
	findOneByFieldTemplateFile:            template.FindOneByField,
	findOneByFieldMethodTemplateFile:      template.FindOneByFieldMethod,
	findOneByFieldExtraMethodTemplateFile: template.FindOneByFieldExtraMethod,
	importsTemplateFile:                   template.Imports,
	importsWithNoCacheTemplateFile:        template.ImportsNoCache,
	insertTemplateFile:                    template.Insert,
	insertTemplateMethodFile:              template.InsertMethod,
	modelTemplateFile:                     template.Model,
	modelNewTemplateFile:                  template.New,
	tagTemplateFile:                       template.Tag,
	typesTemplateFile:                     template.Types,
	updateTemplateFile:                    template.Update,
	varTemplateFile:                       template.Vars,
	errTemplateFile:                       template.Error,
}

func Category() string {
	return category
}

func Clean() error {
	return util.Clean(category)
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

func Update() error {
	err := Clean()
	if err != nil {
		return err
	}
	return util.InitTemplates(category, templates)
}
