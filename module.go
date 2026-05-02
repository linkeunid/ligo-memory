package memory

import "github.com/linkeunid/ligo"

// Provider returns a ligo.Provider that registers a *Store[K, V] as a singleton
// in the DI container. Other factories can then receive it as an injected parameter.
//
// Example — declare in your module:
//
//	ligo.Providers(
//	    memory.Provider[string, *entity.User](),
//	    ligo.Factory[repository.UserRepository](NewUserRepository),
//	)
//
// And accept it in the constructor:
//
//	func NewUserRepository(store *memory.Store[string, *entity.User]) repository.UserRepository { ... }
func Provider[K comparable, V any]() ligo.Provider {
	return ligo.Factory[*Store[K, V]](func() *Store[K, V] {
		return New[K, V]()
	})
}

// Module returns a Ligo module that provides a general-purpose
// string-keyed in-memory store (*Store[string, any]) via DI.
// For typed stores, use Provider[K, V]() inside your own module's Providers list.
func Module() ligo.Module {
	return ligo.NewModule("memory",
		ligo.Providers(
			Provider[string, any](),
		),
	)
}
