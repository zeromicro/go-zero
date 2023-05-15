package {{.packageName}}

import (
	"context"
{{if .hasTime}}     "time"{{end}}

	"{{.projectPath}}/internal/svc"
	"{{.projectPath}}/internal/types"
	"{{.projectPath}}/internal/utils/dberrorhandler"

{{if .useI18n}}    "github.com/suyuan32/simple-admin-common/i18n"
{{else}}    "github.com/suyuan32/simple-admin-common/msg/errormsg"
{{end}}{{if or .hasUUID .useUUID}}	"github.com/suyuan32/simple-admin-common/utils/uuidx"{{end}}
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

func (l *Update{{.modelName}}Logic) Update{{.modelName}}(req *types.{{.modelName}}Info) (*types.BaseMsgResp, error) {
    err := l.svcCtx.DB.{{.modelName}}.UpdateOneID({{if .useUUID}}uuidx.ParseUUIDString({{end}}req.Id){{if .useUUID}}){{end}}.
{{.setLogic}}

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

    return &types.BaseMsgResp{Msg: {{if .useI18n}}l.svcCtx.Trans.Trans(l.ctx, i18n.UpdateSuccess){{else}}errormsg.UpdateSuccess{{end}}}, nil
}
