package selector

import (
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

}
