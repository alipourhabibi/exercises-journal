package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/alipourhabibi/exercises-journal/logging/config"
	"github.com/alipourhabibi/exercises-journal/logging/http"
	"github.com/alipourhabibi/exercises-journal/logging/logger"
)

var port = flag.Uint("port", 8000, "port which is used for the http file server")
var prefix = flag.String("prefix", "Server", "prefix for the printed logs")
var route = flag.String("route", "/", "http route which is server on browser")
var path = flag.String("path", "/usr/share/httpfileserver", "directory that the files in that are being served")
var logLevel = flag.String("level", "info", "log level")
var outFlie = flag.String("out", "/dev/stderr", "logs output file")
var printCallerstr = flag.String("printcaller", "false", "prints the caller function and file; usefull for debugging")
var configFile = flag.String("configfile", "/etc/httpfileserver/config.yaml", "path to the config file")

func configuire() {
	flag.Parse()
	if *configFile == "" {
		*configFile = "/etc/logger/config.yaml"
	}
	var err error
	err = readConfigFile()
	if err != nil {
		// config file is mandatory
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
	checkDefaults()
	return nil
}

func checkDefaults() {
	if config.Conf.Logging.Level == "" {
		config.Conf.Logging.Level = "info"
	}

	if config.Conf.Logging.Out == "" {
		config.Conf.Logging.Out = "/dev/stderr"
	}

	if config.Conf.Logging.Format == "" {
		config.Conf.Logging.Format = "2006-01-02 15:04:05"
	}

	if config.Conf.Logging.Prefix == "" {
		config.Conf.Logging.Prefix = "Server"
	}

	if config.Conf.Server.Path == "" {
		config.Conf.Server.Path = "/usr/share/httpfileserver"
	}

	if config.Conf.Server.Route == "" {
		config.Conf.Server.Route = "/"
	}

	if config.Conf.Server.Port == 0 {
		config.Conf.Server.Port = 8000
	}
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
		// Unexpected behaviout happened
		panic(err)
	}
	customLogger.SetPrintCaller(config.Conf.Logging.Printcaller)
	h2, err := http.NewFileServer(*customLogger, config.Conf.Server.Path, config.Conf.Server.Route)
	if err != nil {
		os.Exit(1)
	}
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
			customLogger = customLogger.ReloadLogger(config.Conf.Logging.Out, config.Conf.Logging.Prefix, config.Conf.Logging.Format, level)
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
			customLogger.Info("msg", "signal recieved", "signal", sig)
			switch sig {
			case syscall.SIGHUP:
				state = reconfigureStatus
			case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
				return
			}

		}
	}
}
