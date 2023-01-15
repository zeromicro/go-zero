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
{{if or .hasUUID .useUUID}}	"github.com/suyuan32/simple-admin-core/pkg/uuidx"{{end}}
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrUpdate{{.modelName}}Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrUpdate{{.modelName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrUpdate{{.modelName}}Logic {
	return &CreateOrUpdate{{.modelName}}Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateOrUpdate{{.modelName}}Logic) CreateOrUpdate{{.modelName}}(in *{{.projectName}}.{{.modelName}}Info) (*{{.projectName}}.BaseResp, error) {
    if in.Id == {{if .useUUID}}""{{else}}0{{end}} {
        err := l.svcCtx.DB.{{.modelName}}.Create().
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

        return &{{.projectName}}.BaseResp{Msg: i18n.CreateSuccess}, nil
    } else {
        err := l.svcCtx.DB.{{.modelName}}.UpdateOneID({{if .useUUID}}uuidx.ParseUUIDString({{end}}in.Id){{if .useUUID}}){{end}}.
{{.setLogic}}

        if err != nil {
            switch {
            case ent.IsNotFound(err):
                logx.Errorw(err.Error(), logx.Field("detail", in))
                return nil, statuserr.NewInvalidArgumentError(i18n.TargetNotFound)
            case ent.IsConstraintError(err):
                logx.Errorw(err.Error(), logx.Field("detail", in))
                return nil, statuserr.NewInvalidArgumentError(i18n.UpdateFailed)
            default:
                logx.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
                return nil, statuserr.NewInternalError(i18n.DatabaseError)
            }
        }

        return &{{.projectName}}.BaseResp{Msg: i18n.UpdateSuccess}, nil
    }
}
