package {{.packageName}}

import (
	"context"

	"{{.projectPath}}/ent/{{.modelNameLowerCase}}"
	"{{.projectPath}}/ent/predicate"
	"{{.projectPath}}/internal/svc"
	"{{.projectPath}}/internal/types"
	"{{.projectPath}}/internal/utils/dberrorhandler"

    "github.com/suyuan32/simple-admin-common/i18n"
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

func (l *Get{{.modelName}}ListLogic) Get{{.modelName}}List(req *types.{{.modelName}}ListReq) (*types.{{.modelName}}ListResp, error) {
{{.predicateData}}

	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	resp := &types.{{.modelName}}ListResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.ctx, i18n.Success)
	resp.Data.Total = data.PageDetails.Total

	for _, v := range data.List {
		resp.Data.Data = append(resp.Data.Data,
		types.{{.modelName}}Info{
            Base{{if .useUUID}}UU{{end}}IDInfo:    types.Base{{if .useUUID}}UU{{end}}IDInfo{
                Id: v.ID{{if .useUUID}}.String(){{end}},
                CreatedAt: v.CreatedAt.UnixMilli(),
                UpdatedAt: v.UpdatedAt.UnixMilli(),
            },
{{.listData}}
		})
	}

	return resp, nil
}
