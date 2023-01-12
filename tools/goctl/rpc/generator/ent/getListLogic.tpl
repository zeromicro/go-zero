package {{.packageName}}

import (
	"context"

	"{{.projectPath}}/ent/{{.modelNameLowerCase}}"
	"{{.projectPath}}/ent/predicate"
	"{{.projectPath}}/internal/svc"
    "{{.projectPath}}/{{.serviceName}}"

    "github.com/suyuan32/simple-admin-core/pkg/i18n"
    "github.com/suyuan32/simple-admin-core/pkg/statuserr"{{if .useUUID}}
    "github.com/suyuan32/simple-admin-core/pkg/uuidx"
{{end}}    "github.com/zeromicro/go-zero/core/logx"
)

type Get{{.modelName}}ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGet{{.modelName}}ListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Get{{.modelName}}ListLogic {
	return &Get{{.modelName}}ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *Get{{.modelName}}ListLogic) Get{{.modelName}}List(in *{{.serviceName}}.{{.modelName}}PageReq) (*{{.serviceName}}.{{.modelName}}ListResp, error) {
{{.predicateData}}

	if err != nil {
		logx.Error(err.Error())
		return nil, statuserr.NewInternalError(i18n.DatabaseError)
	}

	resp := &{{.serviceName}}.{{.modelName}}ListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &{{.serviceName}}.{{.modelName}}Info{
			Id:          v.ID{{if .useUUID}}.String(){{end}},
			CreatedAt:   v.CreatedAt.UnixMilli(),
{{.listData}}
		})
	}

	return resp, nil
}
