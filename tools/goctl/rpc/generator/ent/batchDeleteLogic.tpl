package logic

import (
	"context"

    "github.com/suyuan32/simple-admin-core/pkg/i18n"
    "github.com/suyuan32/simple-admin-core/pkg/msg/logmsg"
    "github.com/suyuan32/simple-admin-core/pkg/statuserr"
    "{{.projectPath}}/pkg/ent"
    "{{.projectPath}}/rpc/internal/svc"
    "{{.projectPath}}/rpc/types/{{.serviceName}}"

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

func (l *BatchDeleteApiLogic) BatchDelete{{.modelName}}(in *{{.serviceName}}.IDsReq) (*{{.serviceName}}.BaseResp, error) {
	err := l.svcCtx.DB.{{.modelName}}.Delete().Where(token.IDIn(in.Ids)).Exec(l.ctx)

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
