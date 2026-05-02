package memory

import (
	"reflect"
	"testing"

	"github.com/linkeunid/ligo"
)

func TestProviderType(t *testing.T) {
	p := Provider[string, int]()
	want := reflect.TypeFor[*Store[string, int]]()
	if p.Type() != want {
		t.Fatalf("expected type %v, got %v", want, p.Type())
	}
}

func TestProviderDistinctTypes(t *testing.T) {
	p1 := Provider[string, int]()
	p2 := Provider[string, string]()
	if p1.Type() == p2.Type() {
		t.Fatal("providers with different value types must have different reflect types")
	}
}

func TestModuleName(t *testing.T) {
	m := Module()
	if m.Name != "ligo-memory" {
		t.Fatalf("expected module name %q, got %q", "ligo-memory", m.Name)
	}
}

func TestModuleRegistersStringAnyStore(t *testing.T) {
	m := Module()
	want := reflect.TypeFor[*Store[string, any]]()
	for _, raw := range m.Providers {
		if p, ok := raw.(ligo.Provider); ok && p.Type() == want {
			return
		}
	}
	t.Fatalf("Module must register *Store[string, any]; providers: %v", m.Providers)
}

func TestModuleProviderType(t *testing.T) {
	m := Module()
	p := m.Providers[0].(ligo.Provider)
	want := reflect.TypeFor[*Store[string, any]]()
	if p.Type() != want {
		t.Fatalf("Module provider type: expected %v, got %v", want, p.Type())
	}
}
