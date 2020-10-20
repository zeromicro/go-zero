package tpl

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/tools/goctl/api/gogen"
	modelgen "github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	rpcgen "github.com/tal-tech/go-zero/tools/goctl/rpc/gen"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const templateParentPath = "/"

func GenTemplates(ctx *cli.Context) error {
	if err := errorx.Chain(
		func() error {
			return gogen.GenTemplates(ctx)
		},
		func() error {
			return modelgen.GenTemplates(ctx)
		},
		func() error {
			return rpcgen.GenTemplates(ctx)
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

func CleanTemplates(_ *cli.Context) error {
	return errorx.Chain(
		func() error {
			return gogen.Clean()
		},
		func() error {
			return modelgen.Clean()
		},
		func() error {
			return rpcgen.Clean()
		},
	)
}

func UpdateTemplates(ctx *cli.Context) error {
	category := ctx.String("category")
	switch category {
	case gogen.GetCategory():
		return gogen.GenTemplates(ctx)
	case rpcgen.GetCategory():
		return rpcgen.GenTemplates(ctx)
	case modelgen.GetCategory():
		return modelgen.GenTemplates(ctx)
	default:
		return fmt.Errorf("unexpected category: %s", category)
	}
}

func RevertTemplates(ctx *cli.Context) error {
	category := ctx.String("category")
	filename := ctx.String("name")
	switch category {
	case gogen.GetCategory():
		return gogen.RevertTemplate(filename)
	case rpcgen.GetCategory():
		return rpcgen.RevertTemplate(filename)
	case modelgen.GetCategory():
		return modelgen.RevertTemplate(filename)
	default:
		return fmt.Errorf("unexpected category: %s", category)
	}
}
