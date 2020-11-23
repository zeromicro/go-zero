package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tal-tech/go-zero/tools/goctl/util"
	"gopkg.in/yaml.v2"
)

const (
	configFile    = "config.yaml"
	configFolder  = "config"
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

func NewConfig(format string) (*Config, error) {
	if len(format) == 0 {
		format = DefaultFormat
	}
	cfg := &Config{NamingFormat: format}
	err := validate(cfg)
	return cfg, err
}

func InitOrGetConfig() (*Config, error) {
	var (
		defaultConfig Config
	)
	err := yaml.Unmarshal([]byte(defaultYaml), &defaultConfig)
	if err != nil {
		return nil, err
	}

	goctlHome, err := util.GetGoctlHome()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(goctlHome, configFolder)
	configFilename := filepath.Join(configDir, configFile)
	if util.FileExists(configFilename) {
		data, err := ioutil.ReadFile(configFilename)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &defaultConfig)
		if err != nil {
			return nil, err
		}

		err = validate(&defaultConfig)
		if err != nil {
			return nil, err
		}

		return &defaultConfig, nil
	}

	err = util.MkdirIfNotExist(configDir)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(configFilename)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
	}()

	_, err = f.WriteString(defaultYaml)
	if err != nil {
		return nil, err
	}

	err = validate(&defaultConfig)
	if err != nil {
		return nil, err
	}

	return &defaultConfig, nil
}

func validate(cfg *Config) error {
	if len(strings.TrimSpace(cfg.NamingFormat)) == 0 {
		return errors.New("missing namingFormat")
	}
	return nil
}
