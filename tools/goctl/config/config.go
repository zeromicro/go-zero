package config

import (
	_ "embed"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/ctx"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"gopkg.in/yaml.v2"
)

const (
	// DefaultFormat defines a default naming style
	DefaultFormat = "gozero"
	configFile    = "goctl.yaml"
)

//go:embed default.yaml
var defaultConfig []byte

// Config defines the file naming style
type (
	Config struct {
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

	External struct {
		// Model is the configuration for the model code generation.
		Model Model `yaml:"model,omitempty"`
	}

	// Model defines the configuration for the model code generation.
	Model struct {
		// TypesMap: custom Data Type Mapping Table.
		TypesMap map[string]ModelTypeMapOption `yaml:"types_map,omitempty" `
	}

	// ModelTypeMapOption custom Type Options.
	ModelTypeMapOption struct {
		// Type: valid when not using UnsignedType and NullType.
		Type string `yaml:"type"`

		// UnsignedType: valid when not using  NullType.
		UnsignedType string `yaml:"unsigned_type,omitempty"`

		// NullType: priority use.
		NullType string `yaml:"null_type,omitempty"`

		// Pkg defines the package of the custom type.
		Pkg string `yaml:"pkg,omitempty"`
	}
)

// NewConfig creates an instance for Config
func NewConfig(format string) (*Config, error) {
	if len(format) == 0 {
		format = DefaultFormat
	}
	cfg := &Config{NamingFormat: format}
	err := validate(cfg)
	return cfg, err
}

func GetExternalConfig() (*External, error) {
	var cfg External
	err := loadConfig(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadConfig(cfg *External) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfgFile, err := getConfigPath(wd)
	if err != nil {
		return err
	}
	var content []byte
	if pathx.FileExists(cfgFile) {
		content, err = os.ReadFile(cfgFile)
		if err != nil {
			return err
		}
	}
	if len(content) == 0 {
		content = append(content, defaultConfig...)
	}
	return yaml.Unmarshal(content, cfg)
}

// getConfigPath returns the configuration file path, but not create the file.
func getConfigPath(workDir string) (string, error) {
	abs, err := filepath.Abs(workDir)
	if err != nil {
		return "", err
	}

	err = pathx.MkdirIfNotExist(abs)
	if err != nil {
		return "", err
	}

	projectCtx, err := ctx.Prepare(abs)
	if err != nil {
		return "", err
	}
	return filepath.Join(projectCtx.Dir, configFile), nil
}

func validate(cfg *Config) error {
	if len(strings.TrimSpace(cfg.NamingFormat)) == 0 {
		return errors.New("missing namingFormat")
	}
	return nil
}
