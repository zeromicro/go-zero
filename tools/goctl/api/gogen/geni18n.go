package gogen

import (
	_ "embed"
)

//go:embed i18n/var.tpl
var i18nVarTemplate string

//go:embed i18n/locale/zh.json
var zhLocaleFile string

//go:embed i18n/locale/en.json
var enLocaleFile string

func genI18n(dir string) error {
	err := genFile(fileGenConfig{
		dir:             dir,
		subdir:          localeDir,
		filename:        "zh.json",
		templateName:    "zhLocaleTemplate",
		category:        category,
		templateFile:    "",
		builtinTemplate: zhLocaleFile,
		data:            map[string]string{},
	})

	if err != nil {
		return err
	}

	err = genFile(fileGenConfig{
		dir:             dir,
		subdir:          localeDir,
		filename:        "en.json",
		templateName:    "enLocaleTemplate",
		category:        category,
		templateFile:    "",
		builtinTemplate: enLocaleFile,
		data:            map[string]string{},
	})

	if err != nil {
		return err
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          i18nDir,
		filename:        "vars.go",
		templateName:    "i18nVarTemplate",
		category:        category,
		templateFile:    "",
		builtinTemplate: i18nVarTemplate,
		data:            map[string]string{},
	})
}
