package gogen

import (
	_ "embed"
)

//go:embed gitlab.tpl
var gitlabTemplate string

func genGitlab(dir string) error {
	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "",
		filename:        ".gitlab-ci.yml",
		templateName:    "gitlabTemplate",
		category:        category,
		templateFile:    gitlabTemplateFile,
		builtinTemplate: gitlabTemplate,
		data:            map[string]string{},
	})
}
