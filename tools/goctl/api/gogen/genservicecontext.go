package gogen

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"zero/tools/goctl/api/spec"
	"zero/tools/goctl/api/util"
)

const (
	contextFilename = "servicecontext.go"
	contextTemplate = `package svc

import {{.configImport}}

type ServiceContext struct {
	Config {{.config}}
}

func NewServiceContext(config {{.config}}) *ServiceContext {
	return &ServiceContext{Config: config}
}

`
)

func genServiceContext(dir string, api *spec.ApiSpec) error {
	fp, created, err := util.MaybeCreateFile(dir, contextDir, contextFilename)
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	var authNames = getAuths(api)
	var auths []string
	for _, item := range authNames {
		auths = append(auths, fmt.Sprintf("%s config.AuthConfig", item))
	}

	parentPkg, err := getParentPackage(dir)
	if err != nil {
		return err
	}
	var configImport = "\"" + path.Join(parentPkg, configDir) + "\""
	t := template.Must(template.New("contextTemplate").Parse(contextTemplate))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]string{
		"configImport": configImport,
		"config":       "config.Config",
	})
	if err != nil {
		return nil
	}
	formatCode := formatCode(buffer.String())
	_, err = fp.WriteString(formatCode)
	return err
}
