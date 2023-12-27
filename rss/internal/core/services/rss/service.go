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
	logger *zap.Logger
	// links       []string
	interval    time.Duration
	retinterval time.Duration
	db          *memdb.MemDBService
	server      *server.ServerService
	retryDB     *memdb.MemDBService
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
		links := rss["links"]
		db, err := memdb.New(
			memdb.WithNewDB(int(config.Conf.DB.MaxRetDBData), links),
			memdb.WithPersist(false),
			memdb.WithEviction(false),
			memdb.WithPath(config.Conf.DB.DBPath),
			memdb.WithLogger(rs.logger),
		)
		if err != nil {
			return err
		}
		rs.db = db
		if db.GetEviction() {
			go rs.db.RunEvictor()
		}
		if db.GetPersist() {
			go rs.db.RunPersistor()
		}
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

type evictorMax struct{}

func (e evictorMax) Run(m memdb.LRUCache) error {
	return nil
}

func WithNewRetryMemDB() RssConfiguration {
	return func(rs *RssService) error {
		evictorMaxSize := evictorMax{}
		db, err := memdb.New(
			memdb.WithNewDB(int(config.Conf.DB.MaxRetDBData), nil),
			memdb.WithPersist(config.Conf.DB.Persist),
			memdb.WithPersistInterval(config.Conf.DB.PersistInterval),
			memdb.WithEviction(config.Conf.DB.Evict),
			memdb.WithEvictionInterval(config.Conf.DB.EvictionInterval),
			memdb.WithPath(config.Conf.DB.RetryDBPath),
			memdb.WithLogger(rs.logger),
			memdb.WithEvictors(evictorMaxSize),
		)
		if err != nil {
			return err
		}
		rs.retryDB = db
		if db.GetEviction() {
			go rs.retryDB.RunEvictor()
		}
		if db.GetPersist() {
			go rs.retryDB.RunPersistor()
		}
		return nil
	}
}

func WithServerService(server *server.ServerService) RssConfiguration {
	return func(rs *RssService) error {
		rs.server = server
		return nil
	}
}

func (rs *RssService) SetLogger(logger *zap.Logger) {
	rs.logger = logger
}

func (rs *RssService) SetServer(server *server.ServerService) {
	rs.server = server
}

func (rs *RssService) SetRetDB() {
	evictorMaxSize := evictorMax{}
	db, err := memdb.New(
		memdb.WithNewDB(int(config.Conf.DB.MaxRetDBData), nil),
		memdb.WithPersist(config.Conf.DB.Persist),
		memdb.WithPersistInterval(config.Conf.DB.PersistInterval),
		memdb.WithEviction(config.Conf.DB.Evict),
		memdb.WithEvictionInterval(config.Conf.DB.EvictionInterval),
		memdb.WithPath(config.Conf.DB.RetryDBPath),
		memdb.WithLogger(rs.logger),
		memdb.WithEvictors(evictorMaxSize),
	)
	if err != nil {
		rs.logger.Sugar().Errorw("SetRetDB", "error", err)
		return
	}
	rs.retryDB = db
	if db.GetEviction() {
		go rs.retryDB.RunEvictor()
	}
	if db.GetPersist() {
		go rs.retryDB.RunPersistor()
	}
}

func (rs *RssService) SetInterval(interval string) {
	duration, err := time.ParseDuration(interval)
	if err != nil {
		rs.logger.Sugar().Errorw("SetInterval", "error", err)
		return
	}
	rs.interval = duration
}

func (rs *RssService) SetRetInterval(interval string) {
	duration, err := time.ParseDuration(interval)
	if err != nil {
		rs.logger.Sugar().Errorw("SetRetInterval", "error", err)
		return
	}
	rs.retinterval = duration
}

// TODO Think
func (rs *RssService) SetNewFeeds(file string) {
	f, err := os.ReadFile(file)
	if err != nil {
		rs.logger.Sugar().Errorw("SetNewFeeds", "error", err)
		return
	}
	var rss map[string][]string
	err = yaml.Unmarshal(f, &rss)
	if err != nil {
		rs.logger.Sugar().Errorw("SetNewFeeds", "error", err)
		return
	}
	links := rss["links"]
	db, err := memdb.New(
		memdb.WithNewDB(int(config.Conf.DB.MaxRetDBData), links),
		memdb.WithPersist(false),
		memdb.WithEviction(false),
		memdb.WithPath(config.Conf.DB.DBPath),
		memdb.WithLogger(rs.logger),
	)
	if err != nil {
		rs.logger.Sugar().Errorw("SetNewFeeds", "error", err)
		return
	}
	rs.db = db
	if db.GetEviction() {
		go rs.db.RunEvictor()
	}
	if db.GetPersist() {
		go rs.db.RunPersistor()
	}
}

func (rs *RssService) Serve(ch chan struct{}) {
	go rs.retryAll()

	rs.logger.Sugar().Debugw("Rss", "files", rs.db.GetAllKeys())
	rs.logger.Sugar().Debugw("Intervals", "interval", rs.interval, "retry-interval", rs.retinterval)

	for _, v := range rs.db.GetAllKeys() {
		go rs.asyncFeedCheck(v)
	}

	ticker := time.NewTicker(rs.interval)
	for {
		select {
		case <-ticker.C:
			for _, v := range rs.db.GetAllKeys() {
				go rs.asyncFeedCheck(v)
			}
		case <-ch:
			return
		}
	}
}

func (rs *RssService) retryAll() {
	ticker := time.NewTicker(rs.retinterval)
	for {
		select {
		case <-ticker.C:
			for _, v := range rs.retryDB.GetAllKeys() {
				rs.logger.Sugar().Infow("retry", "feed", v)
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
		rs.db.DelKey(feed)
		rs.retryDB.SetKey(feed, memdb.Data{
			T: time.Now(),
		})
		return
	}

	var lastBuildParsed time.Time
	lastBuld, ok := rs.db.GetKey(feed)
	if ok {
		lastBuildParsed = lastBuld.T
		if err != nil {
			rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
			return
		}

		newBuildParsed, err := rss.Channel.LastBuildDate.Parse()
		if err != nil {
			rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
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
			continue
		} else if time.Unix() > lastBuildParsed.Unix() {
			items = append(items, v)
		}
	}
	// rs.logger.Sugar().Debugw("asyncFeedCheck", "items", items)

	// Sending to destinatio
	rs.logger.Sugar().Debugw("asyncFeedCheck", "headers", config.Conf.Http.Headers)
	itemsByte, err := json.Marshal(items)
	if err != nil {
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		return
	}
	err = rs.server.Send(itemsByte, config.Conf.Http.Headers)
	if err != nil {
		rs.logger.Sugar().Debugw("DB", "items", rs.retryDB)
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		rs.db.DelKey(feed)
		rs.retryDB.SetKey(feed, memdb.Data{
			T: time.Now(),
		})
		return
	}

	rs.logger.Sugar().Debugw("setting feed to db", "feed", feed)
	date, err := rss.Channel.LastBuildDate.Parse()
	if err != nil {
		rs.logger.Sugar().Errorw("asyncFeedCheck", "error", err.Error())
		return
	}
	rs.db.SetKey(feed, memdb.Data{
		T: date,
	})
	// removing from retry as it is ok
	// delete if it was success
	// TODO Possible better solution using channels to communicate
	rs.retryDB.DelKey(feed)

	rs.logger.Sugar().Debugw("rss", "last_build_date", rss.Channel.LastBuildDate)
}
