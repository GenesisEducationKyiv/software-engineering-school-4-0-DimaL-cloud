package configs

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Mail   Mail   `yaml:"mail"`
	DB     DB     `yaml:"db"`
	Rate   Rate   `yaml:"rate"`
	Server Server `yaml:"server"`
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
