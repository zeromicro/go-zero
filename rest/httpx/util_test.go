package httpx

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRemoteAddr(t *testing.T) {
	host := "8.8.8.8"
	r, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	assert.Nil(t, err)

	r.Header.Set(xForwardedFor, host)
	assert.Equal(t, host, GetRemoteAddr(r))
}

func TestGetRemoteAddrNoHeader(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "/", strings.NewReader(""))
	assert.Nil(t, err)

	assert.True(t, len(GetRemoteAddr(r)) == 0)
}

func Test_filterFormValues(t *testing.T) {
	values, valid := filterFormValues([]string{"1", "", "2"})
	assert.Equal(t, []string{"1", "2"}, values)
	assert.True(t, valid)

	values, valid = filterFormValues([]string{"1", ""})
	assert.Equal(t, "1", values)
	assert.True(t, valid)

	values, valid = filterFormValues([]string{""})
	assert.Equal(t, nil, values)
	assert.False(t, valid)
}
