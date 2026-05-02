package memory

import "sync"

// Store is a generic, thread-safe in-memory key-value store.
// Designed for fast testing and local development without external dependencies.
type Store[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// New creates a new empty Store.
func New[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{data: make(map[K]V)}
}

// Get retrieves a value by key. Returns the value and true if found.
func (s *Store[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

// Set stores a value under the given key.
func (s *Store[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Delete removes a key. Returns true if the key existed.
func (s *Store[K, V]) Delete(key K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[key]; !ok {
		return false
	}
	delete(s.data, key)
	return true
}

// All returns all stored values as a slice.
func (s *Store[K, V]) All() []V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]V, 0, len(s.data))
	for _, v := range s.data {
		out = append(out, v)
	}
	return out
}

// Keys returns all stored keys as a slice.
func (s *Store[K, V]) Keys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]K, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of entries.
func (s *Store[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// Clear removes all entries.
func (s *Store[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[K]V)
}
