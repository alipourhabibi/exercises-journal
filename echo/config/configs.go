package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Confs config

type config struct {
	Server server `yaml:"server"`
	Logger logger `yaml:"logger"`
}

type server struct {
	Port uint `yaml:"port"`
}

type logger struct {
	AddSource bool   `yaml:"add_source"`
	Level     string `yaml:"level"`
}

func Load(path string) error {
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(f, &Confs)
}
