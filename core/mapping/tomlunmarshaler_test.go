package mapping

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalToml(t *testing.T) {
	const input = `a = "foo"
b = 1
c = "${FOO}"
d = "abcd!@#$112"
`
	var val struct {
		A string `json:"a"`
		B int    `json:"b"`
		C string `json:"c"`
		D string `json:"d"`
	}
	assert.NoError(t, UnmarshalTomlReader(strings.NewReader(input), &val))
	assert.Equal(t, "foo", val.A)
	assert.Equal(t, 1, val.B)
	assert.Equal(t, "${FOO}", val.C)
	assert.Equal(t, "abcd!@#$112", val.D)
}

func TestUnmarshalTomlErrorToml(t *testing.T) {
	const input = `foo"
b = 1
c = "${FOO}"
d = "abcd!@#$112"
`
	var val struct {
		A string `json:"a"`
		B int    `json:"b"`
		C string `json:"c"`
		D string `json:"d"`
	}
	assert.Error(t, UnmarshalTomlReader(strings.NewReader(input), &val))
}

func TestUnmarshalTomlBadReader(t *testing.T) {
	var val struct {
		A string `json:"a"`
	}
	assert.Error(t, UnmarshalTomlReader(new(badReader), &val))
}
