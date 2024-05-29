package config

import (
	"errors"
	"mock-server/net"
	"os"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Name      string     `yaml:"name"`
	Port      int        `yaml:"port"`
	Endpoints []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	URL      string            `yaml:"url"`
	Method   string            `yaml:"method"`
	Response string            `yaml:"response"`
	Status   int               `yaml:"status"`
	Delay    int               `yaml:"delay"`
	Enabled  bool              `yaml:"enabled"`
	Headers  map[string]string `yaml:"headers"`
}

type NacosConfig struct {
	Enabled     bool   `yaml:"enabled"`
	ServerAddr  string `yaml:"server-addr"`
	ServerPort  int64  `yaml:"server-port"`
	NamespaceId string `yaml:"namespace-id"`
}

type Config struct {
	Server []Server    `yaml:"server"`
	Nacos  NacosConfig `yaml:"nacos"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) InitConfig() error {
	for _, server := range c.Server {
		if server.Port == 0 {
			port, err := net.GetFreePort()
			if err != nil {
				return errors.New("failed to get free port")
			}
			server.Port = port
		}
	}
	return nil
}
