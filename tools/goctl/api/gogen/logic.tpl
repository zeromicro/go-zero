package {{.pkgName}}

import (
	{{.imports}}
)

type {{.logic}} struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

{{if .hasDoc}}{{.doc}}{{end}}
func New{{.logic}}(ctx context.Context, svcCtx *svc.ServiceContext) *{{.logic}} {
	return &{{.logic}}{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *{{.logic}}) {{.function}}({{.request}}) {{.responseType}} {
	// todo: add your logic here and delete this line
	{{ range $k,$v:=.openFiles }}
	{{$v}},{{$v}}Header, err := r.FormFile("{{$v}}")
        if err != nil {
		return 
	}
	defer {{$v}}.Close()
	_ = {{$v}}Header
	{{end}}

	{{.returnString}}
}
