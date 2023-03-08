package {{.modelNameLowerCase}}

import (
	"context"

	"{{.projectPackage}}/internal/svc"
	"{{.projectPackage}}/internal/types"
	"{{.rpcPackage}}"

	"github.com/zeromicro/go-zero/core/logx"
)

type Update{{.modelName}}Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdate{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Update{{.modelName}}Logic {
	return &Update{{.modelName}}Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Update{{.modelName}}Logic) Update{{.modelName}}(req *types.{{.modelName}}Info) (resp *types.BaseMsgResp, err error) {
	data, err := l.svcCtx.{{.rpcName}}Rpc.Update{{.modelName}}(l.ctx,
		&{{.rpcPbPackageName}}.{{.modelName}}Info{
			Id:          req.Id,{{.setLogic}}
		})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, data.Msg)}, nil
}
