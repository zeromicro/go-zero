package gogen

import (
	_ "embed"
	"path/filepath"

	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed dberrorhandler.tpl
var dbErrorHandlerTemplateText string

func GenErrorHandler(dir, style, rootPkg string) error {
	filename, err := format.FileNamingFormat(style, "error_handler.go")
	if err != nil {
		return err
	}

	if err := pathx.MkdirIfNotExist(filepath.Join(dir, utilsDir)); err != nil {
		return err
	}

	if err := pathx.MkdirIfNotExist(filepath.Join(dir, utilsDir, "dberrorhandler")); err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          filepath.Join(utilsDir, "dberrorhandler"),
		filename:        filename,
		templateName:    "dbErrorHandlerFile",
		category:        category,
		templateFile:    dbErrorHandlerTemplateFile,
		builtinTemplate: dbErrorHandlerTemplateText,
		data: map[string]string{
			"package": rootPkg,
		},
	})
}
