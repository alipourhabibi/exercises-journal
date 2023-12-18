package rss

import "go.uber.org/zap"

type RssConfiguration func(*RssService) error

type RssService struct {
	logger *zap.Logger
}

func New(cfgs ...RssConfiguration) (*RssService, error) {
	rs := &RssService{}
	for _, cfg := range cfgs {
		err := cfg(rs)
		if err != nil {
			return nil, err
		}
	}
	return rs, nil
}

func WithLogger(logger *zap.Logger) RssConfiguration {
	return func(rs *RssService) error {
		rs.logger = logger
		return nil
	}
}

func (rs *RssService) Serve() error {
	rs.logger.Sugar().Debugw("Serve", "status", "starting...")
	return nil
}
