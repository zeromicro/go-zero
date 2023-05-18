package {{.packageName}}

import (
	"context"

    "{{.projectPath}}{{.importPrefix}}/ent/{{.modelNameLowerCase}}"
    "{{.projectPath}}{{.importPrefix}}/internal/svc"
    "{{.projectPath}}{{.importPrefix}}/internal/types"
    "{{.projectPath}}{{.importPrefix}}/internal/utils/dberrorhandler"

{{if .useI18n}}    "github.com/suyuan32/simple-admin-common/i18n"
{{else}}    "github.com/suyuan32/simple-admin-common/msg/errormsg"
{{end}}{{if .useUUID}}    "github.com/suyuan32/simple-admin-common/utils/uuidx"
{{end}}    "github.com/zeromicro/go-zero/core/logx"
)

type Delete{{.modelName}}Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelete{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Delete{{.modelName}}Logic {
	return &Delete{{.modelName}}Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Delete{{.modelName}}Logic) Delete{{.modelName}}(req *types.{{if .useUUID}}UU{{end}}IDsReq) (*types.BaseMsgResp, error) {
	_, err := l.svcCtx.DB.{{.modelName}}.Delete().Where({{.modelNameLowerCase}}.IDIn({{if .useUUID}}uuidx.ParseUUIDSlice({{end}}req.Ids{{if .useUUID}}){{end}}...)).Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

    return &types.BaseMsgResp{Msg: {{if .useI18n}}l.svcCtx.Trans.Trans(l.ctx, i18n.DeleteSuccess){{else}}errormsg.DeleteSuccess{{end}}}, nil
}
