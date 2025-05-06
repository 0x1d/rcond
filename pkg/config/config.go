package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rcond RcondConfig `yaml:"rcond"`
}

type RcondConfig struct {
	Addr     string `yaml:"addr"`
	ApiToken string `yaml:"api_token"`
}

func LoadConfig(path string) (*Config, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveConfig(path string, config *Config) error {
	yamlFile, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, yamlFile, 0644)
}
