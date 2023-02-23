package {{.packageName}}

import (
	"context"
{{if .hasTime}}     "time"{{end}}

	"{{.projectPath}}/ent"
	"{{.projectPath}}/internal/svc"
    "{{.projectPath}}/{{.projectName}}"

    "github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/suyuan32/simple-admin-core/pkg/msg/logmsg"
	"github.com/suyuan32/simple-admin-core/pkg/statuserr"
{{if .hasUUID}}    "github.com/suyuan32/simple-admin-core/pkg/uuidx"
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

func (l *Create{{.modelName}}Logic) Create{{.modelName}}(in *{{.projectName}}.{{.modelName}}Info) (*{{.projectName}}.Base{{if .useUUID}}UU{{end}}IDResp, error) {
    result, err := l.svcCtx.DB.{{.modelName}}.Create().
{{.setLogic}}

    if err != nil {
        switch {
        case ent.IsConstraintError(err):
            logx.Errorw(err.Error(), logx.Field("detail", in))
            return nil, statuserr.NewInvalidArgumentError(i18n.CreateFailed)
        default:
            logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
            return nil, statuserr.NewInternalError(i18n.DatabaseError)
        }
    }

    return &{{.projectName}}.Base{{if .useUUID}}UU{{end}}IDResp{Id: result.ID{{if .useUUID}}.String(){{end}}, Msg: i18n.CreateSuccess}, nil
}
