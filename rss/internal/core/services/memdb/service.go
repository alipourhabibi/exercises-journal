package memdb

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

type MemDBConfiguration func(*MemDBService) error

type Data struct {
	T   time.Time
	ttl int64
}

type LRUCache struct {
	dataMap   map[string]Data
	dataOrder []string
	capacity  int
}

func (lru *LRUCache) updateOrder(key string) {
	// Remove the key from the current position
	for i, k := range lru.dataOrder {
		if k == key {
			lru.dataOrder = append(lru.dataOrder[:i], lru.dataOrder[i+1:]...)
			break
		}
	}

	// Add the key to the end
	lru.dataOrder = append(lru.dataOrder, key)
}

type Evictor interface {
	Run(LRUCache) error
}

type MemDBService struct {
	sync.RWMutex
	db               LRUCache
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

func WithNewDB(capacity int, initialKeys []string) MemDBConfiguration {
	return func(ms *MemDBService) error {
		if initialKeys == nil {
			initialKeys = []string{}
		}
		db := map[string]Data{}
		for _, v := range initialKeys {
			db[v] = Data{
				T:   time.Time{},
				ttl: 0,
			}
		}
		ms.db = LRUCache{
			dataMap:   db,
			dataOrder: initialKeys,
			capacity:  capacity,
		}
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

func (m *MemDBService) GetKey(key string) (Data, bool) {
	m.RLock()
	defer m.RUnlock()
	data, ok := m.db.dataMap[key]
	if ok {
		m.db.updateOrder(key)
		return data, ok
	}
	return data, false
}

func (m *MemDBService) SetKey(key string, value Data) {
	m.Lock()
	defer m.Unlock()
	/*
		value.t = time.Now()
		m.db[key] = value
	*/
	if len(m.db.dataOrder) >= m.db.capacity {
		// Remove the oldest item
		oldestKey := m.db.dataOrder[0]
		delete(m.db.dataMap, oldestKey)
		m.db.dataOrder = m.db.dataOrder[1:]
	}

	// Add the new item
	if _, ok := m.db.dataMap[key]; ok {
		return
	}
	m.db.dataMap[key] = value
	m.db.dataOrder = append(m.db.dataOrder, key)
	m.logger.Sugar().Debugw("SetKey", "map", m.db.dataMap, "ordered", m.db.dataOrder)
}

func (m *MemDBService) GetAllKeys() []string {
	m.RLock()
	defer m.RUnlock()
	return maps.Keys(m.db.dataMap)
}

func (m *MemDBService) DelKey(key string) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.db.dataMap[key]; ok {
		// Delete the key from the map
		delete(m.db.dataMap, key)

		// Remove the key from the order slice
		for i, k := range m.db.dataOrder {
			if k == key {
				m.db.dataOrder = append(m.db.dataOrder[:i], m.db.dataOrder[i+1:]...)
				break
			}
		}
	}
}

// Delete entry with the given and key and its index int the dataOrder
func (m *MemDBService) DelKeyWithIndexKey(key string, index int) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.db.dataMap[key]; ok {
		if key != m.db.dataOrder[index] {
			// Key and index does not match
			return
		}
		// Delete the key from the map
		delete(m.db.dataMap, key)

		m.db.dataOrder = append(m.db.dataOrder[:index], m.db.dataOrder[index+1:]...)
	}
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

// write persists the data into the disk
func (m *MemDBService) write(path string, db LRUCache) error {
	m.Lock()
	defer m.Unlock()
	m.logger.Sugar().Debugw("Write", "path", path, "db", db)
	return nil
}
