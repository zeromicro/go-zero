package generator

import (
	_ "embed"
	"path/filepath"

	conf "github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/rpc/parser"
	"github.com/zeromicro/go-zero/tools/goctl/util"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed gitlab.tpl
var gitlabTemplate string

// GenGitlab generates the Gitlab-ci.yml file, which is for CI/CD
func (g *Generator) GenGitlab(ctx DirContext, _ parser.Proto, cfg *conf.Config, c *ZRpcContext) error {
	dir := ctx.GetMain()

	fileName := filepath.Join(dir.Filename, ".gitlab-ci.yml")
	text, err := pathx.LoadTemplate(category, gitlabTemplateFile, gitlabTemplate)
	if err != nil {
		return err
	}

	return util.With("gitlab").Parse(text).SaveTo(map[string]any{}, fileName, false)
}
