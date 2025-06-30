package gogen

import (
	_ "embed"
	"fmt"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"path"
	"strconv"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

const (
	defaultPort = 8888
	etcDir      = "etc"
)

//go:embed etc.tpl
var etcTemplate string

func genEtc(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	baseName, err := format.FileNamingFormat(cfg.NamingFormat, api.Service.Name)
	if err != nil {
		return err
	}

	service := api.Service
	host := "0.0.0.0"
	port := strconv.Itoa(defaultPort)

	etcTypes := []string{
		"yaml",
		"json",
		"yml",
		"toml",
	}
	filename := baseName + "." + etcTypes[0]
	for _, etcType := range etcTypes {
		currentFilename := fmt.Sprintf("%s.%s", baseName, etcType)
		fpath := path.Join(dir, etcDir, currentFilename)
		if pathx.FileExists(fpath) {
			filename = currentFilename
			break
		}
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          etcDir,
		filename:        filename,
		templateName:    "etcTemplate",
		category:        category,
		templateFile:    etcTemplateFile,
		builtinTemplate: etcTemplate,
		data: map[string]string{
			"serviceName": service.Name,
			"host":        host,
			"port":        port,
		},
	})
}
