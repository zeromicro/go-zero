package swagger

import (
	"encoding/json"
	"errors"
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
	filename := strings.TrimSuffix(VarStringAPI, filepath.Ext(VarStringAPI)) + ".json"
	data, err := json.MarshalIndent(swagger, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}
