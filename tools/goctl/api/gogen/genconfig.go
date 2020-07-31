package gogen

import (
	"bytes"
	"text/template"

	"zero/tools/goctl/api/spec"
	"zero/tools/goctl/api/util"
)

const (
	configFile     = "config.go"
	configTemplate = `package config

import (
	"zero/rest"
	{{.authImport}}
)

type Config struct {
	rest.RestConf
}
`
)

func genConfig(dir string, api *spec.ApiSpec) error {
	fp, created, err := util.MaybeCreateFile(dir, configDir, configFile)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	var authImportStr = ""
	t := template.Must(template.New("configTemplate").Parse(configTemplate))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]string{
		"authImport": authImportStr,
	})
	if err != nil {
		return nil
	}
	formatCode := formatCode(buffer.String())
	_, err = fp.WriteString(formatCode)
	return err
}
