package gen

import (
	"fmt"
	"path/filepath"

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

func (g *defaultRpcGenerator) genSvc() error {
	svcPath := g.dirM[dirSvc]
	fileName := filepath.Join(svcPath, fileServiceContext)
	text, err := util.LoadTemplate(category, svcTemplateFile, svcTemplate)
	if err != nil {
		return err
	}

	return util.With("svc").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"imports": fmt.Sprintf(`"%v"`, g.mustGetPackage(dirConfig)),
	}, fileName, false)
}
