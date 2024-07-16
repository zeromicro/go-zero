package logic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/types"
)

type ApiGenerateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApiGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApiGenerateLogic {
	return &ApiGenerateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApiGenerateLogic) ApiGenerate(req *types.APIGenerateRequest) (resp *types.APIGenerateResponse, err error) {


	return
}
