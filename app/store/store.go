package store

import (
	"sync"
	"time"
)

type Item struct {
	Value     string
	ExpiresAt time.Time
}

type DBStore struct {
	Mu   sync.RWMutex
	Data map[string]Item
}

func (s *DBStore) Set(key, value string, ttl int64) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(time.Duration(ttl) * time.Second)
	}

	s.Data[key] = Item{
		Value:     value,
		ExpiresAt: expiresAt,
	}
}

func (s *DBStore) Get(key string) (string, bool) {
	s.Mu.RLock()
	item, ok := s.Data[key]
	s.Mu.RUnlock()

	if !ok {
		return "", false
	}
	if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
		s.Mu.Lock()
		delete(s.Data, key)
		s.Mu.Unlock()
		return "", false
	}
	return item.Value, ok
}
