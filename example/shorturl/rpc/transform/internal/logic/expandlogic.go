package logic

import (
	"context"

	"shorturl/rpc/transform/internal/svc"
	transform "shorturl/rpc/transform/pb"

	"github.com/tal-tech/go-zero/core/logx"
)

type ExpandLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewExpandLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExpandLogic {
	return &ExpandLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ExpandLogic) Expand(in *transform.ExpandReq) (*transform.ExpandResp, error) {
	res, err := l.svcCtx.Model.FindOne(in.Shorten)
	if err != nil {
		return nil, err
	}

	return &transform.ExpandResp{
		Url: res.Url,
	}, nil
}
