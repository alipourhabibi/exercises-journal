package logger

import (
	"github.com/alipourhabibi/exercises-journal/rss/config"
	"go.uber.org/zap"
)

type LoggerConfiguration func(*LoggerService) error

type LoggerService struct {
	logger *zap.Logger
}

func New(cfgs ...LoggerConfiguration) (*LoggerService, error) {
	ls := &LoggerService{}
	for _, cfg := range cfgs {
		err := cfg(ls)
		if err != nil {
			return nil, err
		}
	}
	return ls, nil

}

func WithConfig() LoggerConfiguration {
	return func(ls *LoggerService) error {
		logger, err := config.Conf.ZapLogger.Build()
		if err != nil {
			return err
		}
		ls.logger = logger
		return nil
	}
}

func (ls *LoggerService) Logger() *zap.Logger {
	return ls.logger
}
