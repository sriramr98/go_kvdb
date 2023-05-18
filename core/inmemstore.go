package core

import (
	"errors"
	"sync"
)

type InMemoryStore struct {
	mu   sync.RWMutex
	data map[string][]byte
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		data: make(map[string][]byte),
	}
}

var ErrKeyNotFound = errors.New("key not found")

func (s *InMemoryStore) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *InMemoryStore) Get(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]

	if !ok {
		return nil, ErrKeyNotFound
	}

	return value, nil
}

func (s *InMemoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
