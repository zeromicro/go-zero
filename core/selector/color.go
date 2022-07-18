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
	c := &Colors{colors: make([]string, 0, len(colors))}
	c.Add(colors...)

	return c
}

// Add adds a set of colors.
func (c *Colors) Add(colors ...string) {
	for _, color := range colors {
		if color == "" {
			continue
		}

		c.colors = append(c.colors, color)
	}
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
