package config

import (
    {{.authImport}}
)

type Config struct {
    rest.RestConf
	{{.rabbitmqConfig}}
}
