package configs

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DB       DB       `yaml:"db"`
	Rate     Rate     `yaml:"rate"`
	Server   Server   `yaml:"server"`
	RabbitMQ RabbitMQ `yaml:"rabbitmq"`
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
