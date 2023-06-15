package {{.packageName}}

import (
	"context"
{{if .hasTime}}     "time"{{end}}

	"{{.projectPath}}{{.importPrefix}}/internal/svc"
	"{{.projectPath}}{{.importPrefix}}/internal/utils/dberrorhandler"
	"{{.projectPath}}{{.importPrefix}}/types/{{.projectName}}"

{{if .useI18n}}    "github.com/suyuan32/simple-admin-common/i18n"
{{else}}    "github.com/suyuan32/simple-admin-common/msg/errormsg"
{{end}}{{if or .hasUUID .useUUID}}	"github.com/suyuan32/simple-admin-common/utils/uuidx"{{end}}
	"github.com/suyuan32/simple-admin-common/utils/pointy"
	"github.com/zeromicro/go-zero/core/logx"
)

type Update{{.modelName}}Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdate{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Update{{.modelName}}Logic {
	return &Update{{.modelName}}Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Update{{.modelName}}Logic) Update{{.modelName}}(in *{{.projectName}}.{{.modelName}}Info) (*{{.projectName}}.BaseResp, error) {
    err := l.svcCtx.DB.{{.modelName}}.UpdateOneID({{if .useUUID}}uuidx.ParseUUIDString({{end}}*in.Id){{if .useUUID}}){{end}}.
{{.setLogic}}

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &{{.projectName}}.BaseResp{Msg: {{if .useI18n}}i18n.CreateSuccess{{else}}errormsg.UpdateSuccess{{end}} }, nil
}
