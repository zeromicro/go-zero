package tpl

import (
	"fmt"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/tools/goctl/api/apigen"
	"github.com/zeromicro/go-zero/tools/goctl/api/gogen"
	apinew "github.com/zeromicro/go-zero/tools/goctl/api/new"
	"github.com/zeromicro/go-zero/tools/goctl/docker"
	"github.com/zeromicro/go-zero/tools/goctl/gateway"
	"github.com/zeromicro/go-zero/tools/goctl/kube"
	mongogen "github.com/zeromicro/go-zero/tools/goctl/model/mongo/generate"
	modelgen "github.com/zeromicro/go-zero/tools/goctl/model/sql/gen"
	rpcgen "github.com/zeromicro/go-zero/tools/goctl/rpc/generator"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const templateParentPath = "/"

// genTemplates writes the latest template text into file which is not exists
func genTemplates(_ *cobra.Command, _ []string) error {
	path := varStringHome
	if len(path) != 0 {
		pathx.RegisterGoctlHome(path)
	}

	if err := errorx.Chain(
		func() error {
			return gogen.GenTemplates()
		},
		func() error {
			return modelgen.GenTemplates()
		},
		func() error {
			return rpcgen.GenTemplates()
		},
		func() error {
			return docker.GenTemplates()
		},
		func() error {
			return kube.GenTemplates()
		},
		func() error {
			return mongogen.Templates()
		},
		func() error {
			return apigen.GenTemplates()
		},
		func() error {
			return apinew.GenTemplates()
		},
		func() error {
			return gateway.GenTemplates()
		},
	); err != nil {
		return err
	}

	dir, err := pathx.GetTemplateDir(templateParentPath)
	if err != nil {
		return err
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}

	fmt.Printf("Templates are generated in %s, %s\n", color.Green.Render(abs),
		color.Red.Render("edit on your risk!"))

	return nil
}

// cleanTemplates deletes all templates
func cleanTemplates(_ *cobra.Command, _ []string) error {
	path := varStringHome
	if len(path) != 0 {
		pathx.RegisterGoctlHome(path)
	}

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
		func() error {
			return docker.Clean()
		},
		func() error {
			return kube.Clean()
		},
		func() error {
			return mongogen.Clean()
		},
		func() error {
			return apigen.Clean()
		},
		func() error {
			return apinew.Clean()
		},
		func() error {
			return gateway.Clean()
		},
	)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", color.Green.Render("templates are cleaned!"))
	return nil
}

// updateTemplates writes the latest template text into file,
// it will delete the older templates if there are exists
func updateTemplates(_ *cobra.Command, _ []string) (err error) {
	path := varStringHome
	category := varStringCategory
	if len(path) != 0 {
		pathx.RegisterGoctlHome(path)
	}

	defer func() {
		if err == nil {
			fmt.Println(color.Green.Sprintf("%s template are update!", category))
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
	case mongogen.Category():
		return mongogen.Update()
	case apigen.Category():
		return apigen.Update()
	case apinew.Category():
		return apinew.Update()
	case gateway.Category():
		return gateway.Update()
	default:
		err = fmt.Errorf("unexpected category: %s", category)
		return
	}
}

// revertTemplates will overwrite the old template content with the new template
func revertTemplates(_ *cobra.Command, _ []string) (err error) {
	path := varStringHome
	category := varStringCategory
	filename := varStringName
	if len(path) != 0 {
		pathx.RegisterGoctlHome(path)
	}

	defer func() {
		if err == nil {
			fmt.Println(color.Green.Sprintf("%s template are reverted!", filename))
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
	case mongogen.Category():
		return mongogen.RevertTemplate(filename)
	case apigen.Category():
		return apigen.RevertTemplate(filename)
	case apinew.Category():
		return apinew.RevertTemplate(filename)
	case gateway.Category():
		return gateway.RevertTemplate(filename)
	default:
		err = fmt.Errorf("unexpected category: %s", category)
		return
	}
}
