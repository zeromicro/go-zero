package tsgen

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringDir describes a directory.
	VarStringDir string
	// VarStringAPI describes an API file.
	VarStringAPI string
	// VarStringWebAPI describes a web API file.
	VarStringWebAPI string
	// VarStringCaller describes a caller.
	VarStringCaller string
	// VarBoolUnWrap describes whether wrap or not.
	VarBoolUnWrap bool
)

// TsCommand provides the entry to generate typescript codes
func TsCommand(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI
	dir := VarStringDir
	webAPI := VarStringWebAPI
	caller := VarStringCaller
	unwrapAPI := VarBoolUnWrap
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}

	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	if len(webAPI) == 0 {
		webAPI = "."
	}

	api, err := parser.Parse(apiFile)
	if err != nil {
		fmt.Println(color.Red.Render("Failed"))
		return err
	}

	if err := api.Validate(); err != nil {
		return err
	}

	api.Service = api.Service.JoinPrefix()
	logx.Must(pathx.MkdirIfNotExist(dir))
	logx.Must(genRequest(dir))
	logx.Must(genHandler(dir, webAPI, caller, api, unwrapAPI))
	logx.Must(genComponents(dir, api))

	fmt.Println(color.Green.Render("Done."))
	return nil
}
