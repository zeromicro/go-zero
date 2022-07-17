package selector

import (
	"context"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

// selectorMap available selectors.
var selectorMap = make(map[string]Selector)

// Register registers a selector.
func Register(selector Selector) {
	selectorMap[selector.Name()] = selector
}

// Get get a selector.
func Get(name string) Selector {
	selector, ok := selectorMap[name]
	if !ok {
		return noneSelector{}
	}

	return selector
}

type (
	// Conn represents a gRPC connection.
	Conn interface {
		// Address returns a server the client connects to.
		Address() resolver.Address
	}

	// Selector represents a selector.
	Selector interface {
		// Select returns a callable connection
		Select(conns []Conn, info balancer.PickInfo) []Conn
		// Name returns a selector name.
		Name() string
	}

	colorKey  struct{}
	selectKey struct{}
)

// NewSelectorContext new a selector context.
func NewSelectorContext(ctx context.Context, selectorName string) context.Context {
	if selectorName == "" {
		return ctx
	}

	return context.WithValue(ctx, selectKey{}, selectorName)
}

// SelectorFromContext get the current selector from the context.
func SelectorFromContext(ctx context.Context) string {
	value := ctx.Value(selectKey{})
	if value == nil {
		return ""
	}

	return value.(string)
}
