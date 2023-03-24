package {{.packageName}}

import (
	"context"
{{if .hasTime}}     "time"{{end}}

	"{{.projectPath}}/internal/svc"
	"{{.projectPath}}/internal/types"
	"{{.projectPath}}/internal/utils/dberrorhandler"

    "github.com/suyuan32/simple-admin-common/i18n"
{{if .hasUUID}}    "github.com/suyuan32/simple-admin-common/utils/uuidx"
{{end}}
	"github.com/zeromicro/go-zero/core/logx"
)

type Create{{.modelName}}Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreate{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *Create{{.modelName}}Logic {
	return &Create{{.modelName}}Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Create{{.modelName}}Logic) Create{{.modelName}}(req *types.{{.modelName}}Info) (*types.BaseMsgResp, error) {
    _, err := l.svcCtx.DB.{{.modelName}}.Create().
{{.setLogic}}

    if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, req)
	}

    return &types.BaseMsgResp{Msg: l.svcCtx.Trans.Trans(l.ctx, i18n.CreateSuccess)}, nil
}
