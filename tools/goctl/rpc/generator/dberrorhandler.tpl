package dberrorhandler

import (
    "github.com/zeromicro/go-zero/core/errorx"
	"github.com/zeromicro/go-zero/core/logx"

{{if .useI18n}}    "github.com/suyuan32/simple-admin-common/i18n"
{{else}}    "github.com/suyuan32/simple-admin-common/msg/errormsg"
{{end}}	"github.com/suyuan32/simple-admin-common/msg/logmsg"

	"{{ .package}}/ent"
)

// DefaultEntError returns errors dealing with default functions.
func DefaultEntError(logger logx.Logger, err error, detail any) error {
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return errorx.NewInvalidArgumentError({{if .useI18n}}i18n.TargetNotFound{{else}}errormsg.TargetNotFound{{end}})
		case ent.IsConstraintError(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return errorx.NewInvalidArgumentError({{if .useI18n}}i18n.ConstraintError{{else}}errormsg.ConstraintError{{end}})
		case ent.IsValidationError(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return errorx.NewInvalidArgumentError({{if .useI18n}}i18n.ValidationError{{else}}errormsg.ValidationError{{end}})
		case ent.IsNotSingular(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return errorx.NewInvalidArgumentError({{if .useI18n}}i18n.NotSingularError{{else}}errormsg.NotSingularError{{end}})
		default:
			logger.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
			return errorx.NewInternalError({{if .useI18n}}i18n.DatabaseError{{else}}errormsg.DatabaseError{{end}})
		}
	}
	return err
}
