package gogen

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	ctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
)

const (
	contextFilename = "servicecontext.go"
	contextTemplate = `package svc

import {{.configImport}}

type ServiceContext struct {
	Config {{.config}}
}

func NewServiceContext(c {{.config}}) *ServiceContext {
	return &ServiceContext{Config: c}
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
	var configImport = "\"" + ctlutil.JoinPackages(parentPkg, configDir) + "\""
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
