package {{.packageName}}

import (
	"context"

	"{{.projectPath}}/internal/svc"
	"{{.projectPath}}/internal/utils/dberrorhandler"
	"{{.projectPath}}/types/{{.projectName}}"

{{if .useUUID}}    "github.com/suyuan32/simple-admin-common/utils/uuidx"
{{end}}	"github.com/zeromicro/go-zero/core/logx"
)

type Get{{.modelName}}ByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGet{{.modelName}}ByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Get{{.modelName}}ByIdLogic {
	return &Get{{.modelName}}ByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Get{{.modelName}}ByIdLogic) Get{{.modelName}}ById(in *{{.projectName}}.{{if .useUUID}}UU{{end}}IDReq) (*{{.projectName}}.{{.modelName}}Info, error) {
	result, err := l.svcCtx.DB.{{.modelName}}.Get(l.ctx, {{if .useUUID}}uuidx.ParseUUIDString({{end}}in.Id{{if .useUUID}}){{end}})
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &{{.projectName}}.{{.modelName}}Info{
		Id:          result.ID{{if .useUUID}}.String(){{end}},
		CreatedAt:   result.CreatedAt.UnixMilli(),
		UpdatedAt:   result.UpdatedAt.UnixMilli(),
{{.listData}}
	}, nil
}

