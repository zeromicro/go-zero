package swagger

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"gopkg.in/yaml.v2"
)

var (
	// VarStringAPI specifies the API filename.
	VarStringAPI string

	// VarStringDir specifies the directory to generate swagger file.
	VarStringDir string

	// VarStringFilename specifies the generated swagger file name without the extension.
	VarStringFilename string

	// VarBoolYaml specifies whether to generate a YAML file.
	VarBoolYaml bool
)

func Command(_ *cobra.Command, _ []string) error {
	if len(VarStringAPI) == 0 {
		return errors.New("missing -api")
	}

	if len(VarStringDir) == 0 {
		return errors.New("missing -dir")
	}

	api, err := parser.Parse(VarStringAPI, "")
	if err != nil {
		return err
	}

	fillAllStructs(api)

	if err := api.Validate(); err != nil {
		return err
	}
	swagger, err := spec2Swagger(api)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return err
	}

	err = pathx.MkdirIfNotExist(VarStringDir)
	if err != nil {
		return err
	}

	filename := VarStringFilename
	if filename == "" {
		base := filepath.Base(VarStringAPI)
		filename = strings.TrimSuffix(base, filepath.Ext(base))
	}

	if VarBoolYaml {
		filePath := filepath.Join(VarStringDir, filename+".yaml")

		var jsonObj interface{}
		if err := yaml.Unmarshal(data, &jsonObj); err != nil {
			return err
		}

		data, err := yaml.Marshal(jsonObj)
		if err != nil {
			return err
		}
		return os.WriteFile(filePath, data, 0644)
	}

	// generate json swagger file
	filePath := filepath.Join(VarStringDir, filename+".json")
	return os.WriteFile(filePath, data, 0644)
}
