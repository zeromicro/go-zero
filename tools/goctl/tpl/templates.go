package tpl

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/tools/goctl/api/gogen"
	"github.com/tal-tech/go-zero/tools/goctl/docker"
	"github.com/tal-tech/go-zero/tools/goctl/kube"
	modelgen "github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	rpcgen "github.com/tal-tech/go-zero/tools/goctl/rpc/generator"
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
		func() error {
			return docker.GenTemplates(ctx)
		},
		func() error {
			return kube.GenTemplates(ctx)
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
	err := errorx.Chain(
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
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", aurora.Green("template are clean!"))
	return nil
}

func UpdateTemplates(ctx *cli.Context) (err error) {
	category := ctx.String("category")
	defer func() {
		if err == nil {
			fmt.Println(aurora.Green(fmt.Sprintf("%s template are update!", category)).String())
		}
	}()
	switch category {
	case docker.Category():
		return docker.Update()
	case gogen.Category():
		return gogen.Update()
	case kube.Category():
		return kube.Update()
	case rpcgen.Category():
		return rpcgen.Update()
	case modelgen.Category():
		return modelgen.Update()
	default:
		err = fmt.Errorf("unexpected category: %s", category)
		return
	}
}

func RevertTemplates(ctx *cli.Context) (err error) {
	category := ctx.String("category")
	filename := ctx.String("name")
	defer func() {
		if err == nil {
			fmt.Println(aurora.Green(fmt.Sprintf("%s template are reverted!", filename)).String())
		}
	}()
	switch category {
	case docker.Category():
		return docker.RevertTemplate(filename)
	case kube.Category():
		return kube.RevertTemplate(filename)
	case gogen.Category():
		return gogen.RevertTemplate(filename)
	case rpcgen.Category():
		return rpcgen.RevertTemplate(filename)
	case modelgen.Category():
		return modelgen.RevertTemplate(filename)
	default:
		err = fmt.Errorf("unexpected category: %s", category)
		return
	}
}
