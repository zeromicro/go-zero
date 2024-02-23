package cache

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ StudentModel = (*customStudentModel)(nil)

type (
	// StudentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customStudentModel.
	StudentModel interface {
		studentModel
	}

	customStudentModel struct {
		*defaultStudentModel
	}
)

// NewStudentModel returns a model for the database table.
func NewStudentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) StudentModel {
	return &customStudentModel{
		defaultStudentModel: newStudentModel(conn, c, opts...),
	}
}
