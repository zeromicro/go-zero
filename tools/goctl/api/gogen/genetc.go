package gogen

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"

	"github.com/tal-tech/go-zero/tools/goctl/api/spec"
	"github.com/tal-tech/go-zero/tools/goctl/api/util"
	"github.com/tal-tech/go-zero/tools/goctl/config"
	ctlutil "github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/format"
)

const (
	defaultPort = 8888
	etcDir      = "etc"
	etcTemplate = `Name: {{.serviceName}}
Host: {{.host}}
Port: {{.port}}
`
)

func genEtc(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, api.Service.Name)
	if err != nil {
		return err
	}

	fp, created, err := util.MaybeCreateFile(dir, etcDir, fmt.Sprintf("%s.yaml", filename))
	if err != nil {
		return err
	}
	if !created {
		return nil
	}
	defer fp.Close()

	service := api.Service
	host, ok := util.GetAnnotationValue(service.Groups[0].Annotations, "server", "host")
	if !ok {
		host = "0.0.0.0"
	}
	port, ok := util.GetAnnotationValue(service.Groups[0].Annotations, "server", "port")
	if !ok {
		port = strconv.Itoa(defaultPort)
	}

	text, err := ctlutil.LoadTemplate(category, etcTemplateFile, etcTemplate)
	if err != nil {
		return err
	}

	t := template.Must(template.New("etcTemplate").Parse(text))
	buffer := new(bytes.Buffer)
	err = t.Execute(buffer, map[string]string{
		"serviceName": service.Name,
		"host":        host,
		"port":        port,
	})
	if err != nil {
		return err
	}

	formatCode := formatCode(buffer.String())
	_, err = fp.WriteString(formatCode)
	return err
}
