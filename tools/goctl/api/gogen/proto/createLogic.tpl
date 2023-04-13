package {{.modelNameLowerCase}}

import (
	"context"

	"{{.projectPackage}}/internal/svc"
	"{{.projectPackage}}/internal/types"
	"{{.rpcPackage}}"

	"github.com/zeromicro/go-zero/core/logx"
)

type Create{{.modelName}}Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreate{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Create{{.modelName}}Logic {
	return &Create{{.modelName}}Logic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Create{{.modelName}}Logic) Create{{.modelName}}(req *types.{{.modelName}}Info) (resp *types.BaseMsgResp, err error) {
	data, err := l.svcCtx.{{.rpcName}}Rpc.Create{{.modelName}}(l.ctx,
		&{{.rpcPbPackageName}}.{{.modelName}}Info{ {{.setLogic}}
		})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, data.Msg)}, nil
}
