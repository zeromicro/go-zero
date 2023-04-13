package generator

import (
	_ "embed"
	"path/filepath"

	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed dockerfile.tpl
var dockerfileTemplate string

// GenDockerfile generates the Makefile file, which is for quick command
func (g *Generator) GenDockerfile(ctx DirContext, _ parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetMain()

	fileName := filepath.Join(dir.Filename, "Dockerfile")
	text, err := pathx.LoadTemplate(category, dockerfileTemplateFile, dockerfileTemplate)
	if err != nil {
		return err
	}

	serviceName, err := format.FileNamingFormat(cfg.NamingFormat, ctx.GetServiceName().Source())
	if err != nil {
		return err
	}

	return util.With("dockerfile").Parse(text).SaveTo(map[string]any{
		"serviceName": serviceName,
		"port":        c.Port,
		"imageTag":    "golang:1.20.2-alpine3.17",
	}, fileName, false)
}
