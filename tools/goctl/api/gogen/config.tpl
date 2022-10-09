package config

{{.authImport}}

type Config struct {
	rest.RestConf
	{{.auth}}
	{{.jwtTrans}}
}
