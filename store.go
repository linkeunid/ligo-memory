// Package memory provides a simple, thread-safe in-memory key-value store
// for the Ligo framework, designed for fast testing and local development
// without any external dependencies.
package memory

import "sync"

// Store is a generic, thread-safe in-memory key-value store.
//
// K must be a comparable type (e.g. string, int, UUID alias).
// V can be any type, typically a pointer to a domain entity.
// All methods are safe for concurrent use.
//
// Example:
//
//	store := memory.New[string, *entity.User]()
//	store.Set("u1", &entity.User{ID: "u1", Name: "Alice"})
//
//	user, ok := store.Get("u1")
//	users := store.All()
type Store[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

// New creates a new empty Store.
//
// Example:
//
//	store := memory.New[string, *entity.User]()
func New[K comparable, V any]() *Store[K, V] {
	return &Store[K, V]{data: make(map[K]V)}
}

// Get retrieves the value associated with key.
// Returns the value and true if found, the zero value and false otherwise.
func (s *Store[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

// Set stores value under key, overwriting any existing entry.
func (s *Store[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Delete removes the entry for key.
// Returns true if the key existed and was removed, false if it was not found.
func (s *Store[K, V]) Delete(key K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[key]
	delete(s.data, key)
	return ok
}

// All returns a snapshot of all stored values as a slice.
// The order of elements is not guaranteed.
func (s *Store[K, V]) All() []V {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]V, 0, len(s.data))
	for _, v := range s.data {
		out = append(out, v)
	}
	return out
}

// Keys returns a snapshot of all stored keys as a slice.
// The order of elements is not guaranteed.
func (s *Store[K, V]) Keys() []K {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]K, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of entries currently in the store.
func (s *Store[K, V]) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.data)
}

// Clear removes all entries from the store.
func (s *Store[K, V]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.data)
}
