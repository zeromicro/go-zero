package dberrorhandler

import (
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-core/pkg/i18n"
	"github.com/suyuan32/simple-admin-core/pkg/msg/logmsg"
	"github.com/suyuan32/simple-admin-core/pkg/statuserr"

	"{{ .package}}/ent"
)

// DefaultEntError returns errors dealing with default functions.
func DefaultEntError(logger logx.Logger, err error, detail any) error {
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return statuserr.NewInvalidArgumentError(i18n.TargetNotFound)
		case ent.IsConstraintError(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return statuserr.NewInvalidArgumentError(i18n.ConstraintError)
		case ent.IsValidationError(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return statuserr.NewInvalidArgumentError(i18n.ValidationError)
		case ent.IsNotSingular(err):
			logger.Errorw(err.Error(), logx.Field("detail", detail))
			return statuserr.NewInvalidArgumentError(i18n.NotSingularError)
		default:
			logger.Errorw(logmsg.DatabaseError, logx.Field("detail", err.Error()))
			return statuserr.NewInternalError(i18n.DatabaseError)
		}
	}
	return err
}

