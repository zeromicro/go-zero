package {{.modelNameLowerCase}}

import (
	"context"
	"net/http"

    "{{.projectPackage}}/internal/svc"
	"{{.projectPackage}}/internal/types"
	"{{.rpcPackage}}"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDelete{{.modelName}}Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewBatchDelete{{.modelName}}Logic(r *http.Request, svcCtx *svc.ServiceContext) *BatchDelete{{.modelName}}Logic {
	return &BatchDelete{{.modelName}}Logic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		lang:   r.Header.Get("Accept-Language"),
	}
}

func (l *BatchDelete{{.modelName}}Logic) BatchDelete{{.modelName}}(req *types.IDsReq) (resp *types.BaseMsgResp, err error) {
	result, err := l.svcCtx.{{.rpcName}}Rpc.BatchDelete{{.modelName}}(l.ctx, &{{.rpcNameLowerCase}}.IDReq{
		Ids: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.lang, result.Msg)}, nil
}
