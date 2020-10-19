package tpl

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/tools/goctl/api/gogen"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const templateParentPath = "/"

func GenTemplates(ctx *cli.Context) error {
	if err := errorx.Chain(
		func() error {
			return gogen.GenTemplates(ctx)
		},
	); err != nil {
		return err
	}

	dir, err := util.GetTemplateDir(templateParentPath)
	if err != nil {
		return err
	}

	fmt.Printf("Templates are generated in %s, %s\n", aurora.Green(dir),
		aurora.Red("edit on your risk!"))

	return nil
}
