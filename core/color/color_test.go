package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithColor(t *testing.T) {
	output := WithColor("Hello", BgRed)
	expected := BgRed + "Hello" + Reset
	assert.Equal(t, expected, output)
}

func TestWithColorPadding(t *testing.T) {
	output := WithColorPadding("Hello", BgRed)
	expected := BgRed + " Hello " + Reset
	assert.Equal(t, expected, output)
}
