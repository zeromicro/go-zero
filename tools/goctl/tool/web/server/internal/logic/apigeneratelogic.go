package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/tools/goctl/util"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/tool/web/server/internal/types"
)

var errMissingServiceName = errors.New("missing service name")

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
	if err := l.validateAPIGenerateRequest(req); err != nil {
		return nil, err
	}
	return
}

func (l *ApiGenerateLogic) validateAPIGenerateRequest(req *types.APIGenerateRequest) error {
	if util.IsEmptyStringOrWhiteSpace(req.Name) {
		return errMissingServiceName
	}

}
