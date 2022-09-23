package config

{{.authImport}}

type Config struct {
	rest.RestConf `yaml:",inline"`
	{{.auth}}
	{{.jwtTrans}}
}

type ConsulConfig struct {
	Consul consul.Conf
}
