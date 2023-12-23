package memdb

import "golang.org/x/exp/maps"

type MemDBConfiguration func(*MemDBService) error

type MemDBService struct {
	db map[string]string
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
		ms.db = map[string]string{}
		return nil
	}
}

func (m *MemDBService) GetKey(key string) string {
	return m.db[key]
}

func (m *MemDBService) SetKey(key string, value string) {
	m.db[key] = value
}

func (m *MemDBService) GetAllKeys() []string {
	return maps.Keys(m.db)
}

func (m *MemDBService) DelKey(key string) {
	delete(m.db, key)
}
