package cgen

import (
	"errors"
	"fmt"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/api/cgen/gen"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	// VarStringDir describes a directory.
	VarStringDir string
	// VarStringAPI describes an API file.
	VarStringAPI string
)

func CCommand(_ *cobra.Command, _ []string) error {
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
	gen.AllTypes = api.Types
	logx.Must(pathx.MkdirIfNotExist(dir))
	logx.Must(gen.GenMessages(dir, api))
	logx.Must(gen.GenClient(dir, api))

	fmt.Println(color.Green.Render("Done."))
	return nil
}
