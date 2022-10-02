package config

import {{.authImport}}

type Config struct {
	rest.RestConf `mapstructure:",squash"`
	{{.auth}}
	{{.jwtTrans}}
}
