{{.head}}

package {{.filePackage}}

import (
	{{.pbPackage}}
	{{if ne .pbPackage .protoGoPackage}}{{.protoGoPackage}}{{end}}

	"github.com/zeromicro/go-zero/zrpc"
)

type {{.serviceName}} = {{if .isCallPkgSameToGrpcPkg}}{{else}}{{.package}}.{{end}}{{.rpcServiceName}}Client

func New{{.serviceName}}(cli zrpc.Client) {{.serviceName}} {
	return {{if .isCallPkgSameToGrpcPkg}}{{else}}{{.package}}.{{end}}New{{.rpcServiceName}}Client(cli.Conn())
}