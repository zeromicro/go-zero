package user

import (
	"context"

	"github.com/zeromicro/go-zero/tools/goctl/example/api/test1/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test1/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAllLogic {
	return &UserAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserAllLogic) UserAll(req *types.UserListRequest) (resp []types.Base, err error) {
	// todo: add your logic here and delete this line

	return
}
