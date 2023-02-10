package svc

import (
{{.imports}}
)

type ServiceContext struct {
	Config config.Config
    {{if .isEnt}}DB     *ent.Client
    Redis  *redis.Redis
{{end}}

}

func NewServiceContext(c config.Config) *ServiceContext {
{{if .isEnt}}   db := ent.NewClient(
      ent.Log(logx.Info), // logger
      ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
      ent.Debug(), // debug mode
    )
    
    {{end}}
	return &ServiceContext{
		Config:c,
		{{if .isEnt}}DB:     db,
		Redis:  redis.MustNewRedis(c.RedisConf),{{end}}
	}
}
