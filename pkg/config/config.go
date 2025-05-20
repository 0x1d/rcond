package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Hostname string        `yaml:"hostname" envconfig:"HOSTNAME"`
	Rcond    RcondConfig   `yaml:"rcond"`
	Network  NetworkConfig `yaml:"network"`
	Cluster  ClusterConfig `yaml:"cluster"`
}

type RcondConfig struct {
	Addr     string `yaml:"addr" envconfig:"RCOND_ADDR"`
	ApiToken string `yaml:"api_token" envconfig:"RCOND_API_TOKEN"`
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
	Enabled       bool     `yaml:"enabled" envconfig:"CLUSTER_ENABLED"`
	NodeName      string   `yaml:"node_name" envconfig:"CLUSTER_NODE_NAME"`
	SecretKey     string   `yaml:"secret_key" envconfig:"CLUSTER_SECRET_KEY"`
	Join          []string `yaml:"join" envconfig:"CLUSTER_JOIN"`
	AdvertiseAddr string   `yaml:"advertise_addr" envconfig:"CLUSTER_ADVERTISE_ADDR"`
	AdvertisePort int      `yaml:"advertise_port" envconfig:"CLUSTER_ADVERTISE_PORT"`
	BindAddr      string   `yaml:"bind_addr" envconfig:"CLUSTER_BIND_ADDR"`
	BindPort      int      `yaml:"bind_port" envconfig:"CLUSTER_BIND_PORT"`
	LogLevel      string   `yaml:"log_level" envconfig:"CLUSTER_LOG_LEVEL"`
}

// LoadConfig reads the configuration from a YAML file and environment variables.
func LoadConfig(filename string) (*Config, error) {
	var config Config
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	err = envconfig.Process("", &config)
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
