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
      ent.Driver(c.DatabaseConf.GetCacheDriver(c.RedisConf)),
      ent.Debug(), // debug mode
    )

    rds := c.RedisConf.NewRedis()
    if !rds.Ping() {
        logx.Error("initialize redis failed")
        return nil
    }
    
    {{end}}
	return &ServiceContext{
		Config:c,
		{{if .isEnt}}DB:     db,
		Redis:  rds,{{end}}
	}
}
