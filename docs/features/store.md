# Store

`Store[K, V]` is a generic, thread-safe in-memory key-value store. It is the core type of the `ligo-memory` package and covers the full CRUD surface expected by a repository.

---

## Type parameters

| Parameter | Constraint | Description |
|-----------|------------|-------------|
| `K` | `comparable` | Key type — any type that supports `==` (e.g. `string`, `int`, UUID alias) |
| `V` | `any` | Value type — typically a pointer to a domain entity |

---

## Creating a store

```go
// Typed store — preferred for real repositories
store := memory.New[string, *entity.User]()

// Untyped store — useful for generic caches or tests
store := memory.New[string, any]()
```

---

## API reference

### Get

```go
func (s *Store[K, V]) Get(key K) (V, bool)
```

Returns the value and `true` if `key` exists, or the zero value and `false` otherwise.

```go
user, ok := store.Get("u1")
if !ok {
    // not found
}
```

### Set

```go
func (s *Store[K, V]) Set(key K, value V)
```

Stores `value` under `key`, overwriting any existing entry.

```go
store.Set("u1", &entity.User{ID: "u1", Name: "Alice"})
```

### Delete

```go
func (s *Store[K, V]) Delete(key K) bool
```

Removes `key`. Returns `true` if the key existed, `false` if not found.

```go
deleted := store.Delete("u1") // true
deleted  = store.Delete("u1") // false — already gone
```

### All

```go
func (s *Store[K, V]) All() []V
```

Returns a snapshot of all values. Order is not guaranteed.

```go
users := store.All() // []*entity.User
```

### Keys

```go
func (s *Store[K, V]) Keys() []K
```

Returns a snapshot of all keys. Order is not guaranteed.

### Len

```go
func (s *Store[K, V]) Len() int
```

Returns the number of entries currently in the store.

### Clear

```go
func (s *Store[K, V]) Clear()
```

Removes all entries. Useful between test cases.

```go
t.Cleanup(store.Clear)
```

---

## Ligo DI integration

### Provider[K, V]() — typed store

`Provider[K, V]()` returns a `ligo.Provider` that registers a `*Store[K, V]` singleton in the DI container. Other factories in the same module receive it as an injected parameter.

```go
func UserModule() ligo.Module {
    return ligo.NewModule("user",
        ligo.Providers(
            memory.Provider[string, *entity.User](),          // registers *Store[string, *entity.User]
            ligo.Factory[repository.UserRepository](NewRepo), // injected ↑
            ligo.Factory[*usecase.UserUseCase](usecase.New),
        ),
        ligo.Controllers(controller.NewUserController),
    )
}

func NewRepo(store *memory.Store[string, *entity.User]) repository.UserRepository {
    return &userRepo{store: store}
}
```

### Module() — zero-config drop-in

`Module()` registers a `*Store[string, any]` singleton. Use it when you only need a single untyped key-value store and don't want to declare a provider yourself.

```go
app.Register(memory.Module(), myModule())
```

The consuming module can then declare a dependency on `*memory.Store[string, any]`:

```go
ligo.Factory[*myService](func(s *memory.Store[string, any]) *myService {
    return &myService{cache: s}
})
```

---

## Using Store in a repository

A typical pattern is to embed the store inside a struct that implements a domain repository interface:

```go
type userRepo struct {
    store *memory.Store[string, *entity.User]
}

func NewUserRepository(store *memory.Store[string, *entity.User]) repository.UserRepository {
    return &userRepo{store: store}
}

func (r *userRepo) FindByID(id string) (*entity.User, bool) {
    return r.store.Get(id)
}

func (r *userRepo) FindAll() []*entity.User {
    return r.store.All()
}

func (r *userRepo) Create(name, email string) *entity.User {
    u := &entity.User{ID: newUUID(), Name: name, Email: email}
    r.store.Set(u.ID, u)
    return u
}

func (r *userRepo) Delete(id string) bool {
    return r.store.Delete(id)
}
```

---

## Using Store in tests

Because the store has no external dependencies, it is trivial to use directly in unit tests — no mocking required.

```go
func TestGetUserByID(t *testing.T) {
    store := memory.New[string, *entity.User]()
    store.Set("u1", &entity.User{ID: "u1", Name: "Alice"})

    repo := NewUserRepository(store)
    user, err := repo.FindByID("u1")
    if err != nil {
        t.Fatal(err)
    }
    if user.Name != "Alice" {
        t.Fatalf("got %q", user.Name)
    }
}
```

Use `t.Cleanup(store.Clear)` to reset state between sub-tests:

```go
func TestUserRepo(t *testing.T) {
    store := memory.New[string, *entity.User]()
    repo := NewUserRepository(store)

    t.Run("create", func(t *testing.T) {
        t.Cleanup(store.Clear)
        u := repo.Create("Bob", "bob@example.com")
        if u.ID == "" {
            t.Fatal("expected non-empty ID")
        }
    })

    t.Run("not found", func(t *testing.T) {
        t.Cleanup(store.Clear)
        _, ok := repo.FindByID("ghost")
        if ok {
            t.Fatal("expected not found")
        }
    })
}
```

---

## Concurrency

All methods acquire a `sync.RWMutex` internally — reads share a lock, writes take an exclusive lock. The store is safe to use from multiple goroutines without additional synchronisation.
