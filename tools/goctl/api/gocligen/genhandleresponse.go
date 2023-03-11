package gocligen

import (
	_ "embed"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
)

const (
	handleResponsePackage = "svc"
	handleResponseFile    = "handleresponse"
)

//go:embed handleresponse.tpl
var handleResponseTemplate string

func genHandleResponse(dir string, cfg *config.Config) error {
	handleFilename, err := format.FileNamingFormat(cfg.NamingFormat, handleResponseFile)
	if err != nil {
		return err
	}

	handleFilename = handleFilename + ".go"
	filename := path.Join(dir, handleResponseDir, handleFilename)
	if _, err := os.Stat(filename); os.IsExist(err) {
		return nil
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          handleResponseDir,
		filename:        handleFilename,
		templateName:    "handleResponseTemplate",
		category:        category,
		templateFile:    handleResponseTemplateFile,
		builtinTemplate: handleResponseTemplate,
		data: map[string]interface{}{
			"pkgName": handleResponsePackage,
			"imports": genHandleResponseImports(),
		},
	})
}

func genHandleResponseImports() string {
	var imports []string
	imports = append(imports, `"encoding/json"`)
	imports = append(imports, `"fmt"`)
	imports = append(imports, `"io/ioutil"`)
	imports = append(imports, `"net/http"`)
	imports = append(imports, `"reflect"`)
	return strings.Join(imports, "\n\t")
}
