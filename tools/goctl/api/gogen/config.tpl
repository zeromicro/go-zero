package config

import (
    {{if .useCasbin}}"github.com/suyuan32/simple-admin-common/plugins/casbin"
    "github.com/suyuan32/simple-admin-common/config"
    "github.com/zeromicro/go-zero/core/stores/redis"{{else}}{{if .useEnt}}"github.com/suyuan32/simple-admin-common/config"{{end}}{{end}}
    "github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth   rest.AuthConf
	{{if .useCasbin}}DatabaseConf config.DatabaseConf
    RedisConf    redis.RedisConf
	CasbinConf   casbin.CasbinConf{{else}}{{if .useEnt}}DatabaseConf config.DatabaseConf{{end}}{{end}}
	{{.jwtTrans}}
}
