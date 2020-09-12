package gen

import (
	"fmt"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/tools/goctl/rpc/execx"
)

func (g *defaultRpcGenerator) initGoMod() error {
	if !g.Ctx.IsInGoEnv {
		projectDir := g.dirM[dirTarget]
		cmd := fmt.Sprintf("go mod init %s", g.Ctx.ProjectName.Source())
		output, err := execx.Run(fmt.Sprintf(cmd), projectDir)
		if err != nil {
			logx.Error(err)
			return err
		}
		g.Ctx.Info(output)
	}
	return nil
}
