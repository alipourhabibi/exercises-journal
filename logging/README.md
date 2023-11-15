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
$ ./httpserve --port 8001 --prefix Server -- path /home/ali --level 4 --printcaller true
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
- Change log format to logfmt
- Add Trace log level
- Add option to print caller file and fucntion
- Log levels will be checked in log lib itself
- Change default output to Stderr
