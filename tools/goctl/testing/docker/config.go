package docker

import (
	"time"

	"gopkg.in/yaml.v2"
)

type YamlConfig struct {
	Version  string               `yaml:"version"`
	Services map[string]*ImageCfg `yaml:"services"`
}

type Hooks struct {
	Cmd    []string `yaml:"cmd"`
	Custom string   `yaml:"custom"`
}

type ImageCfg struct {
	Image       string        `yaml:"image"`
	Ports       []string      `yaml:"ports"`
	Environment []string      `yaml:"environment"`
	Command     []string      `yaml:"command"`
	Volumes     []string      `yaml:"volumes"`
	HealthCheck *HealthyCheck `yaml:"healthcheck"`
	Hooks       []*Hooks      `yaml:"hooks"`
}

type HealthyCheck struct {
	Test     []string      `yaml:"test"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Retries  int           `yaml:"retries"`
}

func decodeConfig(data []byte) (*YamlConfig, error) {
	cfg := &YamlConfig{}
	err := yaml.Unmarshal(data, cfg)
	return cfg, err
}
