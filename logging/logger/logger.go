package logger

import (
	"context"
	"fmt"
	"os"
)

var output = os.Stdout

type Level int

var (
	LevelDebug Level = 8
	LevelInfo  Level = 4
	LevelWarn  Level = 2
	LevelError Level = 1
)

var levelStr = map[Level]string{
	LevelDebug: "[DEBUG]",
	LevelInfo:  "[INFO]",
	LevelWarn:  "[WARN]",
	LevelError: "[ERROR]",
}

type logger struct {
	prefix string
	level  Level
}

func New(output *os.File, prefix string, level Level) *logger {
	output = output
	return &logger{
		prefix: prefix,
		level:  level,
	}
}

func (l *logger) SetLever(level Level) {
	l.level = level
}

func (l *logger) SetOutput(out *os.File) {
	output = out
}

func (l *logger) log(ctx context.Context, level Level, msg string, args ...any) {
	strs := []string{}
	for _, v := range args {
		strs = append(strs, fmt.Sprintf("%s", v))
	}
	msg = fmt.Sprintf(msg, args...)
	msg = levelStr[level] + " " + l.prefix + " " + msg + "\n"
	output.Write([]byte(msg))
}

func (l *logger) Info(msg string, args ...any) {
	l.log(context.Background(), LevelInfo, msg, args...)
}

func (l *logger) Debug(msg string, args ...any) {
	l.log(context.Background(), LevelDebug, msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.log(context.Background(), LevelDebug, msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.log(context.Background(), LevelDebug, msg, args...)
}
