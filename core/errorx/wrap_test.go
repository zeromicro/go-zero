package errorx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	assert.Nil(t, Wrap(nil, "test"))
	assert.Equal(t, "foo: bar", Wrap(errors.New("bar"), "foo").Error())
}

func TestWrapf(t *testing.T) {
	assert.Nil(t, Wrapf(nil, "%s", "test"))
	assert.Equal(t, "foo bar: quz", Wrapf(errors.New("quz"), "foo %s", "bar").Error())
}
