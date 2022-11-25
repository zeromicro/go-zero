package user

import (
	"context"

	"github.com/zeromicro/go-zero/tools/goctl/example/api/test2/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test2/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ModifyUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ModifyUserInfoLogic {
	return &ModifyUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyUserInfoLogic) ModifyUserInfo(req *types.ModifyUserInfoRequest) (resp *types.ModifyUserInfoResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
