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
$ ./httpserve --port 8001 --prefix Server -- path /home/ali --level 4
```

## Log Levels
There are 4 log Levels
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

NOTE: the log level is only used to keep track of your logic for logging. so if the level is ERROR and you do a Debug call it will log the Debug args you provide.
