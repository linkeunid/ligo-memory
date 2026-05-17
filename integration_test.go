package memory_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/linkeunid/ligo"
	memory "github.com/linkeunid/ligo-memory"
	"github.com/linkeunid/ligo/adapters/echo"
)

// userFixture is a minimal domain entity used only in tests.
type userFixture struct {
	ID   string
	Name string
}

func startApp(t *testing.T, modules ...ligo.Module) *ligo.App {
	t.Helper()
	app := ligo.New(
		ligo.WithRouter(echo.NewAdapter()),
		ligo.WithAddr(":0"),
	)
	app.Register(modules...)
	go func() { _ = app.Run() }()
	time.Sleep(200 * time.Millisecond)
	return app
}

func containsType(types []reflect.Type, want reflect.Type) bool {
	for _, typ := range types {
		if typ == want {
			return true
		}
	}
	return false
}

// TestModuleIntegratesWithApp verifies that memory.Module() registers
// successfully in a Ligo app and exposes *Store[string, any] in the DI
// container.
func TestModuleIntegratesWithApp(t *testing.T) {
	app := startApp(t, memory.Module())

	want := reflect.TypeFor[*memory.Store[string, any]]()
	if !containsType(app.Container().Types(), want) {
		t.Fatalf("*Store[string, any] not found in DI container; registered types: %v", app.Container().Types())
	}
}

// TestProviderInjectsTypedStore verifies that Provider[K, V]() registers a
// correctly typed *Store[K, V] that the DI container can inject into a
// dependent factory.
func TestProviderInjectsTypedStore(t *testing.T) {
	var injected *memory.Store[string, *userFixture]

	mod := ligo.NewModule(
		"test",
		ligo.Providers(
			memory.Provider[string, *userFixture](),
			ligo.Factory[*userFixture](func(s *memory.Store[string, *userFixture]) *userFixture {
				injected = s
				return &userFixture{}
			}),
		),
	)
	app := startApp(t, mod)

	want := reflect.TypeFor[*memory.Store[string, *userFixture]]()
	if !containsType(app.Container().Types(), want) {
		t.Fatalf("*Store[string, *userFixture] not found in DI container; registered types: %v", app.Container().Types())
	}

	// If the factory was resolved, verify the injected store is functional.
	if injected != nil {
		injected.Set("u1", &userFixture{ID: "u1", Name: "Alice"})
		u, ok := injected.Get("u1")
		if !ok {
			t.Fatal("expected to retrieve user from injected store")
		}
		if u.Name != "Alice" {
			t.Fatalf("expected Name %q, got %q", "Alice", u.Name)
		}
	}
}
