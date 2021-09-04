package tpl

import (
	"fmt"
	"path/filepath"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/core/errorx"
	"github.com/tal-tech/go-zero/tools/goctl/api/gogen"
	"github.com/tal-tech/go-zero/tools/goctl/docker"
	internalErrorx "github.com/tal-tech/go-zero/tools/goctl/internal/errorx"
	"github.com/tal-tech/go-zero/tools/goctl/kube"
	mongogen "github.com/tal-tech/go-zero/tools/goctl/model/mongo/generate"
	modelgen "github.com/tal-tech/go-zero/tools/goctl/model/sql/gen"
	rpcgen "github.com/tal-tech/go-zero/tools/goctl/rpc/generator"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/urfave/cli"
)

const templateParentPath = "/"

// GenTemplates writes the latest template text into file which is not exists
func GenTemplates(ctx *cli.Context) error {
	path := ctx.String("home")
	if len(path) != 0 {
		util.RegisterGoctlHome(path)
	}

	internalErrorx.Must(errorx.Chain(
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
		func() error {
			return mongogen.Templates(ctx)
		},
	))

	dir, err := util.GetTemplateDir(templateParentPath)
	internalErrorx.Must(err)

	abs, err := filepath.Abs(dir)
	internalErrorx.Must(err)

	fmt.Printf("Templates are generated in %s, %s\n", aurora.Green(abs),
		aurora.Red("edit on your risk!"))

	return nil
}

// CleanTemplates deletes all templates
func CleanTemplates(ctx *cli.Context) error {
	path := ctx.String("home")
	if len(path) != 0 {
		util.RegisterGoctlHome(path)
	}

	internalErrorx.Must(errorx.Chain(
		func() error {
			return gogen.Clean()
		},
		func() error {
			return modelgen.Clean()
		},
		func() error {
			return rpcgen.Clean()
		},
		func() error {
			return docker.Clean()
		},
		func() error {
			return kube.Clean()
		},
		func() error {
			return mongogen.Clean()
		},
	))

	fmt.Printf("%s\n", aurora.Green("template are clean!"))
	return nil
}

// UpdateTemplates writes the latest template text into file,
// it will delete the older templates if there are exists
func UpdateTemplates(ctx *cli.Context) (err error) {
	path := ctx.String("home")
	category := ctx.String("category")
	if len(path) != 0 {
		util.RegisterGoctlHome(path)
	}

	defer func() {
		if err == nil {
			fmt.Println(aurora.Green(fmt.Sprintf("%s template are update!", category)).String())
		}
	}()
	switch category {
	case docker.Category():
		internalErrorx.Must(docker.Update())
	case gogen.Category():
		internalErrorx.Must(gogen.Update())
	case kube.Category():
		internalErrorx.Must(kube.Update())
	case rpcgen.Category():
		internalErrorx.Must(rpcgen.Update())
	case modelgen.Category():
		internalErrorx.Must(modelgen.Update())
	case mongogen.Category():
		internalErrorx.Must(mongogen.Update())
	default:
		internalErrorx.Must(fmt.Errorf("unexpected category: %s", category))
	}

	return nil
}

// RevertTemplates will overwrite the old template content with the new template
func RevertTemplates(ctx *cli.Context) (err error) {
	path := ctx.String("home")
	category := ctx.String("category")
	filename := ctx.String("name")
	if len(path) != 0 {
		util.RegisterGoctlHome(path)
	}

	defer func() {
		if err == nil {
			fmt.Println(aurora.Green(fmt.Sprintf("%s template are reverted!", filename)).String())
		}
	}()
	switch category {
	case docker.Category():
		internalErrorx.Must(docker.RevertTemplate(filename))
	case kube.Category():
		internalErrorx.Must(kube.RevertTemplate(filename))
	case gogen.Category():
		internalErrorx.Must(gogen.RevertTemplate(filename))
	case rpcgen.Category():
		internalErrorx.Must(rpcgen.RevertTemplate(filename))
	case modelgen.Category():
		internalErrorx.Must(modelgen.RevertTemplate(filename))
	case mongogen.Category():
		internalErrorx.Must(mongogen.RevertTemplate(filename))
	default:
		internalErrorx.Must(fmt.Errorf("unexpected category: %s", category))
	}

	return nil
}
