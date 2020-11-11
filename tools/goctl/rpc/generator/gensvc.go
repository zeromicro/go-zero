package generator

import (
	"fmt"
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const svcTemplate = `package svc

import {{.imports}}

type ServiceContext struct {
	c config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		c:c,
	}
}
`

func (g *defaultGenerator) GenSvc(ctx DirContext, _ parser.Proto) error {
	dir := ctx.GetSvc()
	fileName := filepath.Join(dir.Filename, formatFilename("service_context")+".go")
	text, err := util.LoadTemplate(category, svcTemplateFile, svcTemplate)
	if err != nil {
		return err
	}

	return util.With("svc").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"imports": fmt.Sprintf(`"%v"`, ctx.GetConfig().Package),
	}, fileName, false)
}
