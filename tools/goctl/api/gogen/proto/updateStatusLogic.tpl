package {{.modelNameLowerCase}}

import (
	"context"
	"net/http"

    "{{.projectPackage}}/internal/svc"
	"{{.projectPackage}}/internal/types"
	"{{.rpcPackage}}"

	"github.com/zeromicro/go-zero/core/logx"
)

type Update{{.modelName}}StatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	lang   string
}

func NewUpdate{{.modelName}}StatusLogic(r *http.Request, svcCtx *svc.ServiceContext) *Update{{.modelName}}StatusLogic {
	return &Update{{.modelName}}StatusLogic{
		Logger: logx.WithContext(r.Context()),
		ctx:    r.Context(),
		svcCtx: svcCtx,
		lang:   r.Header.Get("Accept-Language"),
	}
}

func (l *Update{{.modelName}}StatusLogic) Update{{.modelName}}Status(req *types.StatusCode{{if .useUUID}}UUID{{end}}Req) (resp *types.BaseMsgResp, err error) {
	result, err := l.svcCtx.{{.rpcName}}Rpc.Update{{.modelName}}Status(l.ctx, &{{.rpcPbPackageName}}.StatusCode{{if .useUUID}}UUID{{end}}Req{
		Id: req.Id,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}

	return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.lang, result.Msg)}, nil
}
