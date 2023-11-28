package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Server struct {
	Port  uint16 `yaml:"port"`
	Route string `yaml:"route"`
	Path  string `yaml:"path"`
}

type Logging struct {
	Prefix      string `yaml:"prefix"`
	Level       string `yaml:"level"`
	Out         string `yaml:"out"`
	Printcaller bool   `yaml:"print_caller"`
	Format      string `yaml:"time_format"`
}

type config struct {
	Server  Server  `yaml:"server"`
	Logging Logging `yaml:"logging"`
}

var Conf = config{}

func (g *config) Load(path string) error {
	confFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(confFile, &Conf)
	if err != nil {
		return err
	}
	return nil
}
