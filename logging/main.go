package main

import (
	"flag"
	"io/fs"
	"os"

	"github.com/alipourhabibi/exercises-journal/logging/http"
	"github.com/alipourhabibi/exercises-journal/logging/logger"
)

var port = flag.Uint("port", 8000, "http file server port")
var prefix = flag.String("prefix", "Server", "prefix for logs")
var route = flag.String("route", "/", "http route")
var path = flag.String("path", "/home/ali", "file path in os")
var logLevel = flag.Int("level", 3, "log level")
var outFlie = flag.String("out", "", "output file")

func main() {
	flag.Parse()
	var f = &os.File{}
	var err error
	if *outFlie == "" {
		f = os.Stdout
	} else {
		f, err = os.OpenFile(*outFlie, 644, fs.FileMode(os.O_CREATE))
		if err != nil {
			panic(err)
		}
	}
	level := logger.Level(*logLevel)
	logger := logger.New(f, *prefix, level)
	h2 := http.NewFileServer(*logger, *path, *route)
	h2.Run(uint16(*port))
}
