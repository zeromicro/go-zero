package {{.modelNameLowerCase}}

import (
	"context"

    "{{.projectPackage}}/internal/svc"
	"{{.projectPackage}}/internal/types"
	"{{.rpcPackage}}"

	"github.com/suyuan32/simple-admin-common/i18n"
	"github.com/zeromicro/go-zero/core/logx"
)

type Get{{.modelName}}ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGet{{.modelName}}ListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Get{{.modelName}}ListLogic {
	return &Get{{.modelName}}ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Get{{.modelName}}ListLogic) Get{{.modelName}}List(req *types.{{.modelName}}ListReq) (resp *types.{{.modelName}}ListResp, err error) {
	data, err := l.svcCtx.{{.rpcName}}Rpc.Get{{.modelName}}List(l.ctx,
		&{{.rpcPbPackageName}}.{{.modelName}}ListReq{
			Page:        req.Page,
			PageSize:    req.PageSize,{{.searchKeys}}
		})
	if err != nil {
		return nil, err
	}
	resp = &types.{{.modelName}}ListResp{}
	resp.Msg = l.svcCtx.Trans.Trans(l.ctx, i18n.Success)
	resp.Data.Total = data.GetTotal()

	for _, v := range data.Data {
		resp.Data.Data = append(resp.Data.Data,
			types.{{.modelName}}Info{
				Base{{if .useUUID}}UU{{end}}IDInfo: types.Base{{if .useUUID}}UU{{end}}IDInfo{
					Id:        v.Id,
					CreatedAt: v.CreatedAt,
					UpdatedAt: v.UpdatedAt,
				},{{.setLogic}}
			})
	}
	return resp, nil
}
