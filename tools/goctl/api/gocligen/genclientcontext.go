package gocligen

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

const clientContextFile = "svc"

//go:embed clientcontext.tpl
var clientContextTemplate string

func genClientContext(dir string, cfg *config.Config, api *spec.ApiSpec) error {
	clientContextFilename, err := format.FileNamingFormat(cfg.NamingFormat, clientContextFile)
	if err != nil {
		return err
	}

	clientContextFilename = clientContextFilename + ".go"
	filename := path.Join(dir, clientContextDir, clientContextFilename)
	if _, err := os.Stat(filename); os.IsExist(err) {
		return nil
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          clientContextDir,
		filename:        clientContextFilename,
		templateName:    "clientContextTemplate",
		category:        category,
		templateFile:    clientContextTemplateFile,
		builtinTemplate: clientContextTemplate,
		data: map[string]interface{}{
			"pkgName": clientContextFile,
			"imports": genClientContextImports(),
			"service": api.Service.Name,
		},
	})
}

func genClientContextImports() string {
	var imports []string
	imports = append(imports, `"errors"`)
	imports = append(imports, `"net/http"`)
	imports = append(imports, `"net/url"`)
	imports = append(imports, `"reflect"`+"\n")
	imports = append(imports, fmt.Sprintf("\"%s/rest/httpc\"", vars.ProjectOpenSourceURL))
	return strings.Join(imports, "\n\t")
}
