package svc

import (
	"shorturl/rpc/transform/internal/config"
	"shorturl/rpc/transform/model"

	"github.com/tal-tech/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	c     config.Config
	Model *model.ShorturlModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		c:     c,
		Model: model.NewShorturlModel(sqlx.NewMysql(c.DataSource), c.Cache, c.Table),
	}
}
