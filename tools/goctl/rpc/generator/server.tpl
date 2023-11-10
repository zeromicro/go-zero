{{.head}}

package server

import (
	{{if .notStream}}"context"{{end}}

	{{.imports}}
)

type {{.server}} struct {
	svcCtx *svc.ServiceContext
	{{.unimplementedServer}}
}

func New{{.server}}(svcCtx *svc.ServiceContext) *{{.server}} {
	return &{{.server}}{
		svcCtx: svcCtx,
	}
}

{{.funcs}}
