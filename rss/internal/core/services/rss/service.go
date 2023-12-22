package rss

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/alipourhabibi/exercises-journal/rss/config"
	"github.com/alipourhabibi/exercises-journal/rss/internal/core/services/memdb"
	"github.com/alipourhabibi/exercises-journal/rss/internal/core/services/server"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

type RssConfiguration func(*RssService) error

type RssService struct {
	logger      *zap.Logger
	links       []string
	interval    time.Duration
	retinterval time.Duration
	db          *memdb.MemDBService
	server      *server.ServerService
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

func WithNewMemDB() RssConfiguration {
	return func(rs *RssService) error {
		db, err := memdb.New(
			memdb.WithNewDB(),
		)
		if err != nil {
			return err
		}
		rs.db = db
		return nil
	}
}

func WithServerService(server *server.ServerService) RssConfiguration {
	return func(rs *RssService) error {
		rs.server = server
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

func (rs *RssService) getFeed(feed string) (RssFeed, error) {
	resp, err := http.Get(feed)
	if err != nil {
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		return RssFeed{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("Status Code Not 200: %d", resp.StatusCode)
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		return RssFeed{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		return RssFeed{}, err
	}

	rss := RssFeed{}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		return RssFeed{}, err
	}
	return rss, nil
}

func (rs *RssService) asyncFeedCheck(feed string) {
	rs.logger.Sugar().Debugw("asyncFeedCheck", "feed", feed)

	rss, err := rs.getFeed(feed)
	if err != nil {
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		// TODO
		return
	}

	var lastBuildParsed time.Time
	lastBuld := rs.db.GetKey(feed)
	if lastBuld != "" {
		lastBuildParsed, err = Date(lastBuld).Parse()
		if err != nil {
			rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
			// TODO
			return
		}

		newBuildParsed, err := rss.Channel.LastBuildDate.Parse()
		if err != nil {
			rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
			// TODO
			return
		}

		if lastBuildParsed.Unix() >= newBuildParsed.Unix() {
			rs.logger.Sugar().Infow("asyncFeedCheck", "status", "already checked", "feed", feed)
			return
		}
	}

	items := []Item{}
	for _, v := range rss.Channel.Item {
		if time, err := v.PubDate.Parse(); err != nil {
			rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error(), "action", "continuing reading oterhs")
			continue // TODO
		} else if time.Unix() > lastBuildParsed.Unix() {
			items = append(items, v)
		}
	}
	rs.logger.Sugar().Debugw("asyncFeedCheck", "items", items)

	// Sending to destinatio
	rs.logger.Sugar().Debugw("asyncFeedCheck", "headers", config.Conf.Http.Headers)
	itemsByte, err := json.Marshal(items)
	if err != nil {
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		return
	}
	err = rs.server.Send(itemsByte, config.Conf.Http.Headers)
	if err != nil {
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		return
	}

	rs.db.SetKey(feed, string(rss.Channel.LastBuildDate))

	rs.logger.Sugar().Debugw("rss", "last_build_date", rss.Channel.LastBuildDate)
}
