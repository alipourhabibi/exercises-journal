package config

import (
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type config struct {
	Http      Http       `yaml:"http"`
	ZapLogger zap.Config `yaml:"zaplogger"`
	DB        DB         `yaml:"db"`
}

type Http struct {
	Interval      string              `yaml:"interval"`
	Destination   string              `yaml:"destination"`
	RetryInterval string              `yaml:"retry_interval"`
	Headers       map[string][]string `yaml:"headers"`
	Timeout       string              `yaml:"timeout"`
}

type DB struct {
	Persist          bool   `yaml:"persist"`
	PersistInterval  string `yaml:"persist_interval"`
	Evict            bool   `yaml:"eviction"`
	EvictionInterval string `yaml:"eviction_interval"`
	DBPath           string `yaml:"db_path"`
	RetryDBPath      string `yaml:"retry_db_path"`
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
