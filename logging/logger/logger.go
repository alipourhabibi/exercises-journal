package logger

import (
	"context"
	"fmt"
	"os"
	"time"
)

var output = os.Stdout

type Level int

var (
	LevelDebug Level = 4
	LevelInfo  Level = 3
	LevelWarn  Level = 2
	LevelError Level = 1
)

var levelStr = map[Level]string{
	LevelDebug: "[DEBUG]",
	LevelInfo:  "[INFO]",
	LevelWarn:  "[WARN]",
	LevelError: "[ERROR]",
}

type Logger struct {
	prefix string
	level  Level
}

func New(out *os.File, prefix string, level Level) *Logger {
	output = out
	return &Logger{
		prefix: prefix,
		level:  level,
	}
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) GetLevel() Level {
	return l.level
}

func (l *Logger) SetOutput(out *os.File) {
	output = out
}

func (l *Logger) log(ctx context.Context, level Level, msg string, args ...any) {
	strs := []string{}
	for _, v := range args {
		strs = append(strs, fmt.Sprintf("%s", v))
	}
	msg = fmt.Sprintf(msg, args...)
	msg = fmt.Sprintf("%s %s %s", levelStr[level], l.prefix, msg)
	msg = fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
	output.Write([]byte(msg))
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(context.Background(), LevelInfo, msg, args...)
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log(context.Background(), LevelDebug, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log(context.Background(), LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log(context.Background(), LevelError, msg, args...)
}
