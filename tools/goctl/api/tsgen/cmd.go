package tsgen

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/tsgen/gen"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringDir describes a directory.
	VarStringDir string
	// VarStringAPI describes an API file.
	VarStringAPI string
	// VarStringCaller describes a caller.
	VarStringCaller string
	// VarBoolUnWrap describes whether wrap or not.
	VarBoolUnWrap bool
	// VarStringUrlPrefix request url prefix
	VarStringUrlPrefix string
	// VarBoolCustomBody request custom body
	VarBoolCustomBody bool
)

// TsCommand provides the entry to generate typescript codes
func TsCommand(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI
	dir := VarStringDir
	if len(apiFile) == 0 {
		return errors.New("missing -api")
	}

	if len(dir) == 0 {
		return errors.New("missing -dir")
	}

	api, err := parser.Parse(apiFile)
	if err != nil {
		fmt.Println(color.Red.Render("Failed"))
		return err
	}

	if err := api.Validate(); err != nil {
		return err
	}

	caller := VarStringCaller
	if len(caller) == 0 {
		caller = "webapi"
	}

	api.Service = api.Service.JoinPrefix()
	logx.Must(pathx.MkdirIfNotExist(dir))
	logx.Must(gen.GenRequests(dir, caller))
	logx.Must(gen.GenHandler(dir, caller, api, VarBoolUnWrap, VarBoolCustomBody, VarStringUrlPrefix))
	logx.Must(gen.GenComponents(dir, api))

	fmt.Println(color.Green.Render("Done."))
	return nil
}
