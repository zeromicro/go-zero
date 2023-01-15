package {{.packageName}}

import (
	"context"

    "{{.projectPath}}/ent"
    "{{.projectPath}}/internal/svc"
    "{{.projectPath}}/{{.projectName}}"

    "github.com/suyuan32/simple-admin-core/pkg/i18n"
    "github.com/suyuan32/simple-admin-core/pkg/msg/logmsg"
    "github.com/suyuan32/simple-admin-core/pkg/statuserr"
{{if .useUUID}}    "github.com/suyuan32/simple-admin-core/pkg/uuidx"
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

func (l *Delete{{.modelName}}Logic) Delete{{.modelName}}(in *{{.projectName}}.{{if .useUUID}}UU{{end}}IDReq) (*{{.projectName}}.BaseResp, error) {
	err := l.svcCtx.DB.{{.modelName}}.DeleteOneID({{if .useUUID}}uuidx.ParseUUIDString({{end}}in.Id){{if .useUUID}}){{end}}.Exec(l.ctx)

	if err != nil {
		switch {
		case ent.IsNotFound(err):
			logx.Errorw(err.Error(), logx.Field("detail", in))
			return nil, statuserr.NewInvalidArgumentError(i18n.TargetNotFound)
		default:
			logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
			return nil, statuserr.NewInternalError(i18n.DatabaseError)
		}
	}

	return &{{.projectName}}.BaseResp{Msg: i18n.DeleteSuccess}, nil
}
