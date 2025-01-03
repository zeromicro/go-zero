package httpx

import (
	"fmt"
	"net/http"
	"net/url"
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

func TestGetFormValues_TooManyValues(t *testing.T) {
	form := url.Values{}

	// Add more values than the limit
	for i := 0; i < maxFormParamCount+10; i++ {
		form.Add("param", fmt.Sprintf("value%d", i))
	}

	// Create a new request with the form data
	req, err := http.NewRequest("POST", "/test", strings.NewReader(form.Encode()))
	assert.NoError(t, err)

	// Set the content type for form data
	req.Header.Set(ContentType, "application/x-www-form-urlencoded")

	_, err = GetFormValues(req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too many form values")
}
