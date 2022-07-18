package selector

import (
	"context"
	"strings"
)

type (
	// Colors represents a set of colors.
	Colors struct {
		colors []string
	}
	colorKey struct{}
)

// NewColors new a Colors.
func NewColors(colors ...string) *Colors {
	return &Colors{colors: append([]string(nil), colors...)}
}

// Add adds a set of colors.
func (c *Colors) Add(colors ...string) {
	c.colors = append(c.colors, colors...)
}

// Equal returns whether c and o are equivalent.
func (c *Colors) Equal(o interface{}) bool {
	if c == nil && o == nil {
		return true
	}
	if c == nil || o == nil {
		return false
	}

	var colors *Colors
	switch v := o.(type) {
	case *Colors:
		colors = v
	case Colors:
		colors = &v
	default:
		return false
	}

	if colors == nil {
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

// Colors returns a color slice.
func (c *Colors) Colors() []string {
	if len(c.colors) == 0 {
		return nil
	}

	cloneColors := make([]string, len(c.colors))
	copy(cloneColors, c.colors)
	return cloneColors
}

// Clone clones a Colors.
func (c *Colors) Clone() *Colors {
	return &Colors{colors: c.Colors()}
}

// Size returns size of the colors.
func (c *Colors) Size() int {
	return len(c.colors)
}

// Empty return ture if the length of colors is 0.
func (c *Colors) Empty() bool {
	return len(c.colors) == 0
}

// String returns a string representation.
func (c *Colors) String() string {
	return "[" + strings.Join(c.colors, ", ") + "]"
}

// NewColorsContext new a colors context.
func NewColorsContext(ctx context.Context, colors ...string) context.Context {
	if len(colors) == 0 {
		return ctx
	}

	return context.WithValue(ctx, colorKey{}, NewColors(colors...))
}

// ColorsFromContext get the current colors from the context.
func ColorsFromContext(ctx context.Context) *Colors {
	value := ctx.Value(colorKey{})
	if value == nil {
		return &Colors{}
	}

	return value.(*Colors)
}
