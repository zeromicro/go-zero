package {{.packageName}}

import (
	"context"

    "{{.projectPath}}/ent/{{.modelNameLowerCase}}"
    "{{.projectPath}}/internal/svc"
    "{{.projectPath}}/internal/utils/dberrorhandler"
    "{{.projectPath}}/types/{{.projectName}}"

    "github.com/suyuan32/simple-admin-common/i18n"
{{if .useUUID}}    "github.com/suyuan32/simple-admin-common/utils/uuidx"
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

func (l *Delete{{.modelName}}Logic) Delete{{.modelName}}(in *{{.projectName}}.{{if .useUUID}}UU{{end}}IDsReq) (*{{.projectName}}.BaseResp, error) {
	_, err := l.svcCtx.DB.{{.modelName}}.Delete().Where({{.modelNameLowerCase}}.IDIn({{if .useUUID}}uuidx.ParseUUIDSlice({{end}}in.Ids{{if .useUUID}}){{end}}...)).Exec(l.ctx)

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

    return &{{.projectName}}.BaseResp{Msg: i18n.DeleteSuccess}, nil
}
