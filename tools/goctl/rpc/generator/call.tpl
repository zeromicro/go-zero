{{.head}}

package {{.filePackage}}

import (
	"context"

	{{.imports}}

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (

	{{.serviceName}} interface {
		{{.interface}}
	}

	default{{.serviceName}} struct {
		cli zrpc.Client
	}
)

func New{{.serviceName}}(cli zrpc.Client) {{.serviceName}} {
	return &default{{.serviceName}}{
		cli: cli,
	}
}

{{.functions}}
