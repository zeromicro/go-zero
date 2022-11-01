package gormsql

import (
	"gorm.io/gorm/logger"

	"github.com/zeromicro/go-zero/core/logx"
)

type gormWriter struct {
	logger.Writer
}

func (l gormWriter) Printf(message string, data ...interface{}) {
	logx.Errorf(message, data...)
}
