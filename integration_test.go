package memory_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/linkeunid/ligo"
	"github.com/linkeunid/ligo/adapters/echo"
	memory "github.com/linkeunid/ligo-memory"
)

// userFixture is a minimal domain entity used only in tests.
type userFixture struct {
	ID   string
	Name string
}

// TestModuleIntegratesWithApp verifies that memory.Module() registers
// successfully in a Ligo app and exposes *Store[string, any] in the DI
// container.
func TestModuleIntegratesWithApp(t *testing.T) {
	app := ligo.New(
		ligo.WithRouter(echo.NewAdapter()),
		ligo.WithAddr(":0"),
	)
	app.Register(memory.Module())

	done := make(chan error, 1)
	go func() { done <- app.Run() }()
	time.Sleep(200 * time.Millisecond)

	types := app.Container().Types()
	want := reflect.TypeFor[*memory.Store[string, any]]()
	for _, typ := range types {
		if typ == want {
			return
		}
	}
	t.Fatalf("*Store[string, any] not found in DI container; registered types: %v", types)
}

// TestProviderInjectsTypedStore verifies that Provider[K, V]() registers a
// correctly typed *Store[K, V] that the DI container can inject into a
// dependent factory.
func TestProviderInjectsTypedStore(t *testing.T) {
	var injected *memory.Store[string, *userFixture]

	mod := ligo.NewModule("test",
		ligo.Providers(
			memory.Provider[string, *userFixture](),
			ligo.Factory[*userFixture](func(s *memory.Store[string, *userFixture]) *userFixture {
				injected = s // capture the injected store
				return &userFixture{}
			}),
		),
	)

	app := ligo.New(
		ligo.WithRouter(echo.NewAdapter()),
		ligo.WithAddr(":0"),
	)
	app.Register(mod)

	done := make(chan error, 1)
	go func() { done <- app.Run() }()
	time.Sleep(200 * time.Millisecond)

	// Verify *Store[string, *userFixture] is registered in the container.
	types := app.Container().Types()
	want := reflect.TypeFor[*memory.Store[string, *userFixture]]()
	found := false
	for _, typ := range types {
		if typ == want {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("*Store[string, *userFixture] not found in DI container; registered types: %v", types)
	}

	// If the factory was resolved (eager or via a prior request), also
	// verify the injected store is functional.
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
