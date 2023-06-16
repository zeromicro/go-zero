package {{.packageName}}

import (
	"context"

	"{{.projectPath}}{{.importPrefix}}/ent/{{.modelNameLowerCase}}"
	"{{.projectPath}}{{.importPrefix}}/ent/predicate"
	"{{.projectPath}}{{.importPrefix}}/internal/svc"
	"{{.projectPath}}{{.importPrefix}}/internal/types"
	"{{.projectPath}}{{.importPrefix}}/internal/utils/dberrorhandler"

{{if .useI18n}}    "github.com/suyuan32/simple-admin-common/i18n"
{{else}}    "github.com/suyuan32/simple-admin-common/msg/errormsg"
{{end}}
	"github.com/suyuan32/simple-admin-common/utils/pointy"
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
	resp.Msg = {{if .useI18n}}l.svcCtx.Trans.Trans(l.ctx, i18n.Success){{else}}errormsg.Success{{end}}
	resp.Data.Total = data.PageDetails.Total

	for _, v := range data.List {
		resp.Data.Data = append(resp.Data.Data,
		types.{{.modelName}}Info{
            Base{{if .useUUID}}UU{{end}}IDInfo:    types.Base{{if .useUUID}}UU{{end}}IDInfo{
				Id:          {{if .useUUID}}pointy.GetPointer(v.ID.String()){{else}}&v.ID{{end}},
				CreatedAt:    pointy.GetPointer(v.CreatedAt.Unix()),
				UpdatedAt:    pointy.GetPointer(v.UpdatedAt.Unix()),
            },
{{.listData}}
		})
	}

	return resp, nil
}
