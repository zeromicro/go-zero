package iox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirectInOut(t *testing.T) {
	restore, err := RedirectInOut()
	assert.Nil(t, err)
	defer restore()
}
