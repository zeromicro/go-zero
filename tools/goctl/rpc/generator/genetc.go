package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/rpc/parser"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const etcTemplate = `Name: {{.serviceName}}.rpc
ListenOn: 127.0.0.1:8080
Etcd:
  Hosts:
  - 127.0.0.1:2379
  Key: {{.serviceName}}.rpc
`

func (g *defaultGenerator) GenEtc(ctx DirContext, _ parser.Proto, namingStyle NamingStyle) error {
	dir := ctx.GetEtc()
	serviceNameLower := formatFilename(ctx.GetMain().Base, namingStyle)
	fileName := filepath.Join(dir.Filename, fmt.Sprintf("%v.yaml", serviceNameLower))

	text, err := util.LoadTemplate(category, etcTemplateFileFile, etcTemplate)
	if err != nil {
		return err
	}

	return util.With("etc").Parse(text).SaveTo(map[string]interface{}{
		"serviceName": strings.ToLower(stringx.From(ctx.GetMain().Base).ToCamel()),
	}, fileName, false)
}
