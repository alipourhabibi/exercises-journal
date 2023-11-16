package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/alipourhabibi/exercises-journal/logging/http"
	"github.com/alipourhabibi/exercises-journal/logging/logger"
	yaml "gopkg.in/yaml.v2"
)

var port = flag.Uint("port", 0, "http file server port")
var prefix = flag.String("prefix", "Server", "prefix for logs")
var route = flag.String("route", "/", "http route")
var path = flag.String("path", "/home/ali", "file path in os")
var logLevel = flag.Int("level", 3, "log level")
var outFlie = flag.String("out", "", "output file")
var printCaller = flag.Bool("printcaller", false, "prints the caller function and file")
var configFile = flag.String("configfile", "", "config file")

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

var conf = config{}

func main() {
	flag.Parse()
	if *configFile == "" {
		*configFile = "config.yaml"
	}
	var err error
	confFile, err := os.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(confFile, &conf)
	if err != nil {
		panic(err)
	}
	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "port":
			port, err := strconv.ParseUint(f.Value.String(), 10, 16)
			if err != nil {
				panic(err)
			}
			conf.Port = uint16(port)
		case "prefix":
			conf.Prefix = f.Value.String()
		case "route":
			conf.Route = f.Value.String()
		case "path":
			conf.Path = f.Value.String()
		case "level":
			level, err := strconv.Atoi(f.Value.String())
			if err != nil {
				panic(err)
			}
			conf.Level = level
		case "out":
			conf.Out = f.Value.String()
		case "printcaller":
			pc, err := strconv.ParseBool(f.Value.String())
			if err != nil {
				panic(err)
			}
			conf.Printcaller = pc
		}
	})
	level := logger.Level(conf.Level)
	logger, err := logger.New(conf.Out, conf.Prefix, conf.Format.Time, level)
	logger.SetPrintCaller(conf.Printcaller)
	h2 := http.NewFileServer(*logger, conf.Path, conf.Route)
	h2.Run(uint16(conf.Port))
}
