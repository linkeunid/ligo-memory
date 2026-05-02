package memory

import "github.com/linkeunid/ligo"

// Provider returns a [ligo.Provider] that registers a *[Store][K, V] as a
// singleton in the DI container. Other factories declared in the same module
// can then receive the store as an injected parameter.
//
// Use this inside your module's [ligo.Providers] list:
//
//	func User() ligo.Module {
//	    return ligo.NewModule("user",
//	        ligo.Providers(
//	            memory.Provider[string, *entity.User](),
//	            ligo.Factory[repository.UserRepository](NewUserRepository),
//	        ),
//	    )
//	}
//
// Accept the store in your constructor:
//
//	func NewUserRepository(store *memory.Store[string, *entity.User]) repository.UserRepository {
//	    return &UserRepository{store: store}
//	}
func Provider[K comparable, V any]() ligo.Provider {
	return ligo.Factory[*Store[K, V]](func() *Store[K, V] {
		return New[K, V]()
	})
}

// Module returns a Ligo module that registers a *[Store][string, any] as a
// singleton provider via DI. It is a convenient zero-config drop-in for
// applications that only need a single untyped key-value store.
//
// For typed stores (the common case), declare [Provider][K, V]() directly
// inside your own module's [ligo.Providers] list instead.
//
//	app.Register(memory.Module(), myModule())
func Module() ligo.Module {
	return ligo.NewModule("ligo-memory",
		ligo.Providers(
			Provider[string, any](),
		),
	)
}
