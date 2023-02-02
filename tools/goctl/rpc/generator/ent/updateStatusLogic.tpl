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

type Update{{.modelName}}StatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdate{{.modelName}}StatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Update{{.modelName}}StatusLogic {
	return &Update{{.modelName}}StatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Update{{.modelName}}StatusLogic) Update{{.modelName}}Status(in *{{.projectName}}.StatusCodeReq) (*{{.projectName}}.BaseResp, error) {
	err := l.svcCtx.DB.{{.modelName}}.UpdateOneID({{if .useUUID}}uuidx.ParseUUIDString({{end}}in.Id){{if .useUUID}}){{end}}.SetStatus(uint8(in.Status)).Exec(l.ctx)

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

	return &{{.projectName}}.BaseResp{Msg: i18n.UpdateSuccess}, nil
}
