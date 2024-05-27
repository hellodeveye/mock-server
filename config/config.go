package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Name      string     `yaml:"name"`
	Port      uint64     `yaml:"port"`
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
	ServerPort  uint64 `yaml:"server-port"`
	NamespaceId string `yaml:"namespace-id"`
}

type Config struct {
	Server []Server    `yaml:"server"`
	Nacos  NacosConfig `yaml:"nacos"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
