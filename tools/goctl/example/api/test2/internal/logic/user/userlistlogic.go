package user

import (
	"context"

	"github.com/zeromicro/go-zero/tools/goctl/example/api/test2/internal/svc"
	"github.com/zeromicro/go-zero/tools/goctl/example/api/test2/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListLogic {
	return &UserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserListLogic) UserList(req *types.UserListRequest) (resp []types.Base, err error) {
	// todo: add your logic here and delete this line

	return
}
