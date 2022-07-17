package selector

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

var (
	selectorMap       = make(map[string]Selector)
	ColorAttributeKey = attribute.Key("selector.color")
)

func Register(selector Selector) {
	selectorMap[selector.Name()] = selector
}

func Get(name string) (Selector, bool) {
	selector, ok := selectorMap[name]
	return selector, ok
}

type (
	selectKey struct{}
	colorKey  struct{}
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

func NewColorContext(ctx context.Context, colors ...string) context.Context {
	return context.WithValue(ctx, colorKey{}, NewColors(colors...))
}

func ColorFromContext(ctx context.Context) (*Colors, bool) {
	value := ctx.Value(colorKey{})
	if value == nil {
		return nil, false
	}

	return value.(*Colors), true
}

type (
	Selector interface {
		Select(conns []Conn, info balancer.PickInfo) []Conn
		Name() string
	}
	Conn interface {
		Address() resolver.Address
		SubConn() balancer.SubConn
	}
)

type Colors struct {
	colors []string
}

func NewColors(colors ...string) *Colors {
	return &Colors{colors: append([]string(nil), colors...)}
}

func (c *Colors) Add(colors ...string) {
	c.colors = append(c.colors, colors...)
}

func (c *Colors) Range(f func(color string) bool) {
	for _, color := range c.colors {
		if !f(color) {
			break
		}
	}
}

func (c *Colors) Equal(o interface{}) bool {
	var colors *Colors
	switch v := o.(type) {
	case *Colors:
		colors = v
	case Colors:
		colors = &v
	default:
		return false
	}

	if len(colors.colors) != len(c.colors) {
		return false
	}

	for i := 0; i < len(c.colors); i++ {
		if c.colors[i] != colors.colors[i] {
			return false
		}
	}

	return true
}

func (c Colors) Colors() []string {
	cloneColors := make([]string, len(c.colors))
	copy(cloneColors, c.colors)
	return cloneColors
}

func (c Colors) Size() int {
	return len(c.colors)
}
