package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Port        uint16 `yaml:"port"`
	Prefix      string `yaml:"prefix"`
	Route       string `yaml:"route"`
	Path        string `yaml:"path"`
	Level       int    `yaml:"level"`
	Out         string `yaml:"out"`
	Printcaller bool   `yaml:"printcaller"`
	Format      Format `yaml:"format"`
}

type Format struct {
	Time string `yaml:"time"`
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
