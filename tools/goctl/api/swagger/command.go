package swagger

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
)

var (
	// VarStringAPI specifies the API filename.
	VarStringAPI string

	// VarStringDir specifies the directory to generate swagger file.
	VarStringDir string

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

	if VarBoolYaml {
		filename := strings.TrimSuffix(VarStringAPI, filepath.Ext(VarStringAPI)) + ".yaml"

		var jsonObj interface{}
		if err := yaml.Unmarshal(data, &jsonObj); err != nil {
			return err
		}

		data, err := yaml.Marshal(jsonObj)
		if err != nil {
			return err
		}
		return os.WriteFile(filename, data, 0644)
	}
	// generate json swagger file
	filename := filepath.Join(VarStringDir, strings.TrimSuffix(VarStringAPI, filepath.Ext(VarStringAPI))+".json")

	return os.WriteFile(filename, data, 0644)
}
