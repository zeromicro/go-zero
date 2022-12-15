package selector

import (
	"context"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

// selectorMap available selectors.
var selectorMap = make(map[string]Selector)

type (
	// Conn represents a gRPC connection.
	Conn interface {
		// Address returns a server the client connects to.
		Address() resolver.Address
	}

	// Selector represents a selector.
	Selector interface {
		// Name returns a selector name.
		Name() string
		// Select returns a callable connection
		Select(conns []Conn, info balancer.PickInfo) []Conn
	}

	selectKey struct{}
)

// Get get a selector.
func Get(name string) Selector {
	selector, ok := selectorMap[name]
	if !ok {
		return noneSelector{}
	}

	return selector
}

// Register registers a selector.
func Register(selector Selector) {
	selectorMap[selector.Name()] = selector
}

// NewContext new a selector context.
func NewContext(ctx context.Context, selectorName string) context.Context {
	if selectorName == "" {
		return ctx
	}

	return context.WithValue(ctx, selectKey{}, selectorName)
}

// FromContext get the current selector from the context.
func FromContext(ctx context.Context) string {
	value := ctx.Value(selectKey{})
	if value == nil {
		return ""
	}

	return value.(string)
}
