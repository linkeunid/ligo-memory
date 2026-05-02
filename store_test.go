package memory

import (
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	s := New[string, int]()
	if s == nil {
		t.Fatal("New returned nil")
	}
	if s.Len() != 0 {
		t.Fatalf("expected empty store, got len %d", s.Len())
	}
}

func TestSetAndGet(t *testing.T) {
	s := New[string, string]()
	s.Set("key", "value")

	got, ok := s.Get("key")
	if !ok {
		t.Fatal("expected key to exist")
	}
	if got != "value" {
		t.Fatalf("expected %q, got %q", "value", got)
	}
}

func TestGetMissing(t *testing.T) {
	s := New[string, string]()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestSetOverwrite(t *testing.T) {
	s := New[string, int]()
	s.Set("k", 1)
	s.Set("k", 2)

	got, _ := s.Get("k")
	if got != 2 {
		t.Fatalf("expected 2 after overwrite, got %d", got)
	}
}

func TestDelete(t *testing.T) {
	s := New[string, string]()
	s.Set("k", "v")

	if !s.Delete("k") {
		t.Fatal("expected Delete to return true for existing key")
	}
	if _, ok := s.Get("k"); ok {
		t.Fatal("key should not exist after Delete")
	}
}

func TestDeleteMissing(t *testing.T) {
	s := New[string, string]()
	if s.Delete("ghost") {
		t.Fatal("expected Delete to return false for missing key")
	}
}

func TestAll(t *testing.T) {
	s := New[string, int]()
	s.Set("a", 1)
	s.Set("b", 2)
	s.Set("c", 3)

	all := s.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 values, got %d", len(all))
	}
}

func TestAllEmpty(t *testing.T) {
	s := New[string, int]()
	if len(s.All()) != 0 {
		t.Fatal("expected empty slice for empty store")
	}
}

func TestKeys(t *testing.T) {
	s := New[string, int]()
	s.Set("x", 10)
	s.Set("y", 20)

	keys := s.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestLen(t *testing.T) {
	s := New[string, bool]()
	if s.Len() != 0 {
		t.Fatal("expected Len 0 on empty store")
	}
	s.Set("a", true)
	s.Set("b", false)
	if s.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", s.Len())
	}
}

func TestClear(t *testing.T) {
	s := New[string, int]()
	s.Set("a", 1)
	s.Set("b", 2)
	s.Clear()

	if s.Len() != 0 {
		t.Fatalf("expected empty store after Clear, got len %d", s.Len())
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	s := New[int, int]()
	const goroutines = 50
	const ops = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			for j := range ops {
				s.Set(i*ops+j, j)
			}
		}(i)

		go func(i int) {
			defer wg.Done()
			for j := range ops {
				s.Get(i*ops + j)
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentDelete(t *testing.T) {
	s := New[int, int]()
	for i := range 100 {
		s.Set(i, i)
	}

	var wg sync.WaitGroup
	wg.Add(100)
	for i := range 100 {
		go func(i int) {
			defer wg.Done()
			s.Delete(i)
		}(i)
	}
	wg.Wait()
}
