package store

import (
	"errors"
	"sync"
	"time"
)

type DataObj struct {
	value     []byte
	expiresAt int // unix timestamp of when the key needs to be deleted, -1 if no expiry
}

func (do DataObj) Value() []byte {
	return do.value
}

func (do DataObj) ExpiresAt() int {
	return do.expiresAt
}

func (do DataObj) IsExpired(now time.Time) bool {
	if do.expiresAt <= 0 {
		return false
	}
	return now.After(time.UnixMilli(int64(do.expiresAt)))
}

func NewDataObj(value []byte, expiresAt int) DataObj {
	return DataObj{
		value:     value,
		expiresAt: expiresAt,
	}
}

type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string]DataObj
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string]DataObj),
	}
}

var ErrKeyNotFound = errors.New("key not found")

func (s *InMemoryStore) Set(key string, value DataObj) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *InMemoryStore) Get(key string) (DataObj, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]

	if !ok {
		return DataObj{}, ErrKeyNotFound
	}

	return value, nil
}

func (s *InMemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

// Deletes all expired keys from X random keys and return the percentage of keys deleted
func (s *InMemoryStore) DeleteExpiredKeysFromXRandom(x int) float32 {
	counter := 0
	deleted := 0

	s.mu.Lock()
	defer s.mu.Unlock()

	for key, value := range s.data {
		if value.IsExpired(time.Now()) {
			delete(s.data, key)
			deleted += 1
		}
		counter += 1
		if counter == x {
			return float32(deleted) / float32(counter)
		}
	}

	return float32(deleted) / float32(len(s.data))
}
