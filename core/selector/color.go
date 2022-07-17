package selector

import (
	"context"
	"strings"
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

func (c *Colors) Colors() []string {
	cloneColors := make([]string, len(c.colors))
	copy(cloneColors, c.colors)
	return cloneColors
}

func (c *Colors) Clone() *Colors {
	return &Colors{colors: c.Colors()}
}

func (c *Colors) Size() int {
	return len(c.colors)
}

func (c *Colors) String() string {
	return "[" + strings.Join(c.colors, ", ") + "]"
}

func NewColorsContext(ctx context.Context, colors ...string) context.Context {
	return context.WithValue(ctx, colorKey{}, NewColors(colors...))
}

func ColorsFromContext(ctx context.Context) (*Colors, bool) {
	value := ctx.Value(colorKey{})
	if value == nil {
		return nil, false
	}

	return value.(*Colors), true
}
