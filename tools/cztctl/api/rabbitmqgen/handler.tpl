package {{.PkgName}}

import (
    "context"
	{{.ImportPackages}}
)

{{if .HasDoc}}{{.Doc}}{{end}}
func {{.HandlerName}}(ctx context.Context, svcCtx *svc.ServiceContext) service.Service {
	return rabbitmq.MustNewListener(ctx, svcCtx.Config.{{.RabbitmqConfName}}, {{.LogicName}}.New{{.LogicType}}(ctx, svcCtx))
}

