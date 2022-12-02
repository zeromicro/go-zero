package {{.modelNameLowerCase}}

import (
	"context"
	"net/http"

	"{{.projectPackage}}/internal/svc"
	"{{.projectPackage}}/internal/types"
	"{{.rpcPackage}}"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdate{{.modelName}}Logic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewCreateOrUpdate{{.modelName}}Logic(r *http.Request, svcCtx *svc.ServiceContext) *CreateOrUpdate{{.modelName}}Logic {
	return &CreateOrUpdate{{.modelName}}Logic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		lang:   r.Header.Get("Accept-Language"),
	}
}

func (l *CreateOrUpdate{{.modelName}}Logic) CreateOrUpdate{{.modelName}}(req *types.CreateOrUpdate{{.modelName}}Req) (resp *types.BaseMsgResp, err error) {
	data, err := l.svcCtx.{{.rpcName}}Rpc.CreateOrUpdate{{.modelName}}(l.ctx,
		&{{.rpcNameLowerCase}}.{{.modelName}}Info{
			Id:          req.Id,{{.setLogic}}
		})
	if err != nil {
		return nil, err
	}
	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.lang, data.Msg)}, nil
}
