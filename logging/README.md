# HTTP-FILE-SERVE
First assignment: \
A simple http file server which will get a folder and we will serve the content \

## Goals:
The main goal for this project is logging.
Important notes:
- Error handling is IMPORTANT!
- Project should be config easily
- Can use library for http server but the reason for choice is important

## Http Server:
Library used in this project is the standard golang library for http which is "net/http" \
Also some codes are use from [http-file-server](https://github.com/sgreben/http-file-server) repo but becase the functions where private i could not use the repo directly in my codes so i used some code snippets from it.

## How to use
Clone the repo
```sh
$ go build -o httpserve .
$ ./httpserve --port 8001 --prefix Server --path /usr/share/httpfileserver --level debug --printcaller true
```
docker
```sh
$ docker run -it --network=host ghcr.io/alipourhabibi/exercise-journals-logging:latest /bin/httpfileserver --port 8000 --path /path/to/file --level debug --printcaller true
```

## Log Levels
We have 5 log Levels in this repo
- TRACE
- DEBUG
- INFO
- WARN
- ERROR

## Use in you code
```go
logger := logger.New(os.Stdout, "PREFIX", logger.LevelDebug)
logger.Info("%s", "SOME TEXT")
logger.Debug("%s", "SOME OTHER TEXT")
```

## Changelogs
November-15, 2023
- Change log format to logfmt
- Add Trace log level
- Add option to print caller file and fucntion
- Log levels will be checked in log lib itself
- Change default output to Stderr

November-16, 2023
- Add config file

November-21, 2023
- Change log level to string in config
- Decrease Docker image size
- Add dynamic conifg reload with SIGHUP
- Change config file style
- Change docker image files folder based on linux hierarchy

## Config file
We have a configuration file that exists near the cli tool which configs the app, \
The flags will override config in config.yaml file. \

default of logging.out is /dev/stderr which is standard error

config file is yaml, because,
- It is human-friendly
- It is a superset of JSON
- It is very convenient format for specifying hierarchical configuration data
- Compact syntax
Which make it a good choice for config files.

The tool will panic if there is no config file.
