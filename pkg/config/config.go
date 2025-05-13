package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rcond   RcondConfig   `yaml:"rcond"`
	Cluster ClusterConfig `yaml:"cluster"`
}

type RcondConfig struct {
	Addr     string `yaml:"addr"`
	ApiToken string `yaml:"api_token"`
}

type ClusterConfig struct {
	Enabled       bool     `yaml:"enabled"`
	NodeName      string   `yaml:"node_name"`
	SecretKey     string   `yaml:"secret_key"`
	Join          []string `yaml:"join"`
	AdvertiseAddr string   `yaml:"advertise_addr"`
	AdvertisePort int      `yaml:"advertise_port"`
	BindAddr      string   `yaml:"bind_addr"`
	BindPort      int      `yaml:"bind_port"`
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
