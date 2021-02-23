package config

import (
	"errors"
	"strings"
)

const (
	// DefaultFormat defines a default naming style
	DefaultFormat = "gozero"
)

const defaultYaml = `# namingFormat is used to define the naming format of the generated file name.
# just like time formatting, you can specify the formatting style through the
# two format characters go, and zero. for example: snake format you can
# define as go_zero, camel case format you can it is defined as goZero,
# and even split characters can be specified, such as go#zero. in theory,
# any combination can be used, but the prerequisite must meet the naming conventions
# of each operating system file name. if you want to independently control the file 
# naming style of the api, rpc, and model layers, you can set it through apiNamingFormat, 
# rpcNamingFormat, modelNamingFormat, and independent control is not enabled by default. 
# for more information, please see #{apiNamingFormat},#{rpcNamingFormat},#{modelNamingFormat}
# Note: namingFormat is based on snake or camel string
namingFormat: gozero
`

// Config defines the file naming style
type Config struct {
	// NamingFormat is used to define the naming format of the generated file name.
	// just like time formatting, you can specify the formatting style through the
	// two format characters go, and zero. for example: snake format you can
	// define as go_zero, camel case format you can it is defined as goZero,
	// and even split characters can be specified, such as go#zero. in theory,
	// any combination can be used, but the prerequisite must meet the naming conventions
	// of each operating system file name.
	// Note: NamingFormat is based on snake or camel string
	NamingFormat string `yaml:"namingFormat"`
}

// NewConfig creates an instance for Config
func NewConfig(format string) (*Config, error) {
	if len(format) == 0 {
		format = DefaultFormat
	}
	cfg := &Config{NamingFormat: format}
	err := validate(cfg)
	return cfg, err
}

func validate(cfg *Config) error {
	if len(strings.TrimSpace(cfg.NamingFormat)) == 0 {
		return errors.New("missing namingFormat")
	}
	return nil
}
