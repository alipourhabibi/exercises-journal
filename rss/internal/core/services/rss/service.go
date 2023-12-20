package rss

import (
	"os"
	"time"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type RssConfiguration func(*RssService) error

type RssService struct {
	logger      *zap.Logger
	links       []string
	interval    time.Duration
	retinterval time.Duration
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

func WithFeeds(file string) RssConfiguration {
	return func(rs *RssService) error {
		f, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		var rss map[string][]string
		err = yaml.Unmarshal(f, &rss)
		if err != nil {
			return err
		}
		rs.links = rss["links"]
		return nil
	}
}

func WithInterval(interval string) RssConfiguration {
	return func(rs *RssService) error {
		duration, err := time.ParseDuration(interval)
		if err != nil {
			return err
		}
		rs.interval = duration
		return nil
	}
}

func WithRetInterval(interval string) RssConfiguration {
	return func(rs *RssService) error {
		duration, err := time.ParseDuration(interval)
		if err != nil {
			return err
		}
		rs.retinterval = duration
		return nil
	}
}

func (rs *RssService) Serve() error {
	rs.logger.Sugar().Debugw("Rss", "files", rs.links)
	rs.logger.Sugar().Debugw("Intervals", "interval", rs.interval, "retry-interval", rs.retinterval)

	for _, v := range rs.links {
		go rs.asyncFeedCheck(v)
	}

	ticker := time.NewTicker(rs.interval)
	for {
		select {
		case <-ticker.C:
			for _, v := range rs.links {
				go rs.asyncFeedCheck(v)
			}
		}
	}
}

func (rs *RssService) asyncFeedCheck(feed string) {
	rs.logger.Sugar().Debugw("asyncFeedCheck", "feed", feed)
}
