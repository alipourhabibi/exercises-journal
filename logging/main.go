package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alipourhabibi/exercises-journal/logging/config"
	"github.com/alipourhabibi/exercises-journal/logging/http"
	"github.com/alipourhabibi/exercises-journal/logging/logger"
)

var port = flag.Uint("port", 0, "http file server port")
var prefix = flag.String("prefix", "Server", "prefix for logs")
var route = flag.String("route", "/", "http route")
var path = flag.String("path", "/home/ali", "file path in os")
var logLevel = flag.Int("level", 3, "log level")
var outFlie = flag.String("out", "", "output file")
var printCallerstr = flag.String("printcaller", "false", "prints the caller function and file")
var configFile = flag.String("configfile", "", "config file")

func configuire() {
	flag.Parse()
	if *configFile == "" {
		*configFile = "/etc/logger/config.yaml"
	}
	var err error
	err = config.Conf.Load(*configFile)
	err = readConfigFile()
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
			config.Conf.Port = uint16(port)
		case "prefix":
			config.Conf.Prefix = f.Value.String()
		case "route":
			config.Conf.Route = f.Value.String()
		case "path":
			config.Conf.Path = f.Value.String()
		case "level":
			level, err := strconv.Atoi(f.Value.String())
			if err != nil {
				panic(err)
			}
			config.Conf.Level = level
		case "out":
			config.Conf.Out = f.Value.String()
		case "printcaller":
			pc, err := strconv.ParseBool(f.Value.String())
			if err != nil {
				panic(err)
			}
			config.Conf.Printcaller = pc
		}
	})
}

func readConfigFile() error {
	err := config.Conf.Load(*configFile)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	configuire()
	level := logger.Level(config.Conf.Level)
	customLogger, err := logger.New(config.Conf.Out, config.Conf.Prefix, config.Conf.Format.Time, level)
	if err != nil {
		panic(err)
	}
	customLogger.SetPrintCaller(config.Conf.Printcaller)
	h2 := http.NewFileServer(*customLogger, config.Conf.Path, config.Conf.Route)
	h2.SetupServer(config.Conf.Port)
	go h2.Run()

	var state byte
	const (
		waitForSignal byte = iota
		reconfigureStatus
	)
	signalCh := make(chan os.Signal, 1)
	for {
		switch state {
		case reconfigureStatus:
			// makes new logger
			// set it to http
			err := readConfigFile()
			if err != nil {
				// TODO
			}
			level := logger.Level(config.Conf.Level)
			customLogger, err = logger.New(config.Conf.Out, config.Conf.Prefix, config.Conf.Format.Time, level)
			if err != nil {
				panic(err)
			}
			customLogger.SetPrintCaller(config.Conf.Printcaller)
			h2.SetLogger(*customLogger)
			state = waitForSignal
		case waitForSignal:
			signal.Notify(signalCh,
				syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT)

			sig := <-signalCh
			log.Println("signal recieved: ", sig)
			switch sig {
			case syscall.SIGHUP:
				state = reconfigureStatus
			case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
				return
			}

		}
	}
}
