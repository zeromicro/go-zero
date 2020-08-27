package gogen

import (
	"fmt"
	"path/filepath"

	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/templatex"
)

var etcTemplate = `{
  "Name": "{{.serviceName}}.rpc",
  "Log": {
    "Mode": "console"
  },
  "ListenOn": "127.0.0.1:8080",
  "Etcd": {
    "Hosts": ["127.0.0.1:6379"],
    "Key": "{{.serviceName}}.rpc"
  }
}
`

func (g *defaultRpcGenerator) genEtc() error {
	etdDir := g.dirM[dirEtc]
	fileName := filepath.Join(etdDir, fmt.Sprintf("%v.json", g.Ctx.ServiceName.Lower()))
	if util.FileExists(fileName) {
		return nil
	}
	return templatex.With("etc").
		Parse(etcTemplate).
		SaveTo(map[string]interface{}{
			"serviceName": g.Ctx.ServiceName.Lower(),
		}, fileName, false)
}
