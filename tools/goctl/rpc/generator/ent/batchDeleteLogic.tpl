package {{.packageName}}

import (
	"context"

    "{{.projectPath}}/ent"
    "{{.projectPath}}/ent/{{.modelNameLowerCase}}"
    "{{.projectPath}}/internal/svc"
    "{{.projectPath}}/{{.serviceName}}"

    "github.com/suyuan32/simple-admin-core/pkg/i18n"
    "github.com/suyuan32/simple-admin-core/pkg/msg/logmsg"
    "github.com/suyuan32/simple-admin-core/pkg/statuserr"
    "github.com/zeromicro/go-zero/core/logx"
)

type BatchDelete{{.modelName}}Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchDelete{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDelete{{.modelName}}Logic {
	return &BatchDelete{{.modelName}}Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchDelete{{.modelName}}Logic) BatchDelete{{.modelName}}(in *{{.serviceName}}.IDsReq) (*{{.serviceName}}.BaseResp, error) {
	_, err := l.svcCtx.DB.{{.modelName}}.Delete().Where({{.modelNameLowerCase}}.IDIn(in.Ids...)).Exec(l.ctx)

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

	return &{{.serviceName}}.BaseResp{Msg: i18n.DeleteSuccess}, nil
}
