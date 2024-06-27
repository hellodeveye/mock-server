package config

import (
	"errors"
	"github.com/mcuadros/go-defaults"
	"mock-server/net"
	"os"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Name      string      `yaml:"name"`
	Port      int         `yaml:"port"`
	Endpoints []*Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	URL      string            `yaml:"url"`
	Method   string            `yaml:"method"`
	Response string            `yaml:"response"`
	Status   int               `yaml:"status"`
	Delay    int               `yaml:"delay"`
	Enabled  bool              `yaml:"enabled" default:"true"`
	Headers  map[string]string `yaml:"headers"`
}

type NacosConfig struct {
	Enabled     bool   `yaml:"enabled"`
	ServerAddr  string `yaml:"server-addr"`
	ServerPort  int64  `yaml:"server-port"`
	NamespaceId string `yaml:"namespace-id"`
}

type RouterConfig struct {
	Enabled     bool   `yaml:"enabled"`
	GatewayAddr string `yaml:"gateway-addr"`
}

type Config struct {
	Server []*Server    `yaml:"server"`
	Nacos  NacosConfig  `yaml:"nacos"`
	Router RouterConfig `yaml:"router"`
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
		// init default values
		for _, ep := range server.Endpoints {
			defaults.SetDefaults(ep)
		}
	}
	return nil
}
