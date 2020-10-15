package gen

import (
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/templatex"
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
	r.Info("查看rpc生成请移步至「https://github.com/tal-tech/go-zero/blob/master/doc/goctl-rpc.md」")
	r.Info("generating template...")
	protoFilename := filepath.Base(r.out)
	serviceName := stringx.From(strings.TrimSuffix(protoFilename, filepath.Ext(protoFilename)))
	err := templatex.With("t").Parse(rpcTemplateText).SaveTo(map[string]string{
		"package":     serviceName.UnTitle(),
		"serviceName": serviceName.Title(),
	}, r.out, false)
	r.Must(err)
	if showState {
		r.Success("Done.")
	}
}
