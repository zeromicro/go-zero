package configurator

import (
	"sync"

	"github.com/zeromicro/go-zero/core/conf"
)

type (
	// unmarshalerRegistry is the registry for unmarshalers.
	unmarshalerRegistry struct {
		unmarshalers map[string]LoaderFn

		mu sync.RWMutex
	}

	// LoaderFn is the function type for loading configuration.
	LoaderFn func([]byte, any) error
)

var defaultRegistry *unmarshalerRegistry

func init() {
	defaultRegistry = &unmarshalerRegistry{
		unmarshalers: map[string]LoaderFn{
			"json": conf.LoadFromJsonBytes,
			"toml": conf.LoadFromTomlBytes,
			"yaml": conf.LoadFromYamlBytes,
		},
	}
}

// RegisterUnmarshaler registers an unmarshaler.
func RegisterUnmarshaler(name string, fn LoaderFn) {
	defaultRegistry.mu.Lock()
	defaultRegistry.unmarshalers[name] = fn
	defaultRegistry.mu.Unlock()
}

// Unmarshaler returns the unmarshaler by name.
func Unmarshaler(name string) (LoaderFn, bool) {
	defaultRegistry.mu.RLock()
	fn, ok := defaultRegistry.unmarshalers[name]
	defaultRegistry.mu.RUnlock()
	return fn, ok
}
