package gogen

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/config"
	ctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/format"
	"github.com/tal-tech/go-zero/tools/goctl/vars"
)

const (
	configFile     = "config"
	configTemplate = `package config

import {{.authImport}}

type Config struct {
	rest.RestConf
	{{.auth}}
}
`

	jwtTemplate = ` struct {
		AccessSecret string
		AccessExpire int64
	}
`
)

func genConfig(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, configFile)
	if err != nil {
		return err
	}

	fp, created, err := util.MaybeCreateFile(dir, configDir, filename+".go")
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
		auths = append(auths, fmt.Sprintf("%s %s", item, jwtTemplate))
	}

	var authImportStr = fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceUrl)
	text, err := ctlutil.LoadTemplate(category, configTemplateFile, configTemplate)
	if err != nil {
		return err
	}

	t := template.Must(template.New("configTemplate").Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]string{
		"authImport": authImportStr,
		"auth":       strings.Join(auths, "\n"),
	})
	if err != nil {
		return err
	}

	formatCode := formatCode(buffer.String())
	_, err = fp.WriteString(formatCode)
	return err
}
