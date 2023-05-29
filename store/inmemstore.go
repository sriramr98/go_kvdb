package store

import (
	"errors"
	"net"
	"sync"
)

type InMemoryStore[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewInMemoryStore() *InMemoryStore[string, []byte] {
	return &InMemoryStore[string, []byte]{
		data: make(map[string][]byte),
	}
}

func NewReplicationStore() *InMemoryStore[net.Conn, struct{}] {
	return &InMemoryStore[net.Conn, struct{}]{
		data: make(map[net.Conn]struct{}),
	}
}

var ErrKeyNotFound = errors.New("key not found")

func (s *InMemoryStore[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *InMemoryStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]

	if !ok {
		return value, ErrKeyNotFound
	}

	return value, nil
}

func (s *InMemoryStore[K, V]) Delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}

func (s *InMemoryStore[K, V]) GetAllKeys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]K, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}

	return keys
}

func (s *InMemoryStore[K, V]) SetAll(data map[K]V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = data
}

func (s *InMemoryStore[K, V]) GetAll() map[K]V {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data
}
