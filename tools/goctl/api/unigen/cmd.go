package unigen

import (
	"errors"

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
	// VarStringAPI describes an PHP namespace.
	VarStringNS string
)

func UniAppCommand(_ *cobra.Command, _ []string) error {
	apiFile := VarStringAPI
	if apiFile == "" {
		return errors.New("missing -api")
	}
	dir := VarStringDir
	if dir == "" {
		return errors.New("missing -dir")
	}

	api, e := parser.Parse(apiFile)
	if e != nil {
		return e
	}

	if err := api.Validate(); err != nil {
		return err
	}

	logx.Must(pathx.MkdirIfNotExist(dir))
	logx.Must(genMessages(dir, api))
	logx.Must(genClient(dir, api))

	return nil
}
