package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/model/sql/template"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const (
	category                              = "model"
	deleteTemplateFile                    = "delete.tpl"
	fieldTemplateFile                     = "filed.tpl"
	findOneTemplateFile                   = "find-one.tpl"
	findOneByFieldTemplateFile            = "find-one-by-field.tpl"
	findOneByFieldExtraMethodTemplateFile = "find-one-by-filed-extra-method.tpl"
	importsTemplateFile                   = "import.tpl"
	importsWithNoCacheTemplateFile        = "import-no-cache.tpl"
	insertTemplateFile                    = "insert.tpl"
	modelTemplateFile                     = "model.tpl"
	modelNewTemplateFile                  = "model-new.tpl"
	tagTemplateFile                       = "tag.tpl"
	typesTemplateFile                     = "types.tpl"
	updateTemplateFile                    = "update.tpl"
	varTemplateFile                       = "var.tpl"
)

var templates = map[string]string{
	deleteTemplateFile:                    template.Delete,
	fieldTemplateFile:                     template.Field,
	findOneTemplateFile:                   template.FindOne,
	findOneByFieldTemplateFile:            template.FindOneByField,
	findOneByFieldExtraMethodTemplateFile: template.FindOneByFieldExtraMethod,
	importsTemplateFile:                   template.Imports,
	importsWithNoCacheTemplateFile:        template.ImportsNoCache,
	insertTemplateFile:                    template.Insert,
	modelTemplateFile:                     template.Model,
	modelNewTemplateFile:                  template.New,
	tagTemplateFile:                       template.Tag,
	typesTemplateFile:                     template.Types,
	updateTemplateFile:                    template.Update,
	varTemplateFile:                       template.Vars,
}

func GenTemplates(_ *cli.Context) error {
	return util.InitTemplates(category, templates)
}
