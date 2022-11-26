package logic

import (
	"context"

	"github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/suyuan32/simple-admin-core/pkg/statuserr"
	"{{.projectPath}}/pkg/ent/{{.modelNameLowerCase}}"
	"{{.projectPath}}/pkg/ent/predicate"
	"{{.projectPath}}/rpc/internal/svc"
	"{{.projectPath}}/rpc/types/{{.serviceName}}"

	"github.com/zeromicro/go-zero/core/logx"
)

type Get{{.modelName}}ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGet{{.modelName}}ListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Get{{.modelName}}ListLogic {
	return &Get{{.modelName}}ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Get{{.modelName}}ListLogic) Get{{.modelName}}List(in *{{.serviceName}}.{{.modelName}}PageReq) (*{{.serviceName}}.{{.modelName}}ListResp, error) {
{{.predicateData}}

	if err != nil {
		logx.Error(err.Error())
		return nil, statuserr.NewInternalError(i18n.DatabaseError)
	}

	resp := &{{.serviceName}}.{{.modelName}}ListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &{{.serviceName}}.{{.modelName}}Info{
			Id:          v.ID,
			CreatedAt:   v.CreatedAt.UnixMilli(),
{{.listData}}
		})
	}

	return resp, nil
}
