package gogen

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/templatex"
	ctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const (
	contextFilename = "servicecontext.go"
	contextTemplate = `package svc

import (
	{{.configImport}}
)

type ServiceContext struct {
	Config {{.config}}
	{{.middleware}}
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

	text, err := templatex.LoadTemplate(category, contextTemplateFile, contextTemplate)
	if err != nil {
		return err
	}

	var middlewareStr string
	for _, item := range getMiddleware(api) {
		middlewareStr += fmt.Sprintf("%s rest.Middleware\n", item)
	}

	var configImport = "\"" + ctlutil.JoinPackages(parentPkg, configDir) + "\""
	if len(middlewareStr) > 0 {
		configImport += fmt.Sprintf("\n\"%s/rest\"", vars.ProjectOpenSourceUrl)
	}

	t := template.Must(template.New("contextTemplate").Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]string{
		"configImport": configImport,
		"config":       "config.Config",
		"middleware":   middlewareStr,
	})
	if err != nil {
		return nil
	}
	formatCode := formatCode(buffer.String())
	_, err = fp.WriteString(formatCode)
	return err
}
