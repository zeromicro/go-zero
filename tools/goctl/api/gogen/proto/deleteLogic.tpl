package {{.modelNameLowerCase}}

import (
	"context"

    "{{.projectPackage}}/internal/svc"
	"{{.projectPackage}}/internal/types"
	"{{.rpcPackage}}"

	"github.com/zeromicro/go-zero/core/logx"
)

type Delete{{.modelName}}Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelete{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Delete{{.modelName}}Logic {
	return &Delete{{.modelName}}Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Delete{{.modelName}}Logic) Delete{{.modelName}}(req *types.{{if .useUUID}}UU{{end}}IDsReq) (resp *types.BaseMsgResp, err error) {
	result, err := l.svcCtx.{{.rpcName}}Rpc.Delete{{.modelName}}(l.ctx, &{{.rpcPbPackageName}}.{{if .useUUID}}UU{{end}}IDsReq{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, result.Msg)}, nil
}
