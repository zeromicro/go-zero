package gogen

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const (
	configFile     = "config.go"
	configTemplate = `package config

import {{.authImport}}

type Config struct {
	rest.RestConf
}
`
)

func genConfig(dir string) error {
	fp, created, err := util.MaybeCreateFile(dir, configDir, configFile)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	var authImportStr = fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceUrl)
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
