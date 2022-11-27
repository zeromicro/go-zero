package config

import (
    {{if .useCasbin}}"github.com/suyuan32/simple-admin-core/pkg/config"
    "github.com/zeromicro/go-zero/core/stores/redis"{{end}}
    "github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth   rest.AuthConf
	{{if .useCasbin}}DatabaseConf config.DatabaseConf
    RedisConf    redis.RedisConf
	CasbinConf   config.CasbinConf{{end}}
	{{.jwtTrans}}
}
