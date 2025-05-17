package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Hostname string        `yaml:"hostname"`
	Rcond    RcondConfig   `yaml:"rcond"`
	Network  NetworkConfig `yaml:"network"`
	Cluster  ClusterConfig `yaml:"cluster"`
}

type RcondConfig struct {
	Addr     string `yaml:"addr"`
	ApiToken string `yaml:"api_token"`
}

type NetworkConfig struct {
	Connections []ConnectionConfig `yaml:"connections"`
}

type ConnectionConfig struct {
	Type        string `yaml:"type,omitempty"`
	UUID        string `yaml:"uuid,omitempty"`
	ID          string `yaml:"id,omitempty"`
	AutoConnect bool   `yaml:"autoconnect,omitempty"`
	SSID        string `yaml:"ssid,omitempty"`
	Mode        string `yaml:"mode,omitempty"`
	Band        string `yaml:"band,omitempty"`
	Channel     uint32 `yaml:"channel,omitempty"`
	KeyMgmt     string `yaml:"keymgmt,omitempty"`
	PSK         string `yaml:"psk,omitempty"`
	IPv4Method  string `yaml:"ipv4method,omitempty"`
	IPv6Method  string `yaml:"ipv6method,omitempty"`
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
	LogLevel      string   `yaml:"log_level"`
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
