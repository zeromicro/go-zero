package gen

import (
	"github.com/tal-tech/go-zero/tools/goctl/util"
	"github.com/tal-tech/go-zero/tools/goctl/util/console"
)

const rpcTemplateText = `syntax = "proto3";

package remote;

message Request {
  string username = 1;
  string password = 2;
}

message Response {
  string name = 1;
  string gender = 2;
}

service User {
  rpc Login(Request) returns(Response);
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

func (r *rpcTemplate) MustGenerate() {
	err := util.With("t").Parse(rpcTemplateText).SaveTo(nil, r.out, false)
	r.Must(err)
	r.Success("Done.")
}
