package rabbitmqgen

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"

	"github.com/lerity-yao/go-zero/tools/cztctl/api/gogen"
	"github.com/lerity-yao/go-zero/tools/cztctl/api/spec"
	"github.com/lerity-yao/go-zero/tools/cztctl/config"
	"github.com/lerity-yao/go-zero/tools/cztctl/util/format"
)

const (
	defaultPort = 8888
	etcDir      = "etc"
)

//go:embed etc.tpl
var etcTemplate string

func genEtc(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, api.Service.Name)
	if err != nil {
		return err
	}

	service := api.Service
	host := "0.0.0.0"
	port := strconv.Itoa(defaultPort)
	rabbitmqNames := generateRabbitmqEtcNames(api)
	return gogen.GenFile(gogen.FileGenConfig{
		Dir:             dir,
		Subdir:          etcDir,
		Filename:        fmt.Sprintf("%s.yaml", filename),
		TemplateName:    "etcTemplate",
		Category:        category,
		TemplateFile:    etcTemplateFile,
		BuiltinTemplate: etcTemplate,
		Data: map[string]string{
			"serviceName":    service.Name,
			"host":           host,
			"port":           port,
			"rabbitmqConfig": strings.Join(rabbitmqNames, "\n"),
		},
	})
}
