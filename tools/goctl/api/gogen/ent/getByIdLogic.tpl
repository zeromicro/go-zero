package {{.packageName}}

import (
	"context"

	"{{.projectPath}}/internal/svc"
	"{{.projectPath}}/internal/types"
	"{{.projectPath}}/internal/utils/dberrorhandler"

    "github.com/suyuan32/simple-admin-common/i18n"
{{if .useUUID}}    "github.com/suyuan32/simple-admin-common/utils/uuidx"
{{end}}	"github.com/zeromicro/go-zero/core/logx"
)

type Get{{.modelName}}ByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGet{{.modelName}}ByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Get{{.modelName}}ByIdLogic {
	return &Get{{.modelName}}ByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Get{{.modelName}}ByIdLogic) Get{{.modelName}}ById(req *types.{{if .useUUID}}UU{{end}}IDReq) (*types.{{.modelName}}InfoResp, error) {
	data, err := l.svcCtx.DB.{{.modelName}}.Get(l.ctx, {{if .useUUID}}uuidx.ParseUUIDString({{end}}req.Id{{if .useUUID}}){{end}})
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

	return &types.{{.modelName}}InfoResp{
	    BaseDataInfo: types.BaseDataInfo{
            Code: 0,
            Msg:  l.svcCtx.Trans.Trans(l.ctx, i18n.Success),
        },
        Data: types.{{.modelName}}Info{
            Base{{if .useUUID}}UU{{end}}IDInfo:    types.Base{{if .useUUID}}UU{{end}}IDInfo{
                Id: data.ID{{if .useUUID}}.String(){{end}},
                CreatedAt: data.CreatedAt.UnixMilli(),
                UpdatedAt: data.UpdatedAt.UnixMilli(),
            },
{{.listData}}
        },
	}, nil
}

