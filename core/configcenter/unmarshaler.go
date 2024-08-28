package configurator

import (
	"sync"

	"github.com/zeromicro/go-zero/core/conf"
)

var registry = &unmarshalerRegistry{
	unmarshalers: map[string]LoaderFn{
		"json": conf.LoadFromJsonBytes,
		"toml": conf.LoadFromTomlBytes,
		"yaml": conf.LoadFromYamlBytes,
	},
}

type (
	// LoaderFn is the function type for loading configuration.
	LoaderFn func([]byte, any) error

	// unmarshalerRegistry is the registry for unmarshalers.
	unmarshalerRegistry struct {
		unmarshalers map[string]LoaderFn
		mu           sync.RWMutex
	}
)

// RegisterUnmarshaler registers an unmarshaler.
func RegisterUnmarshaler(name string, fn LoaderFn) {
	registry.mu.Lock()
	defer registry.mu.Unlock()
	registry.unmarshalers[name] = fn
}

// Unmarshaler returns the unmarshaler by name.
func Unmarshaler(name string) (LoaderFn, bool) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()
	fn, ok := registry.unmarshalers[name]
	return fn, ok
}
