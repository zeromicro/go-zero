package {{.modelNameLowerCase}}

import (
	"context"
	"net/http"

    "{{.projectPackage}}/internal/svc"
	"{{.projectPackage}}/internal/types"
	"{{.rpcPackage}}"

	"github.com/zeromicro/go-zero/core/logx"
)

type Delete{{.modelName}}Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewDelete{{.modelName}}Logic(r *http.Request, svcCtx *svc.ServiceContext) *Delete{{.modelName}}Logic {
	return &Delete{{.modelName}}Logic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		lang:   r.Header.Get("Accept-Language"),
	}
}

func (l *Delete{{.modelName}}Logic) Delete{{.modelName}}(req *types.IDReq) (resp *types.BaseMsgResp, err error) {
	result, err := l.svcCtx.{{.rpcName}}Rpc.Delete{{.modelName}}(l.ctx, &{{.rpcNameLowerCase}}.IDReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.lang, result.Msg)}, nil
}
