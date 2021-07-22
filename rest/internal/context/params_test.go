package context

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVars(t *testing.T) {
	expect := map[string]string{
		"a": "1",
		"b": "2",
	}
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	r = r.WithContext(context.WithValue(context.Background(), pathVars, expect))
	assert.EqualValues(t, expect, Vars(r))
}

func TestVarsNil(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	assert.Nil(t, Vars(r))
}

func TestContextKey(t *testing.T) {
	ck := contextKey("hello")
	assert.True(t, strings.Contains(ck.String(), "hello"))
}
