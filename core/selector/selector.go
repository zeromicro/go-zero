package selector

import (
	"context"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

var selectorMap = make(map[string]Selector)

func Register(selector Selector) {
	selectorMap[selector.Name()] = selector
}

func Get(name string) (Selector, bool) {
	selector, ok := selectorMap[name]
	return selector, ok
}

type (
	Conn interface {
		Address() resolver.Address
	}

	Selector interface {
		Select(conns []Conn, info balancer.PickInfo) []Conn
		Name() string
	}

	colorKey  struct{}
	selectKey struct{}
)

func NewSelectorContext(ctx context.Context, selectorName string) context.Context {
	return context.WithValue(ctx, selectKey{}, selectorName)
}

func SelectorFromContext(ctx context.Context) (string, bool) {
	value := ctx.Value(selectKey{})
	if value == nil {
		return "", false
	}

	return value.(string), true
}
