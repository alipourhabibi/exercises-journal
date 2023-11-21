package logger

import (
	"context"
	"fmt"
	"math"
	"os"
	"runtime"
	"time"
)

type Level int

var (
	LevelTrace Level = 5
	LevelDebug Level = 4
	LevelInfo  Level = 3
	LevelWarn  Level = 2
	LevelError Level = 1
)

var levelStr = map[Level]string{
	LevelTrace: "TRACE",
	LevelDebug: "DEBUG",
	LevelInfo:  "INFO",
	LevelWarn:  "WARN",
	LevelError: "ERROR",
}

var StrLevel = map[string]Level{
	"TRACE": LevelTrace,
	"DEBUG": LevelDebug,
	"INFO":  LevelInfo,
	"WARN":  LevelWarn,
	"ERROR": LevelError,
}

type Logger struct {
	output      *os.File
	prefix      string
	level       Level
	printCaller bool
	timeFormat  string
}

func New(file string, prefix string, timeFormat string, level Level) (*Logger, error) {
	out, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	return &Logger{
		timeFormat: timeFormat,
		output:     out,
		prefix:     prefix,
		level:      level,
	}, nil
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) GetLevel() Level {
	return l.level
}

func (l *Logger) SetOutput(out *os.File) {
	l.output = out
}

func (l *Logger) SetPrintCaller(b bool) {
	l.printCaller = b
}

func (l *Logger) log(ctx context.Context, level Level, args ...any) {
	if l.printCaller {
		// skip 2
		pc, file, line, _ := runtime.Caller(2)
		funcForPC := runtime.FuncForPC(pc)
		args = append(args, "caller_file", fmt.Sprintf("%s:%d", file, line))
		args = append(args, "caller_func", funcForPC.Name())
	}
	mod := math.Mod(float64(len(args)), 2)
	if mod != 0 {
		return
	}
	var msg string
	msg += "time=[" + time.Now().Format(l.timeFormat) + "] "
	msg += "level=" + levelStr[level] + " "
	for len(args) > 0 {
		if strArg, ok := args[0].(string); ok {
			msg += "\"" + strArg + "\""
		} else {
			msg += fmt.Sprintf("%v", args[0])
		}
		if strArg, ok := args[1].(string); ok {
			msg += "=\"" + strArg + "\""
		} else {
			msg += fmt.Sprintf("=%v", args[1])
		}
		msg += " "
		args = args[2:]
	}
	msg += "\n"
	l.output.Write([]byte(msg))
}

func (l *Logger) Trace(args ...any) {
	if l.level >= LevelTrace {
		l.log(context.Background(), LevelTrace, args...)
	}
}

func (l *Logger) Debug(args ...any) {
	if l.level >= LevelDebug {
		l.log(context.Background(), LevelDebug, args...)
	}
}

func (l *Logger) Info(args ...any) {
	if l.level >= LevelInfo {
		l.log(context.Background(), LevelInfo, args...)
	}
}

func (l *Logger) Warn(args ...any) {
	if l.level >= LevelWarn {
		l.log(context.Background(), LevelWarn, args...)
	}
}

func (l *Logger) Error(args ...any) {
	if l.level >= LevelError {
		l.log(context.Background(), LevelError, args...)
	}
}
