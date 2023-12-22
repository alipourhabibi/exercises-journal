package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type config struct {
	Http      Http       `yaml:"http"`
	ZapLogger zap.Config `yaml:"zaplogger"`
}

type Http struct {
	Interval      string              `yaml:"interval"`
	Destination   string              `yaml:"destination"`
	RetryInterval string              `yaml:"retry_interval"`
	Headers       map[string][]string `yaml:"headers"`
	Timeout       string              `yaml:"timeout"`
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
