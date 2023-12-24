package memdb

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type MemDBConfiguration func(*MemDBService) error

type Data struct {
	Value string
	ttl   int64
}

type DB map[string]Data

type Evictor interface {
	Run(DB) error
}

type MemDBService struct {
	sync.RWMutex
	db               DB
	persist          bool
	persistInterval  int64
	eviction         bool
	evictionInterval int64
	evictors         []Evictor
	path             string
	logger           *zap.Logger
}

func New(cfgs ...MemDBConfiguration) (*MemDBService, error) {
	ms := &MemDBService{}
	for _, cfg := range cfgs {
		err := cfg(ms)
		if err != nil {
			return nil, err
		}
	}
	return ms, nil
}

func WithNewDB() MemDBConfiguration {
	return func(ms *MemDBService) error {
		ms.db = map[string]Data{}
		return nil
	}
}

func WithEvictors(evictors ...Evictor) MemDBConfiguration {
	return func(ms *MemDBService) error {
		ms.evictors = evictors
		return nil
	}
}

func WithPersist(persist bool) MemDBConfiguration {
	return func(ms *MemDBService) error {
		ms.persist = persist
		return nil
	}
}

func WithPersistInterval(interval string) MemDBConfiguration {
	return func(ms *MemDBService) error {
		iInterval, err := time.ParseDuration(interval)
		if err != nil {
			return err
		}
		ms.persistInterval = int64(iInterval)
		return nil
	}
}

func WithEviction(eviction bool) MemDBConfiguration {
	return func(ms *MemDBService) error {
		ms.eviction = eviction
		return nil
	}
}

func WithEvictionInterval(interval string) MemDBConfiguration {
	return func(ms *MemDBService) error {
		iInterval, err := time.ParseDuration(interval)
		if err != nil {
			return err
		}
		ms.evictionInterval = int64(iInterval)
		return nil
	}
}

func WithPath(path string) MemDBConfiguration {
	return func(ms *MemDBService) error {
		ms.path = path
		return nil
	}
}

func WithLogger(logger *zap.Logger) MemDBConfiguration {
	return func(ms *MemDBService) error {
		ms.logger = logger
		return nil
	}
}

func (m *MemDBService) GetPersist() bool {
	return m.persist
}

func (m *MemDBService) GetEviction() bool {
	return m.eviction
}

func (m *MemDBService) GetKey(key string) Data {
	m.RLock()
	defer m.RUnlock()
	return m.db[key]
}

func (m *MemDBService) SetKey(key string, value Data) {
	m.Lock()
	defer m.Unlock()
	m.db[key] = value
}

func (m *MemDBService) GetAllKeys() []string {
	m.RLock()
	defer m.RUnlock()
	return maps.Keys(m.db)
}

func (m *MemDBService) DelKey(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.db, key)
}

func (m *MemDBService) RunEvictor() {
	m.logger.Sugar().Debugw("start running evictor")
	ticker := time.NewTicker(time.Duration(m.evictionInterval))
	for {
		select {
		case <-ticker.C:
			for _, v := range m.evictors {
				m.logger.Sugar().Debugw("start evicting", "evictor", v)
				err := v.Run(m.db)
				if err != nil {
					m.logger.Sugar().Errorw("RunEvictor", "error", err)
					continue
				}
			}
			/*
				go func() {
					m.logger.Sugar().Debugw("start evicting", "t", t)
					for k, v := range m.db {
						// 0 means forever
						if v.ttl == 0 {
							continue
						}
						if v.ttl <= time.Now().Unix() {
							m.DelKey(k)
						}
					}
				}()
			*/
		}
	}
}

func (m *MemDBService) RunPersistor() {
	m.logger.Sugar().Debugw("start running persistor")
	ticker := time.NewTicker(time.Duration(m.persistInterval))
	for {
		select {
		case t := <-ticker.C:
			go func() {
				m.logger.Sugar().Debugw("start persisting", "t", t)
				go m.write(m.path, m.db)
			}()
		}
	}
}

func (m *MemDBService) write(path string, db DB) error {
	m.Lock()
	defer m.Unlock()
	m.logger.Sugar().Debugw("Write", "path", path, "db", db)
	return nil
}
