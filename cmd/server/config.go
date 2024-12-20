package server

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/crazytaxii/fake-dcgm-exporter/pkg/dcgm"
)

const (
	defaultConfigPath = "/etc/fake-dcgm-exporter/config.yaml"
	defaultPort       = 9400
)

type Config struct {
	Port                uint32 `yaml:"port,omitempty"`
	*dcgm.FakeGPUConfig `yaml:",inline"`
}

func DefaultConfig(nodeName string) *Config {
	return &Config{
		Port:          defaultPort,
		FakeGPUConfig: dcgm.DefaultGPUConfig(nodeName),
	}
}

func LoadConfig(path, nodeName string) (*Config, error) {
	cfg := DefaultConfig(nodeName)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// default config
			return cfg, nil
		}
		return nil, err
	}
	return cfg, yaml.Unmarshal(data, cfg)
}
