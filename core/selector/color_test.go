package selector

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewColors(t *testing.T) {
	colors := NewColors("1", "2", "3")
	assert.Equal(t, []string{"1", "2", "3"}, colors.colors)
}

func TestColors_Add(t *testing.T) {
	colors := NewColors("1", "2", "3")
	colors.Add("3", "5", "6")
	assert.Len(t, colors.colors, 6)

	colors.Add("", "", "")
	assert.Len(t, colors.colors, 6)
	assert.Len(t, colors.colors, 6)
}

func TestColors_Clone(t *testing.T) {
	colors := NewColors("1", "2", "3")
	assert.Equal(t, []string{"1", "2", "3"}, colors.Clone().colors)
}

func TestColors_Empty(t *testing.T) {
	colors := NewColors()
	assert.True(t, colors.Empty())

	colors.Add("")
	assert.True(t, colors.Empty())

	colors.Add("1")
	assert.False(t, colors.Empty())
}

func TestColors_Size(t *testing.T) {
	colors := NewColors()
	assert.Equal(t, 0, colors.Size())

	colors.Add("1")
	assert.Equal(t, 1, colors.Size())
}

func TestColors_Colors(t *testing.T) {
	colors := NewColors()
	assert.Len(t, colors.Colors(), 0)

	colors.Add("1")
	assert.Equal(t, []string{"1"}, colors.Colors())
}

func TestColors_String(t *testing.T) {
	colors := NewColors()
	assert.Equal(t, "[]", colors.String())

	colors.Add("1")
	assert.Equal(t, "[1]", colors.String())

	colors.Add("2")
	assert.Equal(t, "[1, 2]", colors.String())
}

func TestNewColorsContext(t *testing.T) {
	ctx := NewColorsContext(context.Background(), "1", "2")
	assert.EqualValues(t, []string{"1", "2"}, ctx.Value(colorKey{}).(*Colors).colors)

	assert.True(t, context.Background() == NewColorsContext(context.Background()))
}

func TestColorsFromContext(t *testing.T) {
	assert.Equal(t, []string(nil), ColorsFromContext(context.Background()).Colors())
	assert.Equal(t, []string{"1", "2"}, ColorsFromContext(context.WithValue(context.Background(), colorKey{}, NewColors("1", "2"))).Colors())

}
