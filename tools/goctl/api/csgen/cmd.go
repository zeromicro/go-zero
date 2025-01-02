package csgen

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
	// VarStringAPI describes an C# namespace.
	VarStringNS string
)

func CSharpCommand(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI
	if apiFile == "" {
		return errors.New("missing -api")
	}
	dir := VarStringDir
	if dir == "" {
		return errors.New("missing -dir")
	}

	ns := VarStringNS
	if ns == "" {
		return errors.New("missing -ns")
	}

	api, e := parser.Parse(apiFile)
	if e != nil {
		return e
	}

	if err := api.Validate(); err != nil {
		return err
	}

	logx.Must(pathx.MkdirIfNotExist(dir))
	logx.Must(genMessages(dir, ns, api))
	logx.Must(genClient(dir, ns, api))

	fmt.Println(color.Green.Render("Done."))
	return nil
}
