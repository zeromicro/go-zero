package user

import (
	"context"

	"github.com/zeromicro/go-zero/tools/goctl/example/api/test1/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test1/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FansLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFansLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FansLogic {
	return &FansLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FansLogic) Fans(req *types.FansRequest) (resp *types.FansResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
