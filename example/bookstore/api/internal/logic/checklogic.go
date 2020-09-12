package logic

import (
	"bookstore/api/internal/svc"
	"bookstore/api/internal/types"
	"bookstore/rpc/check/checker"
	"context"

	"github.com/tal-tech/go-zero/core/logx"
)

type CheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) CheckLogic {
	return CheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckLogic) Check(req types.CheckReq) (*types.CheckResp, error) {
	resp, err := l.svcCtx.Checker.Check(l.ctx, &checker.CheckReq{
		Book:  req.Book,
	})
	if err != nil {
		return nil, err
	}

	return &types.CheckResp{
		Found: resp.Found,
		Price: resp.Price,
	}, nil
}