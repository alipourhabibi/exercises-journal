package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/alipourhabibi/exercises-journal/logging/config"
	"github.com/alipourhabibi/exercises-journal/logging/http"
	"github.com/alipourhabibi/exercises-journal/logging/logger"
)

var port = flag.Uint("port", 0, "http file server port")
var prefix = flag.String("prefix", "Server", "prefix for logs")
var route = flag.String("route", "/", "http route")
var path = flag.String("path", "/home/ali", "file path in os")
var logLevel = flag.String("level", "info", "log level")
var outFlie = flag.String("out", "", "output file")
var printCallerstr = flag.String("printcaller", "false", "prints the caller function and file")
var configFile = flag.String("configfile", "", "config file")

func configuire() {
	flag.Parse()
	if *configFile == "" {
		*configFile = "/etc/logger/config.yaml"
	}
	var err error
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
			config.Conf.Server.Port = uint16(port)
		case "route":
			config.Conf.Server.Route = f.Value.String()
		case "path":
			config.Conf.Server.Path = f.Value.String()

		case "prefix":
			config.Conf.Logging.Prefix = f.Value.String()
		case "level":
			config.Conf.Logging.Level = f.Value.String()
		case "out":
			config.Conf.Logging.Out = f.Value.String()
		case "printcaller":
			pc, err := strconv.ParseBool(f.Value.String())
			if err != nil {
				panic(err)
			}
			config.Conf.Logging.Printcaller = pc
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
	strLevel := config.Conf.Logging.Level
	strLevel = strings.ToUpper(strLevel)
	level, ok := logger.StrLevel[strLevel]
	if !ok {
		level = logger.LevelInfo
	}
	customLogger, err := logger.New(config.Conf.Logging.Out, config.Conf.Logging.Prefix, config.Conf.Logging.Format, level)
	if err != nil {
		panic(err)
	}
	customLogger.SetPrintCaller(config.Conf.Logging.Printcaller)
	h2 := http.NewFileServer(*customLogger, config.Conf.Server.Path, config.Conf.Server.Route)
	h2.SetupServer(config.Conf.Server.Port)
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
			strLevel := config.Conf.Logging.Level
			strLevel = strings.ToUpper(strLevel)
			level, ok := logger.StrLevel[strLevel]
			if !ok {
				level = logger.LevelInfo
			}
			customLogger, err = logger.New(config.Conf.Logging.Out, config.Conf.Logging.Prefix, config.Conf.Logging.Format, level)
			if err != nil {
				panic(err)
			}
			customLogger.SetPrintCaller(config.Conf.Logging.Printcaller)
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
