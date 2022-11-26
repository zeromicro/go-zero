package generator

import (
	_ "embed"
	"path/filepath"

	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed makefile.tpl
var makefileTemplate string

// GenMakefile generates the Makefile file, which is for quick command
func (g *Generator) GenMakefile(ctx DirContext, _ parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetMain()

	fileName := filepath.Join(dir.Filename, "Makefile")
	text, err := pathx.LoadTemplate(category, makefileTemplate, makefileTemplate)
	if err != nil {
		return err
	}

	return util.With("makefile").Parse(text).SaveTo(map[string]interface{}{
		"serviceName": ctx.GetServiceName().Lower(),
		"isEnt":       c.Ent,
	}, fileName, false)
}
