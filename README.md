# ligo-memory

A simple in-memory store for [Ligo](https://github.com/linkeunid/ligo), designed for fast testing and local development without external dependencies.

[![Go Version](https://img.shields.io/badge/go-1.21+-blue)](https://go.dev/dl)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Tests](https://img.shields.io/badge/tests-20%20passing-brightgreen)](https://github.com/linkeunid/ligo-memory)
[![No dependencies](https://img.shields.io/badge/dependencies-none-brightgreen)](go.mod)

## Install

```bash
go get github.com/linkeunid/ligo-memory
```

## Quick start

Register a typed store in your module and inject it into a constructor:

```go
import (
    "github.com/linkeunid/ligo"
    memory "github.com/linkeunid/ligo-memory"
)

func UserModule() ligo.Module {
    return ligo.NewModule("user",
        ligo.Providers(
            memory.Provider[string, *User](),
            ligo.Factory[UserRepository](NewUserRepository),
        ),
    )
}

func NewUserRepository(store *memory.Store[string, *User]) UserRepository {
    return &userRepo{store: store}
}
```

For a zero-config drop-in, use `memory.Module()` which registers a `*Store[string, any]`:

```go
app.Register(memory.Module(), myModule())
```

## See also

- [Store API & patterns](docs/features/store.md)
