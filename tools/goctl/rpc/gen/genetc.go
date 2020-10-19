package gen

import (
	"fmt"
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/util"
)

const etcTemplate = `Name: {{.serviceName}}.rpc
ListenOn: 127.0.0.1:8080
Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: {{.serviceName}}.rpc
`

func (g *defaultRpcGenerator) genEtc() error {
	etdDir := g.dirM[dirEtc]
	fileName := filepath.Join(etdDir, fmt.Sprintf("%v.yaml", g.Ctx.ServiceName.Lower()))
	if util.FileExists(fileName) {
		return nil
	}

	return util.With("etc").Parse(etcTemplate).SaveTo(map[string]interface{}{
		"serviceName": g.Ctx.ServiceName.Lower(),
	}, fileName, false)
}
