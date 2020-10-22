package gen

import (
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
	"github.com/tal-tech/go-zero/tools/goctl/util/stringx"
)

const rpcTemplateText = `syntax = "proto3";

package {{.package}};

message Request {
  string ping = 1;
}

message Response {
  string pong = 1;
}

service {{.serviceName}} {
  rpc Ping(Request) returns(Response);
}
`

type rpcTemplate struct {
	out string
	console.Console
}

func NewRpcTemplate(out string, idea bool) *rpcTemplate {
	return &rpcTemplate{
		out:     out,
		Console: console.NewConsole(idea),
	}
}

func (r *rpcTemplate) MustGenerate(showState bool) {
	r.Info(aurora.Blue("-> goctl rpc reference documents: ").String() + "「https://github.com/tal-tech/zero-doc/blob/main/doc/goctl-rpc.md」")
	r.Info("-> generating template...")
	protoFilename := filepath.Base(r.out)
	serviceName := stringx.From(strings.TrimSuffix(protoFilename, filepath.Ext(protoFilename)))
	text, err := util.LoadTemplate(category, rpcTemplateFile, rpcTemplateText)
	r.Must(err)

	err = util.With("t").Parse(text).SaveTo(map[string]string{
		"package":     serviceName.UnTitle(),
		"serviceName": serviceName.Title(),
	}, r.out, false)
	r.Must(err)

	if showState {
		r.Success("Done.")
	}
}
