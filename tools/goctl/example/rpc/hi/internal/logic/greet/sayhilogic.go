package greetlogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/tools/goctl/example/rpc/hi/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/example/rpc/hi/pb/hi"
)

type SayHiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSayHiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SayHiLogic {
	return &SayHiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SayHiLogic) SayHi(in *hi.HiReq) (*hi.HiResp, error) {
	// todo: add your logic here and delete this line

	return &hi.HiResp{}, nil
}
