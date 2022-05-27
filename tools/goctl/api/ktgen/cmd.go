package ktgen

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
)

var (
	// VarStringDir describes a directory.
	VarStringDir string
	// VarStringAPI describes an API.
	VarStringAPI string
	// VarStringPKG describes a package.
	VarStringPKG string
)

// KtCommand generates kotlin code command entrance
func KtCommand(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI
	if apiFile == "" {
		return errors.New("missing -api")
	}
	dir := VarStringDir
	if dir == "" {
		return errors.New("missing -dir")
	}
	pkg := VarStringPKG
	if pkg == "" {
		return errors.New("missing -pkg")
	}

	api, e := parser.Parse(apiFile)
	if e != nil {
		return e
	}

	if err := api.Validate(); err != nil {
		return err
	}

	api.Service = api.Service.JoinPrefix()
	e = genBase(dir, pkg, api)
	if e != nil {
		return e
	}
	e = genApi(dir, pkg, api)
	if e != nil {
		return e
	}
	return nil
}
